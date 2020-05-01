package main

import (
	"fmt"

	"github.com/strava/go.strava"
)

const access_token = "b5df91b19333984bfc5c3b188aac1ec65cb75af9"

// const client_secret = "6b4b568a28f85980acb38c4fc58b8143b9b93b6d"
// const refresh_token = "e62ba8af991dd981742355f4613dee64483a5b44"

func ConnectStravaApi() *strava.Client {

	// Connect to Strava's API and return a Client Token
	client := strava.NewClient(access_token)
	return client

}

func GetTeamDistance(client *strava.Client, team int64) float64 {
	service := strava.NewClubsService(client)
	activities, err := service.ListActivities(team).Page(1).PerPage(100).Do()
	if err != nil {
		fmt.Println(err)
	}
	var totalDistance float64 = 0
	for _, activity := range activities {
		totalDistance += activity.Distance
	}
	return totalDistance / 1605.0

}

func main() {

	fmt.Println(GetTeamDistance(ConnectStravaApi(), 516715))

}
