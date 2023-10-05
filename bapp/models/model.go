package models

import (
	"bapp/types"
	"encoding/json"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	FirstName    string          `json:"first_name" gorm:"not null"`
	LastName     string          `json:"last_name" gorm:"not null"`
	Email        string          `json:"email" gorm:"not null;unique"`
	PasswordHash string          `json:"password" gorm:"not null"`
	CreatedAt    types.Timestamp `json:"created_at" gorm:"type:timestamptz;autoCreateTime"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Bill struct {
	DueDate   types.Timestamp `json:"due_date" gorm:"type:timestamptz"`
	AmountDue float64         `json:"amount_due"`
	Paid      bool            `json:"paid"`
}

type Payday struct {
	Paydate types.Timestamp `json:"paydate" gorm:"type:timestamptz"`
	Amount  float64         `json:"amount"`
}

type Income struct {
	Amount   float64 `json:"amount"`
	Schedule string  `json:"schedule"`
}

type UserConfig struct {
	ID        int             `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	Balance   float64         `json:"balance"`
	Income    types.JSONB     `json:"income" gorm:"type:jsonb"`
	Paydates  json.RawMessage `json:"paydates" gorm:"type:jsonb"`
	Bills     json.RawMessage `json:"bills" gorm:"type:jsonb"`
	CreatedAt types.Timestamp `json:"created_at" gorm:"type:timestamptz;autoCreateTime"`
	UpdateAt  types.Timestamp `json:"updated_at" gorm:"type:timestamptz;autoUpdateTime"`
}
