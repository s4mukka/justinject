package domain

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	Id        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Query         string
	UpperBound    int64
	LowerBound    int64
	NumPartitions int64
	ExtractorId   string
	Extractor     Extractor
}
