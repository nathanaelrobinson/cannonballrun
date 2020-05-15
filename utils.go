package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	strava "github.com/strava/go.strava"
)

func checkInternalServerError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

func (a *App) UpdateWorkouts(jobChan <-chan Webhook) {
	for webhook := range jobChan {
		// Retrieve Strava Users
		fmt.Println("Reached Function")
		var athlete StravaAthlete
		a.DB.Where("strava_id = ?", webhook.OwnerID).First(&athlete)
		fmt.Printf("%v", athlete)
		// Grab the User of the workout and determine if their access token has expired.

		if int64(athlete.ExpiresAt) < time.Now().Unix() {
			baseURL, _ := url.Parse("https://www.strava.com/")
			baseURL.Path = "api/v3/oauth/token"
			params := url.Values{}
			params.Add("client_id", string(os.Getenv("STRAVA_CLIENT_ID")))
			params.Add("client_secret", string(os.Getenv("STRAVA_CLIENT_SECRET")))
			params.Add("grant_type", "refresh_token")
			params.Add("refresh_token", athlete.RefreshToken)
			baseURL.RawQuery = params.Encode()
			fmt.Printf(baseURL.String())
			// var athlete StravaAthlete
			resp, _ := http.Post(baseURL.String(), "application/json", nil)
			body, _ := ioutil.ReadAll(resp.Body)
			responseData := make(map[string]interface{})
			json.Unmarshal(body, &responseData)

			athlete.AccessToken = responseData["access_token"].(string)
			athlete.RefreshToken = responseData["refresh_token"].(string)
			athlete.ExpiresAt = int(responseData["expires_at"].(float64))

			a.DB.Save(&athlete)

		}
		// Use Updated Access Token to create client
		var workout Workout

		client := strava.NewClient(athlete.AccessToken)
		service := strava.NewActivitiesService(client)

		activity, err := service.Get(int64(webhook.ObjectId)).IncludeAllEfforts().Do()
		if err == nil {
			workout.Distance = activity.Distance
			workout.StravaID = int(activity.Id)
			workout.Type = activity.Type.String()
			workout.Time = int64(activity.ElapsedTime)
			workout.UserID = athlete.UserID
			a.DB.Create(&workout)
		}
		//Update, Save, and Exit

	}

}
