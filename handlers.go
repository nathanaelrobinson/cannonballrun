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
		a.DB.Preload("Runners", "status = 'Active'").First(&team, vars["id"])
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
		team.Name = r.FormValue("name")
		team.Description = r.FormValue("description")
		ownerID, _ := strconv.ParseUint(r.FormValue("owner_id"), 10, 64)
		team.UserID = uint(ownerID)
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

// TeamHandlerList provides a list view of Teams
func (a *App) JoinTeam(w http.ResponseWriter, r *http.Request) {
	var team Team
	a.DB.Preload("Runners").Find(&team)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(&team)
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

// LoginHandler has some data
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
	userIDInt := r.Context().Value("user_id")
	var user User
	a.DB.Preload("Teams").Preload("AdminTeams").Preload("Workouts", "distance > 0").Preload("StravaAthlete").First(&user, userIDInt)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(&user)
}

// StravaAuthorization is an endpoint that handles Strava URI direct for OAuth Authorization
func (a *App) StravaAuthorization(w http.ResponseWriter, r *http.Request) {
	code := string(r.URL.Query()["code"][0])
	baseURL, _ := url.Parse("https://www.strava.com/")
	baseURL.Path = "api/v3/oauth/token"
	params := url.Values{}
	params.Add("code", code)
	params.Add("client_id", string(os.Getenv("STRAVA_CLIENT_ID")))
	params.Add("client_secret", string(os.Getenv("STRAVA_CLIENT_SECRET")))
	params.Add("grant_type", "authorization_code")
	baseURL.RawQuery = params.Encode()
	fmt.Printf(baseURL.String())
	// var athlete StravaAthlete

	resp, err := http.Post(baseURL.String(), "application/json", nil)
	checkInternalServerError(err, w)
	body, err := ioutil.ReadAll(resp.Body)
	checkInternalServerError(err, w)

	responseData := make(map[string]interface{})
	json.Unmarshal(body, &responseData)
	innerAthlete := responseData["athlete"].(map[string]interface{})

	var athlete StravaAthlete
	athlete.AccessToken = responseData["access_token"].(string)
	athlete.RefreshToken = responseData["refresh_token"].(string)
	athlete.ExpiresAt = int(responseData["expires_at"].(float64))
	athlete.StravaID = int(innerAthlete["id"].(float64))
	athlete.UserName = innerAthlete["username"].(string)
	athlete.ProfileImage = innerAthlete["profile"].(string)

	var previousAthlete StravaAthlete
	if a.DB.Where("strava_id = ?", athlete.StravaID).First(&previousAthlete).RecordNotFound() {
		a.DB.Save(&athlete)
		tmpl, _ := template.ParseFiles("./templates/strava_landing.html")
		tmpl.Execute(w, athlete)
	}

	// var athlete StravaAthlete
	tmpl, _ := template.ParseFiles("./templates/strava_landing.html")
	tmpl.Execute(w, previousAthlete)

}

// TeamHandlerList provides a list view of Teams
func (a *App) CreateStravaWebhook(w http.ResponseWriter, r *http.Request) {
	hubChallenge := string(r.URL.Query()["hub.challenge"][0])
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var resp = map[string]interface{}{"hub.challenge": hubChallenge}
	json.NewEncoder(w).Encode(&resp)
}

// StravaWebhookHandler takes post requests from Strava and updates workout information
func (a *App) StravaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	var webhook Webhook
	// var workout Workout
	body, err := ioutil.ReadAll(r.Body)
	checkInternalServerError(err, w)
	json.Unmarshal(body, &webhook)
	var stravaathlete StravaAthlete
	a.DB.Where("strava_id = ?", int(webhook.OwnerID)).First(&stravaathlete)
	// Currently do not save workout titles/types
	// updates := make(map[string]interface{})
	// if webhook.Updates != nil {
	// 	updates = webhook.Updates.(map[string]interface{})
	// }

	var workout Workout
	if webhook.ObjectType == "activity" {
		if webhook.AspectType == "create" {
			workout.StravaID = int(webhook.ObjectId)
			workout.UserID = stravaathlete.UserID
			workout.Distance = webhook.Distance
			// Need to Enable a function to grab the workout time and distance
			a.DB.Create(&workout)
		} else if webhook.AspectType == "delete" {
			a.DB.Where("strava_id = ?", webhook.ObjectId).Delete(Workout{})
		} else if webhook.AspectType == "update" {
			//pass
		}
	}
	w.WriteHeader(200)
}
