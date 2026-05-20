package model

import "github.com/google/uuid"

// User maps to the public.users table — an application profile linked to a
// Supabase Auth user. It embeds Base (id, uid, timestamps) and doubles as the
// GORM entity. The JSON output stays flat: Base's fields are promoted.
type User struct {
	Base
	AuthUserID uuid.UUID `json:"authUserId"`
	Email      string    `json:"email"`
	FullName   string    `json:"fullName"`
}

// TableName tells GORM which table this model maps to.
func (User) TableName() string {
	return "users"
}
