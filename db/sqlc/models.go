// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Account struct {
	ID        int64       `db:"id"`
	Owner     string      `db:"owner"`
	Balance   int64       `db:"balance"`
	Currency  string      `db:"currency"`
	CreatedAt time.Time   `db:"created_at"`
	Version   pgtype.Int4 `db:"version"`
}

type Entry struct {
	ID        int64     `db:"id"`
	AccountID int64     `db:"account_id"`
	Amount    int64     `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
}

type Session struct {
	ID           uuid.UUID   `db:"id"`
	Username     string      `db:"username"`
	RefreshToken string      `db:"refresh_token"`
	UserAgent    string      `db:"user_agent"`
	ClientIp     string      `db:"client_ip"`
	IsBlocked    pgtype.Bool `db:"is_blocked"`
	ExpiresAt    time.Time   `db:"expires_at"`
	CreatedAt    time.Time   `db:"created_at"`
}

type Transfer struct {
	ID            int64     `db:"id"`
	FromAccountID int64     `db:"from_account_id"`
	ToAccountID   int64     `db:"to_account_id"`
	Amount        int64     `db:"amount"`
	CreatedAt     time.Time `db:"created_at"`
}

type User struct {
	Username          string    `db:"username"`
	FullName          string    `db:"full_name"`
	Email             string    `db:"email"`
	HashedPassword    string    `db:"hashed_password"`
	PasswordChangedAt time.Time `db:"password_changed_at"`
	CreatedAt         time.Time `db:"created_at"`
	IsEmailActivated  bool      `db:"is_email_activated"`
}

type VerifyEmail struct {
	ID         int64     `db:"id"`
	Username   string    `db:"username"`
	Email      string    `db:"email"`
	SecretCode string    `db:"secret_code"`
	IsUsed     bool      `db:"is_used"`
	CreatedAt  time.Time `db:"created_at"`
	ExpiredAt  time.Time `db:"expired_at"`
}
