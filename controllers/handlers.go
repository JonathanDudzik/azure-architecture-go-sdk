package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/jonathandudzik/gin_session_auth/helpers"
)

// Other actions (routes) in the app will have the proper authorization checks and be directed to Redis
// There should be these routes which should have handlers:
// POST /sign-in
// POST: /sign-up
// POST: /sign-out
// Serverless function(?) that authRequired() should block route if no session OR if no proper permission
// TO remeber, declare var in if statement and export struct fields!!!

type UserSignInDetails struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserSignUpDetails struct {
	Name string `json:"name" binding:"required"`
	UserSignInDetails
}

type DataBase struct {
	Key      string
	Name     string
	Password string
	Email    string
}

type UserSessionDetails struct {
	Token string
	Name  string
	Email string
}

func SignInHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// receive the username and password
		json := UserSignInDetails{}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Println("Retreived user sign-in details:", json)

		// shortcut to get sessions based on name
		sessionCookie := sessions.DefaultMany(c, "cookie")
		sessionRedis := sessions.DefaultMany(c, "redis")

		// check the session to see if the user is already signed-in. if yes: return, if no: continue
		// this is a single instance of THIS user's session (automatically isolated from all other sessions?)
		userSessionToken := sessionCookie.Get("userSessionToken")
		if userSessionToken != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message":        "User already signed-in. To sign-in with a different user, please signout.",
				"signed-in user": sessionRedis.Get("userSessionName"),
			})
			return
		}

		// Salt and hash the password
		// Check the database to make sure the passwords match using the username as reference. If no: return. If yes: continue
		password := json.Password
		hash, _ := helpers.HashPassword(password) // ignore error for the sake of simplicity

		fmt.Println("Password:", password)
		fmt.Println("Hash:    ", hash)

		match := helpers.CheckPasswordHash(password, hash)
		fmt.Println("Match:   ", match)

		// populate the session data on Redis with details retrieved from DB (opaque token, username, permissions, localization, etc.)
		sessionRedis.Set("userSessionToken", "qwerty123")
		sessionRedis.Set("userSessionName", "Jonathan Dudzik")
		sessionRedis.Set("userSessionEmail", "jonD@gmail.com")
		sessionRedis.Save()

		// send the opaque referenceID as a cookie to the client
		sessionCookie.Set("userSessionToken", sessionRedis.Get("userSessionToken"))
		sessionCookie.Save()

		c.JSON(http.StatusOK, gin.H{
			"content": "Successfully authenticated and signed in",
		})
	}
}

func SignOutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionCookie := sessions.DefaultMany(c, "cookie")
		sessionRedis := sessions.DefaultMany(c, "redis")

		if userSessionToken := sessionCookie.Get("userSessionToken"); userSessionToken == nil {
			log.Println("Invalid session token")
			return
		}

		log.Println("logging out user:", sessionRedis.Get("userSessionName"))

		sessionRedis.Clear()
		sessionRedis.Save()
		sessionCookie.Clear()
		sessionCookie.Save()
	}
}
