package handlers

import (
	"encoding/json"
	"life-gamifying/internal/database"
	"life-gamifying/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPlayers(ctx *gin.Context, s database.Service) {
	var players []models.Player
	client := s.RDB()
	db := s.DB()

	// Get player list in cache
	cachedPlayers, err := client.Get(ctx, "players:all").Result()

	if err != nil {
		log.Println(err)
	}

	if cachedPlayers != "" {
		// If cache exists, return cache
		cachedPlayersJSON := []byte(cachedPlayers)
		json.Unmarshal(cachedPlayersJSON, &players)

		ctx.JSON(http.StatusOK, players)
		return
	}

	// If cache does not exist, get from database
	db.Model(&models.Player{}).Preload("Habits").Find(&players)

	// Set player list in cache
	playersJSON, _ := json.Marshal(players)
	client.Set(ctx, "players:all", playersJSON, 0)

	ctx.JSON(http.StatusOK, players)
}

func GetPlayerByID(ctx *gin.Context, s database.Service) {
	var player models.Player
	client := s.RDB()
	db := s.DB()

	// Get player in cache
	cachedPlayer, err := client.Get(ctx, "player:"+ctx.Param("id")).Result()

	if err != nil {
		log.Println(err)
	}

	if cachedPlayer != "" {
		// If cache exists, return cache
		cachedPlayerJSON := []byte(cachedPlayer)
		json.Unmarshal(cachedPlayerJSON, &player)

		ctx.JSON(http.StatusOK, player)
		return
	}

	// If cache does not exist, get from database
	db.Preload("Habits").First(&player, ctx.Param("id"))

	// Set player in cache
	playerJSON, _ := json.Marshal(player)
	client.Set(ctx, "player:"+ctx.Param("id"), playerJSON, 0)

	ctx.JSON(http.StatusOK, player)
}

func CreatePlayer(ctx *gin.Context, s database.Service) {
	var player models.Player
	db := s.DB()

	if err := ctx.ShouldBindJSON(&player); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Create(&player)

	ctx.JSON(http.StatusCreated, player)
}

func UpdatePlayer(ctx *gin.Context, s database.Service) {
	var player models.Player
	db := s.DB()

	if err := ctx.ShouldBindJSON(&player); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&player)

	ctx.JSON(http.StatusOK, player)
}

func DeletePlayer(ctx *gin.Context, s database.Service) {
	var player models.Player
	db := s.DB()

	db.Delete(&player, ctx.Param("id"))

	ctx.JSON(http.StatusOK, gin.H{"message": "Player deleted"})
}

func GetHabitsOfPlayer(ctx *gin.Context, s database.Service) {
	var habits []models.Habit
	client := s.RDB()
	db := s.DB()

	// Get habit in cache
	cachedHabit, err := client.Get(ctx, "player:"+ctx.Param("id")+":habits").Result()

	if err != nil {
		log.Println(err)
	}

	if cachedHabit != "" {
		// If cache exists, return cache
		cachedHabitJSON := []byte(cachedHabit)
		json.Unmarshal(cachedHabitJSON, &habits)

		ctx.JSON(http.StatusOK, habits)
		return
	}

	// If cache does not exist, get from database
	db.Model(&models.Habit{}).Where("player_id = ?", ctx.Param("id")).Find(&habits)

	// Set habit in cache
	habitJSON, _ := json.Marshal(habits)
	client.Set(ctx, "player:"+ctx.Param("id")+":habits", habitJSON, 0)

	ctx.JSON(http.StatusOK, habits)
}
