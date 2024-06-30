package models

type LevelThresholds struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	Level     uint `json:"level" gorm:"unique;not null"`
	ExpToNext uint `json:"exp_to_next" gorm:"not null"`
}
