package handlers

import (
	"encoding/json"
	"life-gamifying/internal/database"
	"life-gamifying/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
