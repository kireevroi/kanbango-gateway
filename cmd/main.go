package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kireevroi/kanbango/gateway/internal/routes"
)


func main() {
	r := gin.Default()

	r.GET("/api/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	a := routes.NewAuth("kanbango.ru", "50051", false)
	err := a.Route()
	if err != nil {
		log.Fatalf("unable to establish route to auth server: %v", err)
	}
	defer a.Close()
	r.POST("/api/login", a.LoginHandler())
	r.POST("/api/logout", a.LogoutHandler())
	r.POST("/api/signup", a.SignupHandler())
	r.POST("/api/delete", a.DeleteHandler())
	r.RunTLS(":5535", "./kbg.crt", "./kbg.key")
}