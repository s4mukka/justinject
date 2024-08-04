package domain

import (
	"time"

	"gorm.io/gorm"
)

type CreateJobRequest struct {
	ExtractorId   string `json:"extractorId"   binding:"required"`
	Query         string `json:"query"         binding:"required"`
	LowerBound    int64  `json:"lowerBound"`
	UpperBound    int64  `json:"upperBound"`
	NumPartitions int64  `json:"numPartitions"`
}

type Job struct {
	Id        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CreateJobRequest
	Extractor Extractor
}

type IJobUseCase interface {
	CreateJob(request CreateJobRequest) (*Job, error)
}

type IJobRepository interface {
	GetExtractorById(extractorId string) (*Extractor, error)
	CreateJob(job *Job) error
}
