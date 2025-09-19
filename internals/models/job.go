package models

import "time"



type Job struct{
	Id string `json:"id" db:"id"`
	Status string `json:"status" db:"status"`
	Type string `json:"type" db:"type"`
	Payload string	`json:"payload" db:"payload"`
	CreatedAt time.Time`json:"created_at" db:"created_at"`
	UpdatedAt time.Time`json:"updated_at" db:"updated_at"`
}


const (
	StatusQueued = "queued"
	StatusProcessing = "processing"
	StatusCompleted = "completed"
	StatusFailed = "failed"
)