package utils

import (
	"crypto/rand"
	"encoding/hex"
	"spinLuck/internal/shared/models"
)

func NewSession(user *models.User) map[string]any {
	authStage := "full"
	mfaVerified := true

	if user.TwoFactorEnabled {
		authStage = "pending_2fa"
		mfaVerified = false
	}

	return map[string]any{
		"user_id":            user.ID,
		"two_factor_enabled": user.TwoFactorEnabled,
		"auth_stage":         authStage,
		"mfa_verified":       mfaVerified,
	}
}

func GenerateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
