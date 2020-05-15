package main

import (
	_ "database/sql"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

// Workout is a struct defining a workout information delivered from Strava API
type Workout struct {
	gorm.Model
	StravaID int     `json:"strava_id,omitempty"`
	Distance float64 `json:"distance,omitempty"`
	UserID   uint    `json:"user_id"`
	Time     int64   `json:"time"`
	Type     string  `json:"type"`
}

// User is a struct defining a User on the site and stores Strava link information
type User struct {
	gorm.Model
	StravaAthlete StravaAthlete
	Username      string `gorm:"size:255" json:"username,omitempty"`
	Email         string `gorm:"size:255" json:"email,omitempty"`
	Password      string `gorm:"size:255" json:"password,omitempty"`

	Teams      []Affiliation
	AdminTeams []Team
	Workouts   []Workout
}

type StravaAthlete struct {
	gorm.Model
	UserID       uint   `json:"user_id,omitempty"`
	StravaID     int    `json:"strava_id,omitempty"`
	ExpiresAt    int    `json:"expires_at,omitempty"`
	RefreshToken string `gorm:"size:255" json:"refresh_token,omitempty"`
	AccessToken  string `gorm:"size:255" json:"access_token,omitempty"`
	UserName     string `gorm:"size:255" json:"username,omitempty"`
	ProfileImage string `gorm:"size:255" json:"profile,omitempty"`
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

type Webhook struct {
	StravaAthleteID uint        `json:"user_id"`
	AspectType      string      `json:"aspect_type"`
	EventTime       float64     `json:"event_time"`
	ObjectId        float64     `json:"object_id"`
	ObjectType      string      `json:"object_type"`
	OwnerID         float64     `json:"owner_id"`
	Distance        float64     `json:"distance"`
	Updates         interface{} `json:"updates"`
}

// Initialize opens our DB connectionn
func (a *App) Initialize(dbDriver string, dbURI string) {
	db, err := gorm.Open(dbDriver, dbURI)
	db.LogMode(true)
	if err != nil {
		panic("failed to connect database")
	}
	a.DB = db
}

// func main() {
// 	// var (
// 	// 	dbUser                 = "nate"
// 	// 	dbPwd                  = "ycombinator"
// 	// 	instanceConnectionName = "cannonballrun:us-central1:cannonballrun"
// 	// 	dbName                 = "test"
// 	// )
// 	// var dbURI string
// 	// dbURI = fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s", dbUser, dbPwd, instanceConnectionName, dbName)
// 	// fmt.Println(dbURI)
// 	db, err := gorm.Open("mysql", "nate:ycombinator@tcp(35.188.216.81:3306)/test?charset=utf8&parseTime=True&loc=UTC")
// 	if err != nil {
// 		panic("failed to connect database")
// 	}
// 	defer db.Close()
//
// 	// Migrate the schema
// 	db.AutoMigrate(&Workout{})
// 	db.AutoMigrate(&User{})
// 	db.AutoMigrate(&StravaAthlete{})
// 	db.AutoMigrate(&Team{})
// 	db.AutoMigrate(&Affiliation{})
//
// 	// Create
// 	nate := User{Username: "Naterob", Email: "natrobsn@gmail.com", Password: "some text"}
// 	db.Create(&nate)
// 	strv := StravaAthlete{UserID: nate.ID, StravaID: 2222222, ExpiresAt: 199999, RefreshToken: "somecrazystring", AccessToken: "someotherlongassstring", UserName: "nathanael_robinson", ProfileImage: "https://someurl.com/"}
// 	db.Create(&strv)
//
// 	nutc := Team{Name: "NUTC", Description: "Jolly band of geeds", UserID: nate.ID}
//
// 	db.Create(&User{Username: "the Bear", Email: "joncohen@.com", Password: "password"})
// 	db.Create(&User{Username: "achey", Email: "fief@gmail.com", Password: "password"})
// 	db.Create(&User{Username: "dictator", Email: "greatleader", Password: "password"})
//
// 	//
// 	db.Create(&nutc)
//
// 	db.Create(&Affiliation{TeamID: nutc.ID, UserID: nate.ID, Status: "Active"})
//
// }
