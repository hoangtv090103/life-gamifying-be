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

func GetAllHabits(ctx *gin.Context, s database.Service) {
	client := s.RDB()

	// Get habit list in cache
	cachedHabits, err := client.Get(ctx, "habits:all").Result()

	if err != nil {
		log.Println(err)
	}

	if cachedHabits != "" {
		// If cache exists, return cache
		cachedHabitsJSON := []byte(cachedHabits)
		var habits []models.Habit
		json.Unmarshal(cachedHabitsJSON, &habits)

		ctx.JSON(http.StatusOK, habits)
		return
	}

	// Set habit list in cache
	LoadHabitAll(ctx, s)
}

func GetHabitByID(ctx *gin.Context, s database.Service) {
	var habit models.Habit
	client := s.RDB()
	db := s.DB()

	// Get habit in cache
	cachedHabit, err := client.Get(ctx, "habit:"+ctx.Param("id")).Result()

	if err != nil {
		log.Println(err)
	}

	if cachedHabit != "" {
		// If cache exists, return cache
		cachedHabitJSON := []byte(cachedHabit)
		json.Unmarshal(cachedHabitJSON, &habit)

		ctx.JSON(http.StatusOK, habit)
		return
	}

	// If cache does not exist, get from database
	db.First(&habit, ctx.Param("id"))

	// Set habit in cache
	habitJSON, _ := json.Marshal(habit)
	client.Set(ctx, "habit:"+ctx.Param("id"), habitJSON, 0)

	ctx.JSON(http.StatusOK, habit)
}

func CreateHabit(ctx *gin.Context, s database.Service) {
	var habit models.Habit
	client := s.RDB()
	db := s.DB()

	if err := ctx.ShouldBindJSON(&habit); err != nil {
		log.Println(habit)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Create(&habit)
	// Update habits:all
	LoadHabitAll(ctx, s)

	habitsJSON, _ := json.Marshal(habit)
	client.Set(ctx, "habit:"+strconv.Itoa(int(habit.ID)), habitsJSON, 0)
	


	ctx.JSON(http.StatusCreated, habit)
}

func UpdateHabit(ctx *gin.Context, s database.Service) {
	// TODO: Update habit in cache
	var habit models.Habit
	client := s.RDB()
	db := s.DB()

	if err := ctx.ShouldBindJSON(&habit); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&models.Habit{}).Where("id = ?", ctx.Param("id")).Updates(&habit)

	habitsJSON, _ := json.Marshal(habit)
	client.Set(ctx, "habit:"+ctx.Param("id"), habitsJSON, 0)

	// Update habits:all
	log.Println("Update player:" + ctx.Param("id") + ":habits")
	LoadHabitAll(ctx, s)

	ctx.JSON(http.StatusOK, habit)
}

func DeleteHabit(ctx *gin.Context, s database.Service) {
	var habit models.Habit
	client := s.RDB()
	db := s.DB()

	db.First(&habit, ctx.Param("id"))
	db.Delete(&habit, ctx.Param("id"))

	client.Del(ctx, "habit:"+ctx.Param("id"))

	ctx.JSON(http.StatusOK, gin.H{"id" + ctx.Param("id"): "deleted"})

	// Update habits:all
	LoadHabitAll(ctx, s)
}

func LoadHabitAll(ctx *gin.Context, s database.Service) {
	var habits []models.Habit
	client := s.RDB()
	db := s.DB()

	player_id := ctx.GetHeader("player_id")
	log.Println(player_id)

	db.Model(&models.Habit{}).Find(&habits)

	if len(habits) == 0 {
		// No content
		ctx.JSON(http.StatusNoContent, habits)
		return
	}

	habitsJSON, err := json.Marshal(habits)

	if err != nil {
		log.Println(err)
	}

	log.Println(ctx.Param("id"))	
	client.Set(ctx, "player:" + player_id + ":habits", habitsJSON, 0)
}
