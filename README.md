# cannonballrun
Cannonball run is a webapp for a virtual road race from New York to Los Angeles.

Traditionally The Cannonball Run is an illegal street race in which participants attempt to drive from the Red Bull Garage
on E 31st Street of Manhattan to the Portfonio Hotel in Redondo Beach, CA. 
See Wikipedia Link Here: https://en.wikipedia.org/wiki/Cannonball_Run_Challenge

This webapp attempts to recreate that race in a virtual footrace setting using Strava.
Runners may log in with their Strava account and join a team. Every time they submit a run to Strava their
workout will automatically get added to their team's relay on the cannonballrun site. 

From the site runners can view the progress of their team against others, see races statistics and create more teams.

This application is powered in the backend by Golang, utilizing GORM as an ORM, Mux for the server, and Strava's own go library for interfacing with their API. 

Please reach out if you're interested in contributing!


