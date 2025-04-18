package types

import (
	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Code   string `gorm:"uniqueIndex" json:"code"`
	URL    string `json:"url"`
	Clicks uint   `json:"clicks"`
}
