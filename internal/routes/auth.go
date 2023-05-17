package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/kireevroi/kanbango/gateway/internal/userproto"
)

type user struct {
	User string `json:"user"`
	Password string `json:"password"`
}

type Auth struct {
	Url string // connection address
	port string // connection port
	secure bool // True - secure, false - insecure
	conn *grpc.ClientConn // connection
}

func NewAuth(url string, port string, secure bool) (*Auth) {
	var a Auth
	a.Url = url
	a.port = port
	a.secure = secure
	a.conn = nil
	return &a
}

func (a *Auth)Route() error {
	var opt grpc.DialOption
	if a.secure {
		opt = grpc.WithTransportCredentials(insecure.NewCredentials()) // temporary, change to SSL/TLS later
	} else {
		opt = grpc.WithTransportCredentials(insecure.NewCredentials())
	}
	conn, err := grpc.Dial(a.Url+":"+a.port, opt)
	if err != nil {
		a.conn = nil
		return err
	}
	a.conn = conn
	return nil
}

func (a *Auth)Close() error {
	return a.conn.Close()
}

func (a *Auth)LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var u user
		if err := c.BindJSON(&u); err != nil {
			log.Printf("ingress login error: %v", err)
		}
		
		cookie, _ := c.Cookie("auth")
		client := pb.NewUserServiceClient(a.conn)
		log.Println(u.User, cookie)
		login, err := client.LoginUser(context.Background(), &pb.LoginUserRequest{Username: u.User, Password: u.Password, Uuid: cookie})
		if err != nil {
			log.Printf("login grpc connection error: %v", err)
			return
		
		}	
		if login.Status == pb.Status_STATUS_OK {
			c.SetCookie("auth", login.Uuid, 3600, "/", "kanbango.ru", true, true)
			c.JSON(http.StatusOK, gin.H{"Status" : "Logged in"})
			return
		} else if login.Status == pb.Status_STATUS_ALRLOGGED {
			c.JSON(http.StatusOK, gin.H{"Status" : "Already logged in"})
		} else if login.Status == pb.Status_STATUS_NOUSER {
			c.JSON(http.StatusForbidden, gin.H{"Status" : "No such user"})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"Status" : "Wrong something"})
			return
		}
	}
}

func (a *Auth)LogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := pb.NewUserServiceClient(a.conn)
		cookie, _ := c.Cookie("auth")
		client.LogoutUser(context.Background(), &pb.LogoutUserRequest{Uuid: cookie})
		c.SetCookie("auth", "", -1, "/", "kanbango.ru", true, true)
		c.JSON(http.StatusOK, gin.H{"Status" : "Logged out"})
	}
}

func (a *Auth)SignupHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var u user
		if err := c.BindJSON(&u); err != nil {
			log.Printf("ingress singup error: %v", err)
		}
		client := pb.NewUserServiceClient(a.conn)
		_, err := c.Cookie("auth")
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"Status" : "Seems like you are already logged in"})
			return
		}
		cuResponse, err := client.CreateUser(context.Background(), &pb.CreateUserRequest{Username: u.User, Password: u.Password})
		if cuResponse.Status != pb.Status_STATUS_OK {
			c.JSON(http.StatusBadRequest, gin.H{"Status" : "Something wrong"})
			return
		}
		c.JSON(http.StatusOK, gin.H {"Status" : "User created"})
	}
}

func (a *Auth)DeleteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := pb.NewUserServiceClient(a.conn)
		cookie, err := c.Cookie("auth")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H {"Status" : "But youre not logged in..."})
			return
		}
		dResponse, err := client.DeleteUser(context.Background(), &pb.LogoutUserRequest{Uuid: cookie})
		if dResponse.Status != pb.Status_STATUS_OK {
			c.JSON(http.StatusBadRequest, gin.H{"Status" : "Something wrong"})
			return
		}
		c.SetCookie("auth", "", -1, "/", "kanbango.ru", true, true)
		c.JSON(http.StatusOK, gin.H {"Status" : "User deleted"})
	}
}
