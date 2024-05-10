package middleware

import (
	"Week2/db"
	"Week2/helper/jwt"
	"Week2/models"
	"database/sql"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)
func getBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}
func AuthMiddleware(c *gin.Context) {
	token, err := getBearerToken(c.GetHeader("Authorization"))
	if err!= nil {
		c.AbortWithStatusJSON(401, gin.H{
			"message": err.Error()})
		return
	}
	id, err := jwt.ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"message":err.Error()})
		return
	}
	// find user
	conn := db.CreateConn()
	var staff models.Staff
	err = conn.QueryRowx("SELECT * FROM staff WHERE id = ? LIMIT 1",id).StructScan(&staff)
	if err != nil && err == sql.ErrNoRows{
		c.AbortWithStatusJSON(404, gin.H{
			"message":"staff not found"})
			return
	}
	c.Set("userId",id)
	c.Next()
}