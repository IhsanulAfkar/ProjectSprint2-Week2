package controllers

import (
	"Week2/db"
	"Week2/forms"
	"Week2/helper"
	"Week2/helper/hash"
	"Week2/helper/jwt"
	"Week2/models"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type StaffController struct{}

func (h StaffController) SignUp(c *gin.Context) {
	
	var staffForm forms.StaffRegister
	if err := c.ShouldBindJSON(&staffForm); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
    }
	// check phone format
	if !helper.IsPhoneNumber(staffForm.PhoneNumber) || helper.ContainSpaces(staffForm.PhoneNumber) {
		c.JSON(400,gin.H{
			"message":"incorrect phone number"})
		return
	}
	conn := db.CreateConn()
	var staff models.Staff
	err := conn.QueryRowx("SELECT * FROM staff WHERE phone = $1 LIMIT 1", staffForm.PhoneNumber).StructScan(&staff)
	if err != nil && err != sql.ErrNoRows {

		c.JSON(500, gin.H{"message": "server error"})
		return
	}
	// check username
	if len(staffForm.Name) < 5 || len(staffForm.Name) > 50 {
	
		c.JSON(400, gin.H{
			"message":"name cannot below 5 nor exceed 15 characters"})
		return
	}
	// check password
	if len(staffForm.Password) < 5 || len(staffForm.Password) >15 {
		
		c.JSON(400, gin.H{
			"message":"password cannot below 5 nor exceed 15 characters"})
		return
	}
	hashedPass, err:= hash.HashPassword(staffForm.Password)
	if err != nil{
		c.JSON(500, gin.H{"Message":err.Error()})
		return
	}
	
	err = conn.QueryRowx("INSERT INTO staff (phone, name, password) VALUES ($1,$2,$3) RETURNING *",staffForm.PhoneNumber, staffForm.Name, hashedPass).StructScan(&staff)
	if err != nil {
		c.JSON(500, gin.H{"message":"failed to create user", "error": err.Error()})
		return
	}
	
	// create access token
	accessToken := jwt.SignJWT(staff)

	data := map[string]string{
		"userId":staff.Id,
		"phoneNumber":staff.Phone,
		"name":staff.Name,
		"accessToken":accessToken}

	c.JSON(201, gin.H{
		"message":"user created successfully",
		"data": data});
}


func (h StaffController)SignIn(c *gin.Context){
	conn:=db.CreateConn()
	var loginForm forms.StaffLogin
	if err:= c.ShouldBindJSON(&loginForm); err != nil{
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
	}

	if loginForm.PhoneNumber == "" || loginForm.Password == ""{
		c.JSON(400, gin.H{"message":"bad request"})
		return
	}
	if !helper.IsPhoneNumber(loginForm.PhoneNumber) || helper.ContainSpaces(loginForm.PhoneNumber) {
		c.JSON(400, gin.H{"message":"invalid phone format"})
		return
	}
	if len(loginForm.Password) < 5 || len(loginForm.Password) >15 {
		c.JSON(400, gin.H{
			"message":"password cannot below 5 nor exceed 15 characters"})
		return
	}
	var staff models.Staff
	err := conn.QueryRowx("SELECT * FROM staff WHERE phone = $1 LIMIT 1", loginForm.PhoneNumber).StructScan(&staff)

	if err != nil {
		if err == sql.ErrNoRows{
			c.JSON(404, gin.H{"message": "staff not found"})
			} else {
			c.JSON(500, gin.H{"message": "server error"})

		}
		return
	}
	if !hash.CheckPasswordHash(loginForm.Password, staff.Password) {
		c.JSON(400, gin.H{"message": "invalid password"})
		return
	}
	accessToken := jwt.SignJWT(staff)

	data := map[string]string{
		"phoneNumber":staff.Phone,
		"name":staff.Name,
		"accessToken":accessToken}

	c.JSON(200, gin.H{
		"message":"Staff logged successfully",
		"data": data});
}