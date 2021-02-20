package main

import (
	"context"
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go"
	jwtMiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

/*
Go-Firebase Template
- HTTP REST JSON API --> DONE
- Parse and validate Google Firebase Auth Token --> DONE
- create valid routes that work with the Auth service --> DONE
- CRUD (Create, read, update, delete) operations on a SQL database --> TODO
- Get claims from auth token --> TODO
- Set claims on user --> TODO
*/

var (
	app        *firebase.App
	// profileSrv profiles.ProfileService
)

const (
	sqlConnString = "postgresql://alpha:alphakilo0001@free-tier.gcp-us-central1.cockroachlabs.cloud:26257/corny-baboon-155.dev?sslmode=verify-full&sslrootcert=secrets/hackthenorth-cockroachdb.crt"
)

func main() {
	fmt.Println("Starting Server")
	r := gin.Default()
	r.Use(corsMiddleware)

	opt := option.WithCredentialsFile("secrets/vue-firebase-key.json") //TO DO UPDATE KEY
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	authMiddleware := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			ctx := context.Background()
			idToken, _ := jwtMiddleware.FromAuthHeader(c.Request)

			// fmt.Println(c.Request, idToken)

			client, err := app.Auth(ctx)
			if err != nil {
				log.Printf("error getting Auth client: %v\n", err)
				c.JSON(401, err)
				return
			}

			token, err := client.VerifyIDToken(ctx, idToken)
			if err != nil {
				log.Printf("error verifying ID token: %v\n", err)
				c.JSON(401, err)
				return
			}

			// log.Printf("Verified ID token: %v\n", token)
			c.Set("token", token)
			c.Set("UID", token.UID)
			// log.Printf("UID %v", token.UID)
		}
	}

	// profileSrv = profiles.NewProfileService(sqlConnString)


	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/auth", authMiddleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": c.MustGet("token"),
		})
	})

	// r.POST("/endpoint", authMiddleware(), endpointDandler)

	r.Run(":8081") // listen and serve on 0.0.0.0:8081 (for windows "localhost:8081")
}

var corsMiddleware = cors.New(cors.Config{
	// AllowOrigins:     []string{"https://wheypal.com", "http://localhost:8080"},
	AllowOrigins: []string{"*"},
	AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
	// AllowMethods:     []string{"*"},
	AllowHeaders:     []string{"Authorization", "Origin", "Content-Length", "Content-Type"},
	AllowCredentials: true,
	MaxAge:           12 * time.Hour,
})
