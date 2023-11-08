package handlers

import (
	"github.com/elina-chertova/loyalty-system/internal/auth/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	Auth *service.UserAuth
}

func NewAuthHandler(userAuth *service.UserAuth) *AuthHandler {
	return &AuthHandler{Auth: userAuth}
}

type Login struct {
	Name     string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (auth *AuthHandler) RegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var login Login
		if err := c.BindJSON(&login); err != nil {
			return
		}

		if len(login.Name) == 0 || len(login.Password) == 0 {
			c.JSON(
				http.StatusBadRequest, Response{
					Message: "Check json input",
					Status:  "Wrong entered data",
				},
			)
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
			return
		}

		c.IndentedJSON(
			http.StatusOK, Response{
				Message: "Registered",
				Status:  "OK",
			},
		)
	}
}
