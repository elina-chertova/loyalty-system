package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type AuthService interface {
	Register(login, password string, isAdmin bool) error
	Login(login, password string) (bool, error)
	SetToken(login string) (string, error)
}

type AuthHandler struct {
	Auth AuthService
}

func NewAuthHandler(userAuth AuthService) *AuthHandler {
	return &AuthHandler{Auth: userAuth}
}

type LoginForm struct {
	Name     string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ResponseWithToken struct {
	Response
	Token string `json:"token"`
}

func (auth *AuthHandler) RegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var login LoginForm
		if err := c.BindJSON(&login); err != nil {
			c.JSON(
				http.StatusBadRequest, Response{
					Message: "Check json input",
					Status:  "Wrong entered data",
				},
			)
			c.Abort()
			return
		}

		if len(login.Name) == 0 || len(login.Password) == 0 {
			c.JSON(
				http.StatusBadRequest, Response{
					Message: "Check json input",
					Status:  "Wrong entered data",
				},
			)
			c.Abort()
			return
		}

		err := auth.Auth.Register(login.Name, login.Password, false)
		if err != nil {
			c.JSON(
				http.StatusConflict, Response{
					Message: err.Error(),
					Status:  "User is already registered",
				},
			)
			c.Abort()
			return
		}

		token, err := auth.Auth.SetToken(login.Name)
		if err != nil {
			c.JSON(
				http.StatusConflict, Response{
					Message: err.Error(),
					Status:  "Error with getting token",
				},
			)
			c.Abort()
			return
		}

		expirationTime := time.Now().Add(72 * time.Hour)
		cookie := http.Cookie{
			Name:     "access_token",
			Value:    token,
			Expires:  expirationTime,
			HttpOnly: true,
			Secure:   true,
		}

		http.SetCookie(c.Writer, &cookie)
		c.IndentedJSON(
			http.StatusOK, ResponseWithToken{
				Response: Response{
					Message: "Registered",
					Status:  "OK",
				},
				Token: token,
			},
		)
	}
}

func (auth *AuthHandler) LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var login LoginForm
		if err := c.BindJSON(&login); err != nil {
			c.JSON(
				http.StatusBadRequest, Response{
					Message: "Check json input",
					Status:  "Wrong entered data",
				},
			)
			c.Abort()
			return
		}

		isEqual, err := auth.Auth.Login(login.Name, login.Password)
		if err != nil {
			c.JSON(
				http.StatusUnauthorized, Response{
					Message: err.Error(),
					Status:  "Login failed",
				},
			)
			c.Abort()
			return
		}
		if isEqual {

		}
		token, err := auth.Auth.SetToken(login.Name)
		if err != nil {
			c.JSON(
				http.StatusConflict, Response{
					Message: err.Error(),
					Status:  "Error with getting token",
				},
			)
			c.Abort()
			return
		}

		expirationTime := time.Now().Add(72 * time.Hour)
		cookie := http.Cookie{
			Name:     "access_token",
			Value:    token,
			Expires:  expirationTime,
			HttpOnly: true,
			Secure:   true,
		}
		http.SetCookie(c.Writer, &cookie)

		c.IndentedJSON(
			http.StatusOK, ResponseWithToken{
				Response: Response{
					Message: "Login success",
					Status:  "OK",
				},
				Token: token,
			},
		)
	}
}
