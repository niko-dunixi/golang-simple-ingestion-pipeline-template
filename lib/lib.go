package lib

import (
	"time"

	"github.com/google/uuid"
)

type PayloadItem struct {
	ID       uuid.UUID `json:"id"`
	Time     time.Time `json:"time"`
	TaskName string    `json:"task_name"`
}
