package main

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Workout is a struct defining a workout information delivered from Strava API
type Workout struct {
	gorm.Model
	StravaID int     `json:"strava_id,omitempty"`
	Distance float64 `json:"distance,omitempty"`
	UserID   uint    `json:"user_id"`
}

// User is a struct defining a User on the site and stores Strava link information
type User struct {
	gorm.Model
	StravaID  int    `json:"strava_id,omitempty"`
	Username  string `gorm:"size:255" json:"username,omitempty"`
	FirstName string `gorm:"size:255" json:"first_name,omitempty"`
	LastName  string `gorm:"size:255" json:"last_name,omitempty"`
	Email     string `gorm:"size:255" json:"email,omitempty"`
	Password  string `gorm:"size:255" json:"password,omitempty"`

	Teams      []Affiliation
	AdminTeams []Team
	Workouts   []Workout
}

// Team is a struct defining a Team affiliation Users may join
type Team struct {
	gorm.Model
	Name        string `gorm:"size:255" json:"name,omitempty"`
	Description string `gorm:"size:2000" json:"description,omitempty"`
	UserID      uint   `json:"admin_id,omitempty"`
	maxSize     int    `gorm:"default:10"`
	Runners     []Affiliation
}

// Affiliation is a struct defining a relationship between a User and a Team
type Affiliation struct {
	gorm.Model
	TeamID uint   `json:"team_id"`
	UserID uint   `json:"user_id"`
	Status string `gorm:"size:12" json:"status"` // Options include Pending, Active, Removed
}

//Token struct declaration
type Token struct {
	UserID   uint
	Username string
	Email    string
	*jwt.StandardClaims
}

// Exceptions send messages to frontend via json
type Exception struct {
	Message string `json:"message"`
}

// App is used to initialize a database and hold our handler functions
type App struct {
	DB *gorm.DB
}

// Initialize opens our DB connectionn
func (a *App) Initialize(dbDriver string, dbURI string) {
	db, err := gorm.Open(dbDriver, dbURI)
	if err != nil {
		panic("failed to connect database")
	}
	a.DB = db
}

// func main() {
// 	db, err := gorm.Open("sqlite3", "test.db")
// 	if err != nil {
// 		panic("failed to connect database")
// 	}
// 	defer db.Close()
//
// 	// Migrate the schema
// 	db.AutoMigrate(&Workout{})
// 	db.AutoMigrate(&User{})
// 	db.AutoMigrate(&Team{})
// 	db.AutoMigrate(&Affiliation{})
//
// 	// Create
// 	nate := User{StravaID: 12345, Username: "Naterob", FirstName: "Nate", LastName: "Robinson", Email: "natrobsn@gmail.com", Password: "some text"}
// 	nutc := Team{Name: "NUTC", Description: "Jolly band of geeds", UserID: nate.ID}
// 	db.Create(&nate)
// 	db.Create(&User{StravaID: 12346, Username: "the Bear", FirstName: "Jon", LastName: "Cohen", Email: "joncohen@.com", Password: "password"})
// 	db.Create(&User{StravaID: 12347, Username: "achey", FirstName: "Andrew", LastName: "Pfeifer", Email: "fief@gmail.com", Password: "password"})
// 	db.Create(&User{StravaID: 12347, Username: "dictator", FirstName: "Luis", LastName: "Sahores", Email: "greatleader", Password: "password"})
//
// 	//
// 	db.Create(&nutc)
//
// 	db.Create(&Affiliation{TeamID: nutc.ID, UserID: nate.ID, Status: "Active"})
//
// }
