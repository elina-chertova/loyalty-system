package handlers

import (
	"net/http"

	"github.com/elina-chertova/loyalty-system/internal/auth/handlers"
	"github.com/elina-chertova/loyalty-system/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Ping @Ping database
// @Description Ping database
// @ID ping-db
// @Tags Database
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /ping [get]
func Ping(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			logger.Logger.Error("failed to get database connection", zap.Error(err))
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				handlers.Response{
					Message: "Failed to get database connection",
					Status:  "Error with database connection",
				},
			)
			return
		}

		if err := sqlDB.Ping(); err != nil {
			logger.Logger.Error("failed to ping the database", zap.Error(err))
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				handlers.Response{
					Message: "Database has not been connected",
					Status:  "Error with database connection",
				},
			)
			return
		}
		c.JSON(
			http.StatusOK,
			handlers.Response{
				Message: "Successfully connected to the database and pinged it",
				Status:  "Pinged",
			},
		)
	}
}
