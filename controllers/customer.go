package controllers

import (
	"Week2/db"
	"Week2/forms"
	"Week2/helper"
	"Week2/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type CustomerController struct{}

func (h CustomerController) CreateCustomer(c *gin.Context){
	var customerForm forms.CustomerRegister
	if err := c.ShouldBindJSON(&customerForm); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
    }
	if !helper.IsPhoneNumber(customerForm.PhoneNumber) || helper.ContainSpaces(customerForm.PhoneNumber){
		c.JSON(400, gin.H{"message":"incorrect phone number"})
		return
	}
	conn := db.CreateConn()
	var customer models.Customer
	// check phone number
	res, err := conn.Exec("SELECT id FROM customer WHERE phone = $1", customerForm.PhoneNumber)
	if err != nil {
		
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	if row, _:= res.RowsAffected(); row > 0 {
		c.JSON(400, gin.H{"message":"phone number already exist"})
		return
	}
	err = conn.QueryRowx("INSERT INTO customer (name, phone) VALUES ($1, $2) RETURNING *", customerForm.Name, customerForm.PhoneNumber).StructScan(&customer)
	if err != nil {
		
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	data := map[string]string{
		"userId":customer.Id,
		"phoneNumber":customer.Phone,
		"name":customer.Name,
	}
	c.JSON(201,gin.H{"message":"success create customer","data":data})
}

func (h CustomerController)GetAllCustomer(c *gin.Context){
	phoneNumber := c.Query("phoneNumber")
	name := strings.ToLower(c.Query("name"))

	var args []interface{}
	var queryParams []string
	argIdx := 1
	if name != ""{
		nameWildcard := "%" + name +"%"
		queryParams = append(queryParams, " name LIKE $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, nameWildcard)
		argIdx += 1
	}
	if phoneNumber!=""{
		phoneWildcard := "+" + phoneNumber + "%"
		queryParams = append(queryParams, " phone LIKE $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, phoneWildcard)
		argIdx += 1
	}
	sqlQuery := "SELECT * FROM customer WHERE \"deletedAt\" IS NULL"
	if len(queryParams)>0{
		allQuery := strings.Join(queryParams, " AND")
		sqlQuery += " AND " + allQuery 
	}
	sqlQuery += " ORDER BY \"createdAt\" DESC"
	conn:=db.CreateConn()
	customers := make([]models.Customer, 0)
	err := conn.Select(&customers, sqlQuery, args...)
	if err != nil {
		
		c.JSON(500,gin.H{"message":"server error"})
		return
	}
	c.JSON(200, gin.H{"message":"success", "data":customers})
}
