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
	var habits []models.Habit
	client := s.RDB()
	db := s.DB()

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

	// If cache does not exist, get from database
	db.Model(&models.Habit{}).Preload("PLayers").Find(&habits)

	// Set habit list in cache
	habitsJSON, _ := json.Marshal(habits)
	client.Set(ctx, "habits:all", habitsJSON, 0)

	ctx.JSON(http.StatusOK, habits)
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Create(&habit)

	habitsJSON, _ := json.Marshal(habit)
	client.Set(ctx, "habit:"+strconv.Itoa(int(habit.ID)), habitsJSON, 0)

	ctx.JSON(http.StatusCreated, habit)
}

func UpdateHabit(ctx *gin.Context, s database.Service) {
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

	ctx.JSON(http.StatusOK, habit)
}

func DeleteHabit(ctx *gin.Context, s database.Service) {
	var habit models.Habit
	client := s.RDB()
	db := s.DB()

	db.First(&habit, ctx.Param("id"))
	db.Delete(&habit)

	client.Del(ctx, "habit:"+ctx.Param("id"))

	ctx.JSON(http.StatusOK, gin.H{"id" + ctx.Param("id"): "deleted"})
}
