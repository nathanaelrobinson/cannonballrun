package main

import (
	"fmt"
	"net/http"
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

func (a *App) updateWorkouts(jobChan <-chan *Webhook) {
	for webhook := range jobChan {
		// Retrieve Strava Users
		var athlete StravaAthlete

		// Grap the User of the workout and determine if their access token has expired.
		a.DB.Where("strava_id = ?", webhook.OwnerID).First(&athlete)
		if int64(athlete.ExpiresAt) > time.Now().Unix() {
			// Handle expired Access Token
		}
		// Use Updated Access Token to create client
		client := strava.NewClient(athlete.AccessToken)
		service := strava.NewActivitiesService(client)
		activity, err := service.Get(int64(webhook.ObjectId)).IncludeAllEfforts().Do()
		if err != nil {
			fmt.Println(err)
		}

		var workout Workout
		workout.Distance = activity.Distance
		workout.StravaID = int(activity.Id)
		// workout.Type = activity.Type
		workout.Time = int64(activity.ElapsedTime)
		//Update, Save, and Exit

	}

}
