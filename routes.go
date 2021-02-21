package main

import (
	"context"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getProfile(c *gin.Context) {
	var prof *Profile
	err := c.Bind(&prof)
	if err != nil {
		c.JSON(501, err)
		return
	}
	err = repo.conn.QueryRow(context.Background(), selectProfileByID, &prof.ID).Scan(&prof.Lat, &prof.Long, &prof.Interests)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, &prof)
}

func updateProfile(c *gin.Context) {
	var prof *Profile
	err := c.Bind(&prof)
	if err != nil {
		c.JSON(501, err)
		return
	}

	var oldProfile *Profile
	_, err = repo.conn.Exec(context.Background(), selectProfileByID, &oldProfile.ID, &oldProfile.Lat, &oldProfile.Long, &oldProfile.Interests)
	if err != nil {
		c.JSON(500, err)
		return
	}

	if prof.ID == "" {
		prof.ID = oldProfile.ID
	}
	if prof.Lat == 0 {
		prof.Lat = oldProfile.Lat
	}
	if prof.Long == 0 {
		prof.Long = oldProfile.Long
	}
	if prof.Interests == nil {
		prof.Interests = oldProfile.Interests
	}

	_, err = repo.conn.Exec(context.Background(), updateProfilebyID, &prof.ID, &prof.Lat, &prof.Long, &prof.Interests)
	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, &prof)
}

func getLoc(c *gin.Context) {
	type locReq struct {
		uid  string  `json:'lat'`
		lat  float64 `json:'lat'`
		long float64 `json:'long'`
	}

	req := &locReq{}
	req.lat, _ = strconv.ParseFloat(c.Query("lat"), 64)
	req.long, _ = strconv.ParseFloat(c.Query("long"), 64)
	req.uid = c.GetString("uid")

	rows, err := repo.conn.Query(context.Background(), getLocsByGeo, &req.uid, &req.lat, &req.long)
	if err != nil {
		c.JSON(500, err)
		return
	}

	locs := make([]*Location, 0)

	for rows.Next() {
		loc := &Location{}
		err = rows.Scan(&loc.ID, &loc.Name, &loc.City, &loc.Distance)
		if err != nil {
			log.Printf("[Error] [getLoc] | %v", err)
			c.JSON(501, err)
			return
		}
		locs = append(locs, loc)
	}

	c.JSON(200, &locs)
}

func postLoc(c *gin.Context) {
	type locResponseReq struct {
		uid    string `json:'uid'`
		loc_id string `json:'loc_id'`
	}

	var req *locResponseReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(501, err)
		return
	}

	_, err = repo.conn.Exec(context.Background(), locMatch, &req.uid, &req.loc_id, 1)
	if err != nil {
		c.JSON(500, err)
		return
	}

	var matchID string
	err = repo.conn.QueryRow(context.Background(), bestMatch, &req.uid).Scan(matchID)
	if err != nil {
		c.JSON(500, err)
		return
	}

	var prof *Profile
	err = repo.conn.QueryRow(context.Background(), selectProfileByID, &prof.ID).Scan(&prof.Lat, &prof.Long, &prof.Interests)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, &prof)
}

func deleteMatch(c *gin.Context) {
	var targetID *string
	_, err := repo.conn.Exec(context.Background(), deleteMatchByID, c.GetString("uid"), &targetID)
	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, true)
}

func getMatches(c *gin.Context) {
	type match struct {
		uid1 string `json:'uid1'`
		uid2 string `json:'uid2'`
	}

	rows, err := repo.conn.Query(context.Background(), getMatchesByID, c.GetString("uid"))
	if err != nil {
		c.JSON(500, err)
		return
	}

	matches := make([]*match, 0)

	for rows.Next() {
		match := &match{}
		err = rows.Scan(&match.uid1, &match.uid2)
		if err != nil {
			c.JSON(501, err)
			return
		}
		matches = append(matches, match)
	}

	c.JSON(200, &matches)
}

const (
	selectProfileByID = "SELECT uid, lat, long, interests FROM profiles WHERE uid $1"
	updateProfilebyID = "UPDATE profiles SET (lat, long, interests) WHERE uid $1"
	getLocsByGeo      = "WITH cte AS ( " +
		"SELECT loc_id, city, name, " +
		"( " +
		"6371 * " +
		"acos(cos(radians($2)) * " +
		"cos(radians(lat)) *  " +
		"cos(radians(long) -  " +
		"radians($3)) +  " +
		"sin(radians($2)) *  " +
		"sin(radians(lat ))) " +
		") AS distance  " +
		"FROM locations ) " +
		"SELECT * " +
		"FROM cte " +
		"WHERE cte.distance < 50 and cte.loc_id <> ALL(SELECT loc_id FROM loc_matches WHERE uid = $1) " +
		"ORDER BY cte.distance LIMIT 20;"
	locMatch  = "INSERT INTO loc_matches (uid, loc_id, status) VALUES ($1,$2,$3);"
	bestMatch = "with cte as ( " +
		"SELECT uid, count(1) as common " +
		"FROM ( " +
		"SELECT uid, unnest(profiles.interests) as intr " +
		"FROM profiles WHERE profiles.uid = any(SELECT uid from loc_matches where loc_id = $1 and uid <> $1 and status = $1) " +
		"and profiles.uid <> ALL(SELECT uid1 from user_matches where uid2 = $1) " +
		"and profiles.uid <> ALL(SELECT uid2 from user_matches where uid1 = $1) " +
		") " +
		"x WHERE intr = any( " +
		"SELECT unnest(interests) FROM profiles WHERE uid = $1) AND uid <> $1 " +
		"GROUP BY uid)  " +
		"SELECT cte.uid FROM cte INNER JOIN profiles p on cte.uid = p.uid ORDER BY common DESC LIMIT 1;"
	deleteMatchByID = "DELETE FROM user_matches WHERE (uid1 = $1 AND uid2 = $2) or (uid1 = $2 AND uid2 = $1)"
	getMatchesByID  = "SELECT uid1, uid2 FROM user_matches WHERE uid1 = $1 or uid2 = $1;"
)
