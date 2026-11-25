package models

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UserContentFilter enables the extended content restrictions for one user.
// Existing list-only settings do not implicitly enable this policy.
type UserContentFilter struct {
	UserID int `gorm:"primaryKey;autoIncrement:false"`
}

// UserHiddenTag stores an exact tag name for an individual user.
type UserHiddenTag struct {
	UserID int    `gorm:"primaryKey;autoIncrement:false"`
	Name   string `gorm:"primaryKey;size:64"`
}

const contentFilterLocal = "user_content_filter_enabled"

func ContentFilterEnabled(c *fiber.Ctx) (bool, error) {
	if enabled, ok := c.Locals(contentFilterLocal).(bool); ok {
		return enabled, nil
	}
	user, err := GetCurrLoginUser(c)
	if err != nil {
		return false, err
	}
	var count int64
	if err := DB.Model(&UserContentFilter{}).Where("user_id = ?", user.ID).Count(&count).Error; err != nil {
		return false, err
	}
	enabled := count != 0
	c.Locals(contentFilterLocal, enabled)
	return enabled, nil
}

// FilterHoles applies account-specific division and tag restrictions to a
// query whose hole table is named "hole".
func FilterHoles(db *gorm.DB, userID int) *gorm.DB {
	return db.Where(`NOT EXISTS (SELECT 1 FROM user_hidden_division uhd WHERE uhd.user_id = ? AND uhd.division_id = hole.division_id)`, userID).
		Where(`NOT EXISTS (SELECT 1 FROM hole_tags ht JOIN tag t ON t.id = ht.tag_id JOIN user_hidden_tag uht ON LOWER(uht.name) = LOWER(t.name) WHERE ht.hole_id = hole.id AND uht.user_id = ?)`, userID)
}

func VisibleHoleIDs(db *gorm.DB, userID int) *gorm.DB {
	return FilterHoles(db.Model(&Hole{}), userID).Select("hole.id")
}
