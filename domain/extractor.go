package domain

import (
	"time"

	"gorm.io/gorm"
)

type Extractor struct {
	Id        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Driver
}

type IExtractorRepository interface {
	GetById(extractorId string) (*Extractor, error)
}
