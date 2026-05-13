package model

import (
	"time"
)

type AuthActivityLog struct {
	ID           string     `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID       *string    `json:"user_id" gorm:"type:uuid;index"`
	ActivityType string     `json:"activity_type" gorm:"type:varchar(50);not null;index"`
	IPAddress    *string    `json:"ip_address" gorm:"type:varchar(45)"`
	UserAgent    *string    `json:"user_agent" gorm:"type:text"`
	Status       string     `json:"status" gorm:"type:varchar(20);not null;default:'success'"`
	ErrorMessage *string    `json:"error_message" gorm:"type:text"`
	Metadata     *string    `json:"metadata" gorm:"type:jsonb"`
	CreatedAt    time.Time  `json:"created_at" gorm:"not null;default:now()"`
	User         *User      `json:"-" gorm:"foreignKey:UserID"`
}

func (AuthActivityLog) TableName() string {
	return "auth_activity_logs"
}

const (
	ActivityLogin              = "login"
	ActivityLoginFailed        = "login_failed"
	ActivityLogout             = "logout"
	ActivityRegister           = "register"
	ActivityPasswordChange     = "password_change"
	ActivityPasswordResetReq   = "password_reset_request"
	ActivityPasswordReset      = "password_reset"
	ActivityTokenRefresh       = "token_refresh"
	ActivityOAuthLogin         = "oauth_login"
	ActivityOAuthLoginFailed   = "oauth_login_failed"

	StatusSuccess = "success"
	StatusFailure = "failure"
	StatusPending = "pending"
)