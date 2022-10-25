package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	//"html/template"
	//"strings"

	globals "github.com/jonathandudzik/gin_session_auth/globals"
	middleware "github.com/jonathandudzik/gin_session_auth/middleware"
	routes "github.com/jonathandudzik/gin_session_auth/routes"
)

func main() {
	router := gin.Default()

	// router.Static("/assets", "./assets")
	// router.LoadHTMLGlob("templates/*.html")

	router.Use(sessions.Sessions("session", cookie.NewStore(globals.Secret)))

	public := router.Group("/")
	routes.PublicRoutes(public)

	private := router.Group("/")
	private.Use(middleware.AuthRequired)
	routes.PrivateRoutes(private)

	router.Run("localhost:8080")
}

// func checkResourceGroup(c *gin.Context) {
// 	// router.GET("/check-resource-group/:resourceGroupName", checkResourceGroup)
// 	resourceGroupName := c.Param("resourceGroupName")

// 	// Authentication for local development
// 	cred, err := azidentity.NewDefaultAzureCredential(nil)
// 	if err != nil {
// 		log.Fatalf("failed to obtain a credential: %v", err)
// 	}

// 	exits, err := checkExistenceResourceGroup(ctx, cred, resourceGroupName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println("resources group already exist:", exits)

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": exits,
// 	})
// }

// func checkExistenceResourceGroup(ctx context.Context, cred azcore.TokenCredential, resourceGroupName string) (bool, error) {
// 	resourceGroupClient, err := armresources.NewResourceGroupsClient(subscriptionID, cred, nil)
// 	if err != nil {
// 		return false, err
// 	}

// 	boolResp, err := resourceGroupClient.CheckExistence(ctx, resourceGroupName, nil)
// 	if err != nil {
// 		return false, err
// 	}

// 	return boolResp.Success, nil
// }
