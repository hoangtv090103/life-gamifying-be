package handlers

import (
	"life-gamifying/internal/database"
	"life-gamifying/internal/models"
	"life-gamifying/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterHandler is a handler for registering a user
func Register(ctx *gin.Context, s database.Service) error {
	var user models.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	user.Password, err = utils.HashPassword(user.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	err = s.DB().Create(&user).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	ctx.JSON(http.StatusCreated, user)

	return nil
}
