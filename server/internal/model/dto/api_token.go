package dto

import (
	"errors"
	"strings"
	"time"

	"github.com/fressive/pocman/server/internal/data"
	"github.com/fressive/pocman/server/internal/util"
	"gorm.io/gorm"
)

type APIToken struct {
	gorm.Model
	Hash        string    `json:"token,omitempty" gorm:"size:128;not null;uniqueIndex"`
	Name        string    `json:"name" gorm:"size:128;not null"`
	Description string    `json:"description" gorm:"size:1024"`
	ValidBefore time.Time `json:"valid_before" gorm:"not null;index"`
	LastUsed    time.Time `json:"last_used" gorm:"not null;index"`
}

func (t APIToken) IsExpired() bool {
	if t.ValidBefore.IsZero() {
		return false
	}

	return t.ValidBefore.Compare(time.Now()) == -1
}

var ErrTokenInvalid = errors.New("token invalid")
var ErrTokenExpired = errors.New("token expired")

func VerifyToken(token string) error {
	hash := util.HashToken(token)

	var record APIToken
	if err := data.DB.Where("hash = ?", hash).First(&record).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrTokenInvalid
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if record.Hash != hash {
		return ErrTokenInvalid
	}

	if record.IsExpired() {
		return ErrTokenExpired
	}

	// touch the token
	data.DB.Model(&record).Updates(APIToken{
		LastUsed: time.Now(),
	})

	return nil
}

func NewAPIToken(name, description string, validFor time.Duration) (*APIToken, string, error) {
	token, err := util.GenerateAPIToken()
	if err != nil {
		return nil, "", err
	}

	if validFor <= 0 {
		validFor = 24 * time.Hour
	}

	// do not save the original
	name = strings.TrimSpace(name)

	hashed := util.HashToken(token)
	if name == "" {
		name = "temp-token-" + hashed[:6]
	}

	tokenModel := &APIToken{
		Hash:        hashed,
		Name:        name,
		Description: strings.TrimSpace(description),
		ValidBefore: time.Now().Add(validFor),
	}

	err = data.DB.Create(tokenModel).Error
	if err != nil {
		return nil, "", err
	}

	return tokenModel, token, nil
}
