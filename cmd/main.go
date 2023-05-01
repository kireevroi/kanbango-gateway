package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kireevroi/kanbango/gateway/internal/routes"
)


func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	a := routes.NewAuth("kanbango.ru", "50051", false)
	err := a.Route()
	if err != nil {
		log.Fatalf("unable to establish route to auth server: %v", err)
	}
	defer a.Close()
	r.POST("/login", a.LoginHandler())


	r.Run(":8080")
}