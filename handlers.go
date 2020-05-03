package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

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
	fmt.Fprintf(w, "Category: %v\n", vars["id"])

}

// RegisterHandler to Register a new user
func (a *App) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	// grab user info
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	firstname := r.FormValue("first_name")
	lastname := r.FormValue("last_name")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	checkInternalServerError(err, w)
	// Check existence of user
	var user User
	if a.DB.Where("username = ?", username).First(&user).RecordNotFound() {
		user.Username = username
		user.Password = string(hashedPassword)
		user.Email = email
		user.FirstName = firstname
		user.LastName = lastname
		a.DB.Create(&user)
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode("Success")
		return
	}

	json.NewEncoder(w).Encode("Record Found")
	// http.Redirect(w, r, "/", 301)

}

// LoginHandler has some data
func (a *App) LoginHandler(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	var user User

	if a.DB.Where("username = ?", username).First(&user).RecordNotFound() {
		json.NewEncoder(w).Encode("Account not Found")
		// http.Redirect(w, r, "/login", 301)
		return
	}

	a.DB.Where("username = ?", username).First(&user)
	// validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		json.NewEncoder(w).Encode("Password Doesn't Match")
		// http.Redirect(w, r, "/login", 301)
		return
	}
	authenticated = true
	json.NewEncoder(w).Encode("Login Successful")
	// http.Redirect(w, r, "/list", 301)

}
