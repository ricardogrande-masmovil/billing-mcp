package persistence

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (bm *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if bm.ID == uuid.Nil {
		bm.ID = uuid.New()
	}
	return
}