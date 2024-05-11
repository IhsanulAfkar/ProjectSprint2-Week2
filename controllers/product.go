package controllers

import (
	"Week2/db"
	"Week2/forms"
	"Week2/helper"
	"Week2/models"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductController struct{}

func (h ProductController) CreateProduct(c *gin.Context) {
	var productForm forms.ProductCreate
	if err := c.ShouldBindJSON(&productForm);err !=nil {
		c.JSON(400, gin.H{"message":err.Error()})
		return
	}
	if productForm.IsAvailable == nil {
		c.JSON(400, gin.H{"message":"bad request"})
		return
	}
	// if productForm[""]
	if productForm.Name == "" || productForm.Sku == "" || productForm.Notes == ""||productForm.Location == ""{
		c.JSON(400, gin.H{"message":"bad request, cannot include empty string"})
		return
	}
	if len(productForm.Name) > 30{
		c.JSON(400, gin.H{"message":"bad request"})
		return
	}
	if len(productForm.Sku) > 30{
		c.JSON(400, gin.H{"message":"bad request"})
		return
	}
	if len(productForm.Notes) > 200{
		c.JSON(400, gin.H{"message":"bad request"})
		return
	}
	if len(productForm.Location) > 200{
		c.JSON(400, gin.H{"message":"bad request"})
		return
	}
	if productForm.Stock < 1 || productForm.Stock > 100000{
		c.JSON(400, gin.H{"message":"stock invalid"})
		return
	}
	if productForm.Price < 1{
		c.JSON(400, gin.H{"message":"price invalid"})
		return
	}
	if !helper.IsURL(productForm.ImageUrl) {
		c.JSON(400, gin.H{"message":"not a valid url"})
		return
	}
	if !helper.Includes(productForm.Category, models.Category[:]){
		c.JSON(400, gin.H{"message":"invalid category"})
		return
	}
	conn := db.CreateConn()
	var newId string
	var createdAt string
	query := "INSERT INTO product (name, sku, category, \"imageUrl\", notes, price, stock, location, \"isAvailable\") VALUES ($1,$2,$3,$4,$5,$6,$7,$8, $9) RETURNING id, \"createdAt\""
	err := conn.QueryRow(query, productForm.Name, productForm.Sku, productForm.Category, productForm.ImageUrl, productForm.Notes, productForm.Price, productForm.Stock, productForm.Location, productForm.IsAvailable).Scan(&newId,&createdAt)
	if err != nil {
		
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	data:= map[string]string{
		"id":newId,
		"createdAt": helper.FormatToIso860(createdAt),
	}
	c.JSON(201, gin.H{"message":"success","data":data})
}
func (h ProductController) GetAllProduct(c *gin.Context){
	limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if errLimit != nil || limit < 0 {
		limit = 5
	}
	offset, errOffset := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if errOffset != nil || offset < 0 {
		offset = 0
	}
	id := c.Query("id")
	// if !helper.IsUUID(id) {
	// 	id = ""
	// }
	name := strings.ToLower(c.Query("name"))
	sku := c.Query("sku")
	category := c.Query("category")
	inStock := c.Query("inStock")
	isAvailable := c.Query("isAvailable")
	createdAt := c.Query("createdAt")
	price := c.Query("price")

	if len(name) > 30 {
		name = ""
	}
	if len(sku) > 30 {
		sku = ""
	}
	// if name
	if isAvailable != "true" && isAvailable !="false"{
		isAvailable = ""
	}
	if inStock != "true" && inStock != "false"{
		inStock = ""
	}
	if createdAt != "asc" && createdAt != "desc"{
		createdAt = ""
	}
	if price != "asc" && price != "desc"{
		price = ""
	}

	// validate query
	if !helper.Includes(category, models.Category[:]){
		category = ""
	}
	// Build the SQL query string
	sqlQuery := "SELECT * FROM product WHERE \"deletedAt\" IS NULL"

	// Prepare parameters slice
	var args []interface{}
	// queryParams := make(map[string]interface{})
	var queryParams []string
	// Add conditions based on parameters
	// WHERE
	argIdx := 1
	if id != "" {
		
		queryParams = append(queryParams, " id::text = $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, id)
		argIdx += 1
	}
	if name != ""{
		nameWildcard := "%" + name +"%"
		queryParams = append(queryParams, " name LIKE $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, nameWildcard)
		argIdx += 1
	}
	if sku != ""{
		
		queryParams = append(queryParams, " sku = $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, sku)
		argIdx += 1
	}
	if category != ""{
		
		queryParams = append(queryParams, " category = $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, category)
		argIdx += 1
	}
	if inStock != ""{
		
		if inStock =="true"{
			queryParams = append(queryParams, " stock > 0 ") 
		} else {
			queryParams = append(queryParams, " stock = 0 ") 
		}
	}
	if isAvailable != "" {
		
		if isAvailable == "true"{
			queryParams = append(queryParams, " \"isAvailable\" = true")
		} else {
			queryParams = append(queryParams, " \"isAvailable\" = false")
		}
	}
	
	
	if len(queryParams) > 0 {
		allQuery := strings.Join(queryParams, " AND")
		sqlQuery += " AND " + allQuery 
	}
	orderBy := true
	var orderQuery []string
	if createdAt == ""{
		orderQuery = append(orderQuery, " \"createdAt\" DESC")
	}else{
		if createdAt == "asc"{
			orderQuery = append(orderQuery, " \"createdAt\" ASC")
		}
	}
	if price != ""{
		orderBy = true
		if price =="asc" {
			orderQuery = append(orderQuery, " price ASC")
		} else {
			orderQuery = append(orderQuery, " price DESC")
		}
	}
	if orderBy {
		sqlQuery += " ORDER BY " + strings.Join(orderQuery, ",")
	}
	sqlQuery += " LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)
	
	conn := db.CreateConn()
	
	products := make([]models.Product, 0)
	err := conn.Select(&products, sqlQuery, args...)
	if err != nil {
		
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	
	c.JSON(200, gin.H{"message": "success","data": products})
}

func (h ProductController)UpdateProduct(c *gin.Context){
	productId := c.Param("productId")
	if !helper.IsUUID(productId){
		c.JSON(404, gin.H{"message":"invalid id"})
		return
	}
	var productForm forms.ProductCreate
	if err := c.ShouldBindJSON(&productForm);err !=nil {
		c.JSON(400, gin.H{"message":err.Error()})
		return
	}
	if productForm.IsAvailable == nil {
		c.JSON(400, gin.H{"message":"bad request"})
		return
	}
	if productForm.Name == "" || productForm.Sku == "" || productForm.Notes == ""||productForm.Location == ""{
		c.JSON(400, gin.H{"message":"bad request, cannot include empty string"})
		return
	}
	if len(productForm.Name) > 30{
		c.JSON(400, gin.H{"message":"bad request"})
		return
	}
	if len(productForm.Sku) > 30{
		c.JSON(400, gin.H{"message":"bad request"})
		return
	}
	if len(productForm.Notes) > 200{
		c.JSON(400, gin.H{"message":"bad request"})
		return
	}
	if len(productForm.Location) > 200{
		c.JSON(400, gin.H{"message":"bad request"})
		return
	}
	
	if productForm.Stock < 1 || productForm.Stock > 100000{
		c.JSON(400, gin.H{"message":"stock invalid"})
		return
	}
	if productForm.Price < 1{
		c.JSON(400, gin.H{"message":"price invalid"})
		return
	}
	if !helper.IsURL(productForm.ImageUrl) {
		c.JSON(400, gin.H{"message":"not a valid url"})
		return
	}
	if !helper.Includes(productForm.Category, models.Category[:]){
		c.JSON(400, gin.H{"message":"invalid category"})
		return
	}
	conn := db.CreateConn()
	var product models.Product
	err := conn.QueryRowx("SELECT * FROM product WHERE id = $1 AND \"deletedAt\" IS NULL", productId).StructScan(&product)
	if err != nil {
		if err == sql.ErrNoRows{
			c.JSON(404, gin.H{"message":"no product found"})
		} else {
			c.JSON(500, gin.H{"message":"server error"})
		}
		return
	}
	current := time.Now()
	query := "UPDATE product SET name = $1, sku = $2, category = $3, \"imageUrl\" = $4, notes = $5, price = $6, stock = $7, location = $8, \"isAvailable\" = $9, \"updatedAt\" = $10 WHERE id = $11"
	res, err := conn.Exec(query, productForm.Name, productForm.Sku, productForm.Category, productForm.ImageUrl, productForm.Notes, productForm.Price, productForm.Stock, productForm.Location, productForm.IsAvailable, current, productId)
	if err != nil {
		
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	if rows, _ := res.RowsAffected(); rows == 0{
		
		c.JSON(500, gin.H{"message":"no product updated"})
		return
	}
	c.JSON(200, gin.H{"message":"success update product"})

}

func (h ProductController)DeleteProduct(c *gin.Context){
	productId := c.Param("productId")
	if !helper.IsUUID(productId){
		c.JSON(404, gin.H{"message":"bad request"})
		return
	}
	// check if exist
	current := time.Now()
	conn:= db.CreateConn()
	query := "UPDATE product SET \"deletedAt\" = $1 WHERE id = $2 AND \"deletedAt\" IS NULL"
	res, err := conn.Exec(query, current, productId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"message":"no product found"})
		} else{
			
			c.JSON(404, gin.H{"message":"server error"})
		}
		return
	}
	if row, _ :=res.RowsAffected(); row == 0 {
		c.JSON(404, gin.H{"message":"no product found"})
		return
	}
	c.JSON(200, gin.H{"message":"product deleted successfully"})
}

func (h ProductController)SearchBySKU(c *gin.Context){
	limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if errLimit != nil || limit < 0 {
		limit = 5
	}
	offset, errOffset := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if errOffset != nil || offset < 0 {
		offset = 0
	}
	name := strings.ToLower(c.Query("name"))
	sku := c.Query("sku")
	category := c.Query("category")
	inStock := c.Query("inStock")
	price := c.Query("price")

	if len(name) > 30 {
		name = ""
	}
	if len(sku) > 30 {
		sku = ""
	}
	if inStock != "true" && inStock != "false"{
		inStock = ""
	}
	if price != "asc" && price != "desc"{
		price = ""
	}

	// validate query
	if !helper.Includes(category, models.Category[:]){
		category = ""
	}
	// Build the SQL query string
	sqlQuery := "SELECT * FROM product WHERE \"deletedAt\" IS NULL AND \"isAvailable\" = true "

	// Prepare parameters slice
	var args []interface{}
	// queryParams := make(map[string]interface{})
	var queryParams []string
	// Add conditions based on parameters
	// WHERE
	argIdx := 1
	if name != ""{
		nameWildcard := "%" + name +"%"
		queryParams = append(queryParams, " name LIKE $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, nameWildcard)
		argIdx += 1
	}
	if sku != ""{
		
		queryParams = append(queryParams, " sku = $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, sku)
		argIdx += 1
	}
	if category != ""{
		
		queryParams = append(queryParams, " category = $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, category)
		argIdx += 1
	}
	if inStock != ""{
		
		if inStock =="true"{
			queryParams = append(queryParams, " stock > 0 ") 
		} else {
			queryParams = append(queryParams, " stock = 0 ") 
		}
	}
	if len(queryParams) > 0 {
		allQuery := strings.Join(queryParams, " AND")
		sqlQuery += " AND " + allQuery 
	}
	// orderBy := false
	// var orderQuery []string
	sqlQuery += " ORDER BY "
	if price != ""{
		// orderBy = true
		if price =="asc" {
			// orderQuery = append(orderQuery, " price ASC")
			sqlQuery += " price ASC"
		} else {
			sqlQuery += " price DESC"
			// orderQuery = append(orderQuery, " price DESC")
		}
	} else {
		sqlQuery += "\"createdAt\" DESC"
	}
	
	
	// if orderBy {
	// 	sqlQuery +=  ", " +strings.Join(orderQuery, ",")
	// }
	sqlQuery += " LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)
	conn := db.CreateConn()
	
	products := make([]models.Product, 0)
	err := conn.Select(&products, sqlQuery, args...)
	if err != nil {
		
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	
	c.JSON(200, gin.H{"message": "success","data": products})
}

func (h ProductController)Checkout(c *gin.Context){
	var checkoutForm forms.Checkout
	if err := c.ShouldBindJSON(&checkoutForm);err !=nil {
		c.JSON(400, gin.H{"message":err.Error()})
		return
	}
	// validating !!
	if !helper.IsUUID(checkoutForm.CustomerId) {
		c.JSON(400, gin.H{"message":"invalid customer id"})
		return
	}
	if checkoutForm.Paid <= 0{
		c.JSON(400, gin.H{"message":"paid cannot below 0"})
		return
	}
	if checkoutForm.Change == nil{
		c.JSON(400, gin.H{"message":"invalid change"})
		return
	}
	if *checkoutForm.Change < 0{
		c.JSON(400, gin.H{"message":"change cannot below 0"})
		return
	}
	if len(checkoutForm.ProductDetails) == 0 {
		c.JSON(400, gin.H{"message":"product cannot empty"})
		return
	}
	// check if product list is unique
	if !helper.IsProductsUnique(checkoutForm.ProductDetails){
		c.JSON(400, gin.H{"message":"duplicate product list"})
		return
	}
	// find customer
	var customer models.Customer
	conn := db.CreateConn()
	query := "SELECT * FROM customer WHERE \"deletedAt\" IS NULL AND id = $1"
	err := conn.QueryRowx(query, checkoutForm.CustomerId).StructScan(&customer)
	if err != nil {
		if err == sql.ErrNoRows{
			c.JSON(404, gin.H{"message":"customer not found"})
		} else {
			
			c.JSON(500, gin.H{"message":"server error"})
		}
		return
	}
	totalPrice := 0
	var products []models.ProductCheckout
	// check each product
	for _, productDetail := range checkoutForm.ProductDetails {
		if productDetail.Quantity <1 {
			c.JSON(400, gin.H{"message":"invalid quantity"})
			return
		}
		// validate
		if productDetail.ProductId == ""{
			c.JSON(400, gin.H{"message":"productId cannot be empty"})
			return
		}
		
		// if !helper.IsUUID(productDetail.ProductId){
		// 	c.JSON(400, gin.H{"message":"invalid product id"})
		// 	return
		// }
		// find product
		var product models.ProductCheckout
		query = "SELECT * FROM product WHERE \"deletedAt\" IS NULL AND id::text = $1 AND \"isAvailable\" = true"
		err = conn.QueryRowx(query, productDetail.ProductId).StructScan(&product)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"message":"product " + productDetail.ProductId + " not found"})
			} else {
				
				c.JSON(500, gin.H{"message":"server error"})
			}
			return
		}
		
		// check quantity 
		if productDetail.Quantity > product.Stock {
			c.JSON(400, gin.H{"message":"product " + productDetail.ProductId+" insufficient stock"})
			return
		}
		product.Stock = product.Stock - productDetail.Quantity
		product.Quantity = productDetail.Quantity
		// sum the price
		totalPrice += productDetail.Quantity * product.Price
		
		products = append(products, product)
	}
	// check if paid and change is valid
	if checkoutForm.Paid - *checkoutForm.Change != totalPrice {
		c.JSON(400, gin.H{"message":"incorrect paid and change value"})
		return
	}
	// create transactions 
	tx, err := conn.Beginx()

	if err != nil {		
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	// create checkout
	query = "INSERT INTO checkout (\"customerId\", total, paid, change) VALUES ($1,$2,$3,$4) RETURNING id"
	var checkoutId string
	err = tx.QueryRowx(query, checkoutForm.CustomerId, totalPrice, checkoutForm.Paid, checkoutForm.Change).Scan(&checkoutId)
	if err != nil {
		tx.Rollback()
		
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	// if rows, _:= res.RowsAffected(); rows ==0{
	// 	tx.Rollback()
	// 	c.JSON(500, gin.H{"message":"server error"})
	// 	return
	// }
	// update all quantity from products and create checkout item
	for _, product := range products {
		
		query = "UPDATE product SET stock = $1 WHERE id = $2"
		res, err := tx.Exec(query, product.Stock, product.Id)
		if err != nil {
			tx.Rollback()
			
			c.JSON(500, gin.H{"message":"server error"})
			return
		}
		if rows, _:= res.RowsAffected(); rows ==0{
			tx.Rollback()
			c.JSON(500, gin.H{"message":"server error"})
			return
		}
		query = "INSERT INTO \"checkoutItem\" (\"checkoutId\", \"productId\", price, quantity) VALUES ($1,$2,$3,$4)"
		res, err = tx.Exec(query, checkoutId, product.Id, product.Price, product.Quantity)
		if err != nil {
			tx.Rollback()
			
			c.JSON(500, gin.H{"message":"server error"})
			return
		}
		if rows, _:= res.RowsAffected(); rows ==0{
			tx.Rollback()
			c.JSON(500, gin.H{"message":"server error"})
			return
		}
	}
	// well done!
	err = tx.Commit()
	if err != nil {
		
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	c.JSON(200, gin.H{"message":"success checkout!"})
}

func (h ProductController)GetAllCheckout(c *gin.Context){
	limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if errLimit != nil || limit < 0 {
		limit = 5
	}
	offset, errOffset := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if errOffset != nil || offset < 0 {
		offset = 0
	}
	customerId := c.Query("customerId")
	createdAt := c.Query("createdAt")
	if createdAt != "asc" && createdAt != "desc"{
		createdAt = ""
	}
	sqlQuery := "SELECT * FROM checkout "
	var args []interface{}
	var queryParams []string
	argIdx := 1
	if customerId != ""{
		queryParams = append(queryParams, " \"customerId\"::text = $"+strconv.Itoa(argIdx) +" ")
		args = append(args, customerId)
	}
	if len(queryParams) > 0 {
		allQuery := strings.Join(queryParams, " AND")
		sqlQuery += " WHERE "+ allQuery
	}
	orderBy := true
	var orderQuery []string
	if createdAt == "" {
		orderQuery = append(orderQuery, " \"createdAt\" DESC")
	}else {		
		if createdAt == "asc"{
			orderQuery = append(orderQuery, " \"createdAt\" ASC")
		} 
	}
	if orderBy {
		sqlQuery += " ORDER BY " + strings.Join(orderQuery, ",")
	}
	sqlQuery += " LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)
	
	conn := db.CreateConn()
	getCheckouts := make([]models.GetCheckout,0)
	checkouts := make([]models.Checkout,0)
	err := conn.Select(&checkouts, sqlQuery, args...)
	if err != nil {
		
		if err ==sql.ErrNoRows{
			c.JSON(200, gin.H{"message":"success", "data":getCheckouts})
		} else {
			c.JSON(500, gin.H{"message":"server error"})
			
		}
		return
	}
	for _, checkout := range checkouts {
		// find checkoutItem
		var getCheckout models.GetCheckout
		var checkoutItem []models.CheckoutItem
		err = conn.Select(&checkoutItem, "SELECT * FROM \"checkoutItem\" WHERE \"checkoutId\" = $1", checkout.Id)
		if err != nil {
			c.JSON(500, gin.H{"message":"server error"})
			return
		}
		var productDetails []forms.ProductDetail
		for _, item := range checkoutItem {
			var productDetail forms.ProductDetail
			productDetail.ProductId = item.ProductId
			productDetail.Quantity = item.Quantity
			productDetails = append(productDetails, productDetail)
		}
		getCheckout.TransactionId = checkout.Id
		getCheckout.CustomerId = checkout.CustomerId
		getCheckout.ProductDetails = productDetails
		getCheckout.Paid = checkout.Paid
		getCheckout.Change = checkout.Change
		getCheckout.CreatedAt = checkout.CreatedAt
		getCheckouts = append(getCheckouts, getCheckout)
	}
	c.JSON(200, gin.H{"message":"success","data":getCheckouts})
}