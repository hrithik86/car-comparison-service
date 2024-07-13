package model

import (
	"github.com/google/uuid"
	"time"
)

type DbTimeAudit struct {
	CreatedOn *time.Time `gorm:"column:created_on;default:now()" json:"created_on"`
	UpdatedOn *time.Time `gorm:"column:updated_on;default:NULL;autoUpdateTime:now()" json:"updated_on"`
}

type DbId struct {
	Id *uuid.UUID `gorm:"column:id;primaryKey;type:uuid" sql:"default:uuid_generate_v4()" json:"id"`
}
