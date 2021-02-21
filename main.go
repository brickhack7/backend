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

var (
	app            *firebase.App
	repo           *Repo
	corsMiddleware = cors.New(cors.Config{
		// AllowOrigins:     []string{"https://wheypal.com", "http://localhost:8080"},
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		// AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Authorization", "Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
)

const (
	sqlConnString = "postgres://allen:jzx5sWtgykoGDe8Y@free-tier.gcp-us-central1.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full&sslrootcert=secrets/cc-ca.crt&options=--cluster=nimble-vole-835"
)

func main() {
	fmt.Println("Starting Server")
	r := gin.Default()
	r.Use(corsMiddleware)

	opt := option.WithCredentialsFile("secrets/firebase-key.json") //TO DO UPDATE KEY
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	repo = NewRepo(sqlConnString)

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

	r.GET("/profile", authMiddleware(), getProfile)
	r.PUT("/profile", authMiddleware(), updateProfile)
	r.GET("/discover", authMiddleware(), getLoc)
	r.POST("/discover", authMiddleware(), postLoc)
	r.DELETE("/match", authMiddleware(), deleteMatch)
	r.GET("/match", authMiddleware(), getMatches)

	// r.POST("/endpoint", authMiddleware(), endpointDandler)

	r.Run(":8081") // listen and serve on 0.0.0.0:8081 (for windows "localhost:8081")
}
