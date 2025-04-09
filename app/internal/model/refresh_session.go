package model

import "time"

type RefreshSession struct {
	ID           uint      `db:"id"`
	UserID       uint      `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	ExpireAt     time.Time `db:"expire_at"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
