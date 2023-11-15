package handlers

import (
	"errors"
	"github.com/elina-chertova/loyalty-system/internal/balance/service"
	"github.com/elina-chertova/loyalty-system/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type AuthService interface {
	Register(login, password string, isAdmin bool) error
	Login(login, password string) (bool, error)
	SetToken(login string) (uuid.UUID, string, error)
}

type AuthHandler struct {
	balanceService *service.UserBalance
	Auth           AuthService
}

func NewAuthHandler(userBal *service.UserBalance, userAuth AuthService) *AuthHandler {
	return &AuthHandler{balanceService: userBal, Auth: userAuth}
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

// RegisterHandler @User Registration
// @Description User Registration and creating an empty user balance
// @ID register-user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body LoginForm true "User login and password"
// @Success 200 {object} ResponseWithToken
// @Failure 400 {object} Response
// @Failure 409 {object} Response
// @Failure 500 {object} Response
// @Router /register [post]
func (auth *AuthHandler) RegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var login LoginForm
		if err := c.BindJSON(&login); err != nil {
			logger.Logger.Error(
				"Wrong entered data",
				zap.String("endpoint", c.Request.URL.Path),
				zap.Error(err),
			)
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
			logger.Logger.Error(
				"Wrong entered data",
				zap.String("endpoint", c.Request.URL.Path),
				zap.Error(errors.New("login or password is empty")),
			)
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
			logger.Logger.Error(
				"User is already registered",
				zap.String("endpoint", c.Request.URL.Path),
				zap.Error(errors.New(err.Error())),
			)
			c.JSON(
				http.StatusConflict, Response{
					Message: err.Error(),
					Status:  "User is already registered",
				},
			)
			c.Abort()
			return
		}

		userID, token, err := auth.Auth.SetToken(login.Name)
		err = auth.balanceService.AddInitialBalance(userID)
		if err != nil {
			logger.Logger.Error(
				"Error initialize balance",
				zap.String("endpoint", c.Request.URL.Path),
				zap.Error(err),
			)
			c.Abort()
			return

		}
		if err != nil {
			logger.Logger.Error(
				"Error with getting token",
				zap.String("endpoint", c.Request.URL.Path),
				zap.Error(errors.New(err.Error())),
			)
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

// LoginHandler @User Login
// @Description User Login and Set token
// @ID login-user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body LoginForm true "User login and password"
// @Success 200 {object} ResponseWithToken
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /login [post]
func (auth *AuthHandler) LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var login LoginForm
		if err := c.BindJSON(&login); err != nil {
			logger.Logger.Error(
				"Wrong entered data",
				zap.String("endpoint", c.Request.URL.Path),
				zap.Error(err),
			)
			c.JSON(
				http.StatusBadRequest, Response{
					Message: "Check json input",
					Status:  "Wrong entered data",
				},
			)
			c.Abort()
			return
		}

		_, err := auth.Auth.Login(login.Name, login.Password)
		if err != nil {
			logger.Logger.Error(
				"Login failed",
				zap.String("endpoint", c.Request.URL.Path),
				zap.Error(err),
			)
			c.JSON(
				http.StatusUnauthorized, Response{
					Message: err.Error(),
					Status:  "Login failed",
				},
			)
			c.Abort()
			return
		}

		_, token, err := auth.Auth.SetToken(login.Name)
		if err != nil {
			logger.Logger.Error(
				"Error with getting token",
				zap.String("endpoint", c.Request.URL.Path),
				zap.Error(err),
			)
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
