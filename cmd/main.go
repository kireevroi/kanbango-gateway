package main

import (
	"context"
	// "fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/kireevroi/kanbango/gateway/internal/userproto"
)



type user struct {
	user string `json:"user"`
	password string `json:"password"`
}


func main() {
	conn, err := grpc.Dial("kanbango.ru:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Print("WTF happened")
	}
	defer conn.Close()
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.POST("/login", func(c *gin.Context) {
		cookie, err := c.Cookie("auth_cookie")
		if err != nil {
			log.Printf("Cookie problems: %v", err)
			// Send cookie data to login or whatever
		}
		log.Printf("Cookie : %v", cookie)
		var newuser user
		if err := c.BindJSON(&newuser); err != nil {
			log.Printf("Something strange happened during login: %v", err)
		}
		client := pb.NewUserServiceClient(conn)
		login, err := client.LoginUser(context.Background(), &pb.LoginUserRequest{Username: newuser.user, Password: newuser.password})
		if err != nil {
			log.Printf("Some kind of error getting GRPC: %v", err)
		}
		if login.Status == pb.Status_STATUS_OK {
			c.SetCookie("auth_cookie", login.Uuid, 3600, "/", "kanbango.ru", false, false)
			c.IndentedJSON(http.StatusOK, gin.H{"Status": "Logged in"})
		} else {
			log.Print(login.Status)
			c.IndentedJSON(http.StatusAccepted, gin.H{"status": login.Status})
		}
	})
	r.Run(":8080")
}