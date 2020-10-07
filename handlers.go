package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// TeamHandlerList provides a list view of Teams
func (a *App) TeamHandlerList(w http.ResponseWriter, r *http.Request) {
	var teams []Team
	a.DB.Preload("Runners").Find(&teams)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(&teams)
}

// TeamHandlerDetail provides a list view of Teams
func (a *App) TeamHandlerDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var team Team
	switch {
	case r.Method == "GET":
		a.DB.First(&team, vars["id"])
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(&team)
		return
	case r.Method == "PUT":
		a.DB.First(&team, vars["id"])
		team.Name = r.FormValue("name")
		team.Description = r.FormValue("description")
		ownerID, _ := strconv.ParseUint(r.FormValue("owner_id"), 10, 64)
		team.UserID = uint(ownerID)
		a.DB.Save(&team)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(&team)
		return
	case r.Method == "POST":
		OwnerID := r.Context().Value("user_id")
		team.Name = r.FormValue("name")
		team.Description = r.FormValue("description")
		team.UserID = OwnerID.(uint)
		a.DB.Create(&team)
		a.DB.Create(&Affiliation{TeamID: team.ID, UserID: team.UserID, Status: "Active"})
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(&team)
		return
	case r.Method == "DELETE":
		a.DB.Where("id = ?", vars["id"]).Delete(&Team{})
		json.NewEncoder(w).Encode("Deleted Record")
		return
	default:
		fmt.Println("Fatal Error")
	}
	w.WriteHeader(http.StatusOK)

}

// TeamHandlerDetail provides a list view of Teams
func (a *App) TeamAthletes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var runners []RunnerOutput
	a.DB.Table("workouts").Select("users.id as user_id, users.username, count(distinct workouts.id) as activities_count, sum(workouts.distance) as total_distance").Joins("left join users on users.id = workouts.user_id").Joins("left join affiliations on affiliations.user_id = users.id").Group("users.id, users.username").Where("affiliations.status = 'Active' and workouts.type = 'Run' and affiliations.team_id = ?", vars["id"]).Order("sum(workouts.distance) desc").Scan(&runners)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(&runners)
	return

}

// TeamHandlerList provides a list view of Teams
func (a *App) JoinTeam(w http.ResponseWriter, r *http.Request) {
	userIDInt := r.Context().Value("user_id")
	vars := mux.Vars(r)
	teamID, _ := strconv.ParseInt(vars["id"], 10, 64)
	var affiliation Affiliation
	affiliation.UserID = userIDInt.(uint)
	affiliation.TeamID = uint(teamID)
	affiliation.Status = "Pending"

	a.DB.Create(&affiliation)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
}

// RegisterHandler to Register a new user
func (a *App) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	// grab user info
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	stravaID := r.FormValue("strava_id")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	checkInternalServerError(err, w)
	// Check existence of user
	var athlete StravaAthlete
	a.DB.Where("strava_id = ?", stravaID).First(&athlete)

	var user User
	if a.DB.Where("username = ?", username).First(&user).RecordNotFound() {
		user.Username = username
		user.Password = string(hashedPassword)
		user.Email = email
		a.DB.Create(&user)
		athlete.UserID = user.ID
		a.DB.Save(&athlete)

		expiresAt := time.Now().Add(time.Minute * 100000).Unix()

		tk := Token{
			UserID:   user.ID,
			Username: user.Username,
			Email:    user.Email,
			StandardClaims: &jwt.StandardClaims{
				ExpiresAt: expiresAt,
			},
		}
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
		tokenString, error := token.SignedString([]byte("secret"))
		if error != nil {
			fmt.Println(error)
		}
		var resp = map[string]interface{}{"status": 200, "message": "logged in"}
		resp["token"] = tokenString //Store the token in the response
		resp["user"] = user
		json.NewEncoder(w).Encode(resp)
		return
	}

	var resp = map[string]interface{}{"status": 403, "message": "Username or Email Already Exists"}
	json.NewEncoder(w).Encode(resp)
	// http.Redirect(w, r, "/", 301)

}

// LoginHandler takes username and password form data. If this data matches the hashed value in the db, we return login credentials via
// a JWT and redirect the user to the home page. Otherwise we return the provided error.
func (a *App) LoginHandler(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	var user User

	if a.DB.Where("username = ?", username).First(&user).RecordNotFound() {
		var resp = map[string]interface{}{"status": 403, "message": "Invalid login credentials. Please try again"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	a.DB.Where("username = ?", username).First(&user)
	// validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		var resp = map[string]interface{}{"status": 403, "message": "Invalid login credentials. Please try again"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	expiresAt := time.Now().Add(time.Minute * 100000).Unix()

	tk := Token{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}
	var resp = map[string]interface{}{"status": 200, "message": "logged in"}
	resp["token"] = tokenString //Store the token in the response
	resp["user"] = user
	json.NewEncoder(w).Encode(resp)
	return

}

// TeamHandlerList provides a list view of Teams
func (a *App) UserHandlerList(w http.ResponseWriter, r *http.Request) {
	var users []User
	a.DB.Preload("StravaAthlete").Find(&users)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(&users)
}

// TeamHandlerList provides a list view of Teams
func (a *App) UserHandlerDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDInt := r.Context().Value("user_id")
	data := vars["data"]
	switch {
	case data == "affiliations":
		var teams []Team
		a.DB.Joins("left join affiliations on affiliations.team_id = teams.id").Where("affiliations.user_id = ? and affiliations.status = 'Active'", userIDInt).Find(&teams)
		json.NewEncoder(w).Encode(&teams)
	case data == "details":
		var user User
		a.DB.Preload("Teams").Preload("AdminTeams").Preload("Workouts", "distance > 0 and type  = 'Run'").Preload("StravaAthlete").First(&user, userIDInt)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(&user)
	case data == "admins":
		var teams []Team
		a.DB.Where("user_id = ?", userIDInt).Find(&teams)
		json.NewEncoder(w).Encode(&teams)
	case data == "workouts":
		var workouts []Workout
		a.DB.Where("distance > 0 and type  = 'Run' and user_id = ?", userIDInt).Order("created_at desc").Find(&workouts)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(&workouts)
	}

	// a.DB.Preload("Teams").Preload("AdminTeams").Preload("Workouts", "distance > 0 and type  = 'Run'").Preload("StravaAthlete").First(&user, userIDInt)
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// json.NewEncoder(w).Encode(&user)
}

// TeamHandlerList provides a list view of Teams
func (a *App) TeamsStatus(w http.ResponseWriter, r *http.Request) {
	var teams []TeamOutput
	a.DB.Table("workouts").Select("teams.id as team_id, teams.name, teams.description, sum(workouts.distance) as distance").Joins("left join affiliations on affiliations.user_id = workouts.user_id").Joins("left join teams on teams.id = affiliations.team_id").Group("teams.id, teams.name, teams.description").Where("affiliations.status = 'Active' and workouts.type = ? and teams.id is not null", "Run").Order("sum(workouts.distance) desc").Scan(&teams)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(&teams)
}

// StravaAuthorization is an endpoint that handles Strava URI direct for OAuth Authorization
func (a *App) StravaAuthorization(w http.ResponseWriter, r *http.Request) {

	// Retreive GET params from Strava URI URL and prepare post payload
	code := string(r.URL.Query()["code"][0])
	baseURL, _ := url.Parse("https://www.strava.com/")
	baseURL.Path = "api/v3/oauth/token"
	params := url.Values{}
	params.Add("code", code)
	params.Add("client_id", string(os.Getenv("STRAVA_CLIENT_ID")))
	params.Add("client_secret", string(os.Getenv("STRAVA_CLIENT_SECRET")))
	params.Add("grant_type", "authorization_code")
	baseURL.RawQuery = params.Encode()

	// Post Auth Code and accompanying payload info to Strava API to retreive Refresh Token and Access Token
	resp, err := http.Post(baseURL.String(), "application/json", nil)
	checkInternalServerError(err, w)
	body, err := ioutil.ReadAll(resp.Body)
	checkInternalServerError(err, w)

	// Unmarshall Json
	var responseData JsonAuthReponse
	json.Unmarshal(body, &responseData)
	fmt.Printf("%+v\n", responseData)

	// Apply JSON to StravaAthlete object
	var athlete StravaAthlete
	athlete.StravaID = responseData.Athlete.ID
	athlete.AccessToken = responseData.AccessToken
	athlete.RefreshToken = responseData.RefreshToken
	athlete.ExpiresAt = responseData.ExpiresAt
	athlete.UserName = responseData.Athlete.Username
	athlete.ProfileImage = responseData.Athlete.Profile

	// Create an object to determine if there is already a strava athlete by this ID. If there is not,
	// create a new strava athlete and serve the landing page. Otherwise send the previous athlete
	var previousAthlete StravaAthlete
	if a.DB.Where("strava_id = ?", athlete.StravaID).First(&previousAthlete).RecordNotFound() {
		a.DB.Create(&athlete)
		tmpl, _ := template.ParseFiles("./templates/strava_landing.html")
		tmpl.Execute(w, athlete)
		return
	} else {
		// var athlete StravaAthlete
		tmpl, _ := template.ParseFiles("./templates/strava_landing.html")
		tmpl.Execute(w, previousAthlete)
		return
	}

}

// CreateStravaWebhook is a helper function to first create a webhook, should no longer be required
func (a *App) CreateStravaWebhook(w http.ResponseWriter, r *http.Request) {
	hubChallenge := string(r.URL.Query()["hub.challenge"][0])
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var resp = map[string]interface{}{"hub.challenge": hubChallenge}
	json.NewEncoder(w).Encode(&resp)
}

// StravaWebhookHandler takes post requests from Strava and updates workout information
func (a *App) StravaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Immediately write 200 header before handling data
	w.WriteHeader(200)

	// Retreview webhook information to a struct, send a copy of that struct to a channel for handling by a concurrent thread
	var webhook Webhook
	body, err := ioutil.ReadAll(r.Body)
	checkInternalServerError(err, w)
	json.Unmarshal(body, &webhook)

	if webhook.ObjectType == "activity" {
		if webhook.AspectType == "create" {

			// Send webhook to job channel to create workout
			jobChan <- webhook

		} else if webhook.AspectType == "delete" {
			// TODO: Handle Deletes concurrently like creations above (although this is less important)
			a.DB.Where("strava_id = ?", int(webhook.ObjectId)).Delete(Workout{})
		}
		// TODO: Account for Updates to workouts via strava
	}

}
