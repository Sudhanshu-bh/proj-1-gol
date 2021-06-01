package main

import (
	"fmt"
	"net/http"

	"proj-backend/db"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	// "/apis/loginApi"
)

func main() {
	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./build", true)))
	router.Use(static.Serve("/login", static.LocalFile("./build", true)))
	router.Use(static.Serve("/dashboard", static.LocalFile("./build", true)))
	router.Use(static.Serve("/profile", static.LocalFile("./build", true)))
	router.Use(static.Serve("/changePassword", static.LocalFile("./build", true)))
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code": "PAGE_NOT_FOUND", "message": "Page not found", //c.File("./public/index.html")
		})
	})

	// Setup route group for the API
	api := router.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "test")
		})

		api.POST("/changePass", func(c *gin.Context) {
			var changePassCreds db.ChangePass

			if err := c.ShouldBindJSON(&changePassCreds); err == nil {

				ResultStruct, errFromCollection, errConnecting := db.ChangePassInDB(changePassCreds)
				if errConnecting != nil {
					c.JSON(500, gin.H{"error": errConnecting.Error()})
				} else if errFromCollection != nil {
					c.JSON(500, gin.H{"error": errFromCollection.Error()})
				} else if ResultStruct.MatchedCount == 0 {
					c.JSON(401, gin.H{"error": "Invalid current password"})
				} else if changePassCreds.NewPass == changePassCreds.CurrPass {
					c.JSON(401, gin.H{"error": "New password cannot be same as old password"})
				} else if ResultStruct.ModifiedCount == 1 {
					c.JSON(200, gin.H{"status": "Password change successful"})
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 400
			}
		})
	}

	validate := router.Group("/validate")
	{
		validate.POST("/credentials", func(c *gin.Context) {
			var credentials db.Details

			if err := c.ShouldBindJSON(&credentials); err == nil {

				errFromCollection, errConnecting := db.CheckUserInDB(credentials.User)
				if errConnecting != nil {
					c.JSON(500, gin.H{"error": errConnecting.Error()})
				} else if errFromCollection == mongo.ErrNoDocuments {
					c.JSON(401, gin.H{"error": "Invalid email"})
				} else if errFromCollection == nil && credentials.Password != db.DBResult.Password {
					c.JSON(401, gin.H{"error": "Invalid password"})
				} else if errFromCollection == nil && credentials.Password == db.DBResult.Password {
					c.JSON(200, gin.H{"status": "Success"})
				} else {
					c.JSON(500, gin.H{"error2": errFromCollection.Error()})
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 400
			}
		})

		validate.POST("/email", func(c *gin.Context) {
			var email db.Email

			if err := c.ShouldBindJSON(&email); err == nil {

				errFromCollection, errConnecting := db.CheckUserInDB(email.User)
				if errConnecting != nil {
					c.JSON(500, gin.H{"error": errConnecting.Error()})
				} else if errFromCollection == mongo.ErrNoDocuments {
					c.JSON(401, gin.H{"error": "email not registered"})
				} else if errFromCollection == nil {
					c.JSON(200, gin.H{"status": "email is valid"})
				} else {
					c.JSON(500, gin.H{"error2": errFromCollection.Error()})
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 400
			}
		})

		router.POST("/fetchProfile", func(c *gin.Context) {
			var email db.Email

			if err := c.ShouldBindJSON(&email); err == nil {

				errFromCollection, errConnecting := db.FetchProfileFromDB(email.User)
				if errConnecting != nil {
					c.JSON(500, gin.H{"error": errConnecting.Error()})
				} else if errFromCollection == mongo.ErrNoDocuments {
					c.JSON(401, gin.H{"error": "email not found"})
				} else if errFromCollection == nil {
					fmt.Println(db.DBProfileResult)
					c.JSON(200, db.DBProfileResult)
				} else {
					c.JSON(500, gin.H{"error2": errFromCollection.Error()})
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 400
			}
		})

		validate.POST("/currPass", func(c *gin.Context) {
			var credentials db.Details

			if err := c.ShouldBindJSON(&credentials); err == nil {

			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 400
			}
		})
	}

	// Start and run the server
	router.Run()

}
