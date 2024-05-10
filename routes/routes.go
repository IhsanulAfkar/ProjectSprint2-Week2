package routes

import (
	"Week2/controllers"
	"Week2/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	staffController := new(controllers.StaffController)
	productController := new(controllers.ProductController)
	customerController := new(controllers.CustomerController)
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})
	v1 := router.Group("/v1")
	{
		staff := v1.Group("/staff")
		{
			staff.POST("/register", staffController.SignUp)
			staff.POST("/login", staffController.SignIn)
		}
		v1.GET("/product/customer", productController.SearchBySKU)
		v1.Use(middleware.AuthMiddleware)
		product := v1.Group("/product")
		{
			product.POST("/", productController.CreateProduct)
			product.GET("/", productController.GetAllProduct)
			product.PUT("/:productId", productController.UpdateProduct)
			product.DELETE("/:productId", productController.DeleteProduct)
			product.POST("/checkout", productController.Checkout)
			product.GET("/checkout/history", productController.GetAllCheckout)
		}
		customer := v1.Group("/customer")
		{
			customer.POST("/register", customerController.CreateCustomer)
			customer.GET("/", customerController.GetAllCustomer)
		} 
	}
	return router
}