// Public and Private routes are defined here and controllers are attached
package routes

import (
	"github.com/gin-gonic/gin"

	controllers "github.com/jonathandudzik/gin_session_auth/controllers"
)

func PublicRoutes(g *gin.RouterGroup) {

	g.POST("/sign-in", controllers.SignInHandler())
	// g.POST("/sign-up", controllers.SignUpHandler())

}

func PrivateRoutes(g *gin.RouterGroup) {

	g.GET("/sign-out", controllers.SignOutHandler())

}
