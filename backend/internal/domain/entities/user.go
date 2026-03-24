package entities

import "time"

// User represents the domain entity for a user
type User struct {
	ID          string    `json:"id"`
	SupabaseUID string    `json:"supabase_uid"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	Role        string    `json:"role"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// IsAdmin checks if the user has admin privileges
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}
