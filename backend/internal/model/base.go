package model

import (
	"time"

	"github.com/google/uuid"
)

// Base holds the standard columns shared by every table: an internal
// autoincrement id, an auto-generated public uid, and creation/update
// timestamps. Embed it (anonymously) in every persisted model so the
// columns and their GORM/JSON mapping are defined in one place.
type Base struct {
	ID        int64     `json:"-"          gorm:"primaryKey"`
	UID       uuid.UUID `json:"uid"        gorm:"not null;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
