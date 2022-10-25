package controllers

import (
	"github.com/gin-contrib/sessions"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	globals "github.com/jonathandudzik/gin_session_auth/globals"
)

// Other actions (routes) in the app will have the proper authorization checks
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

func SignInHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		json := UserSignInDetails{}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Salt and hash and check the DB against the user credentials

		session := sessions.Default(c)
		if user := session.Get(globals.Userkey); user != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"content": "User already logged in.",
				"body":    json.Email,
			})
			return
		}

		session.Set(globals.Userkey, json.Email)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"content": "Failed to save session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"content": "Successfully authenticated and signed in",
			"body":    json.Email,
		})
	}
}

func SignOutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		log.Println("logging out user:", user)
		if user == nil {
			log.Println("Invalid session token")
			return
		}
		session.Delete(globals.Userkey)
		if err := session.Save(); err != nil {
			log.Println("Failed to save session:", err)
			return
		}

		c.Redirect(http.StatusMovedPermanently, "/")
	}
}
