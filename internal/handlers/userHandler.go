package handlers

import (
	"encoding/json"
	"life-gamifying/internal/database"
	"life-gamifying/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUsers(ctx *gin.Context, s database.Service) {
	var users []models.User
	client := s.RDB()
	db := s.DB()

	// Get user list in cache
	cachedUsers, err := client.Get(ctx, "user:all").Result()

	if err != nil {
		log.Println(err)
	}

	if cachedUsers != "" {
		// If cache exists, return cache
		cachedPlayersJSON := []byte(cachedUsers)
		json.Unmarshal(cachedPlayersJSON, &users)

		ctx.JSON(http.StatusOK, users)
		return
	}

	// If cache does not exist, get from database
	db.Model(&models.User{}).Preload("Habits").Find(&users)

	// Set user list in cache
	usersJSON, _ := json.Marshal(users)
	client.Set(ctx, "user:all", usersJSON, 0)

	ctx.JSON(http.StatusOK, users)
}

func GetUserByID(ctx *gin.Context, s database.Service) {
	var user models.User
	client := s.RDB()
	db := s.DB()

	// Get user in cache
	cachedUser, err := client.Get(ctx, "user:"+ctx.Param("id")).Result()

	if err != nil {
		log.Println(err)
	}

	if cachedUser != "" {
		// If cache exists, return cache
		cachedUserJSON := []byte(cachedUser)
		json.Unmarshal(cachedUserJSON, &user)

		ctx.JSON(http.StatusOK, user)
		return
	}

	// If cache does not exist, get from database
	db.Preload("Habits").First(&user, ctx.Param("id"))

	// Set user in cache
	userJSON, _ := json.Marshal(user)
	client.Set(ctx, "user:"+ctx.Param("id"), userJSON, 0)

	ctx.JSON(http.StatusOK, user)
}

func CreateUser(ctx *gin.Context, s database.Service) error {
	var user models.User
	client := s.RDB()
	db := s.DB()

	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	db.Create(&user)

	userJSON, _ := json.Marshal(user)
	client.Set(ctx, "user:"+strconv.Itoa(int(user.ID)), userJSON, 0)

	ctx.JSON(http.StatusCreated, user)

	return nil
}

func UpdateUser(ctx *gin.Context, s database.Service) error {
	var user models.User
	client := s.RDB()
	db := s.DB()

	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	db.Save(&user)

	userJSON, _ := json.Marshal(user)
	client.Set(ctx, "user:"+strconv.Itoa(int(user.ID)), userJSON, 0)

	ctx.JSON(http.StatusOK, user)

	return nil
}

func DeleteUser(ctx *gin.Context, s database.Service) error {
	var user models.User
	client := s.RDB()
	db := s.DB()

	db.First(&user, ctx.Param("id"))
	db.Delete(&user)

	client.Del(ctx, "user:"+ctx.Param("id"))

	ctx.JSON(http.StatusOK, user)

	return nil
}
