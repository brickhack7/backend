package main

import (
	"context"

	"github.com/gin-gonic/gin"
)

func getProfile(c *gin.Context) {
	var prof *profile
	err := c.Bind(&prof)
	if err != nil {
		c.JSON(501, err)
		return
	}
	_, err = repo.conn.Exec(context.Background(), selectProfileByID, &prof.id, &prof.lat, &prof.long, &prof.interests)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, &prof)
}

func updateProfile(c *gin.Context) {
	var prof *profile
	err := c.Bind(&prof)
	if err != nil {
		c.JSON(501, err)
		return
	}

	var oldProfile *profile
	_, err = repo.conn.Exec(context.Background(), selectProfileByID, &oldProfile.id, &oldProfile.lat, &oldProfile.long, &oldProfile.interests)
	if err != nil {
		c.JSON(500, err)
		return
	}

	if prof.id == "" {
		prof.id = oldProfile.id
	}
	if prof.lat == "" {
		prof.lat = oldProfile.lat
	}
	if prof.long == "" {
		prof.long = oldProfile.long
	}
	if prof.interests == nil {
		prof.interests = oldProfile.interests
	}

	_, err = repo.conn.Exec(context.Background(), updateProfilebyID, &prof.id, &prof.lat, &prof.long, &prof.interests)
	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, &prof)
}

func getLoc(c *gin.Context) {
	var resp *int
	c.JSON(200, &resp)
}

func postLoc(c *gin.Context) {
	var resp *int
	c.JSON(200, &resp)
}

func deleteMatch(c *gin.Context) {
	var resp *int
	c.JSON(200, &resp)
}

func getMatches(c *gin.Context) {
	var resp *int
	c.JSON(200, &resp)
}

const (
	selectProfileByID = "SELECT uid, lat, long, interests FROM profiles WHERE uid $1"
	updateProfilebyID = "UPDATE profiles SET (lat, long, interests) WHERE uid $1"
)
