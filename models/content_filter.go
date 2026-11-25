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

// UserHiddenTag stores one full tag name for an individual user.
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

// HiddenTagIDs resolves the account's explicit list of tag names.
func HiddenTagIDs(db *gorm.DB, userID int) *gorm.DB {
	return db.Session(&gorm.Session{NewDB: true}).Model(&Tag{}).Select("tag.id").
		Where(`EXISTS (SELECT 1 FROM user_hidden_tag rule WHERE rule.user_id = ?
			AND LOWER(rule.name) = LOWER(tag.name))`, userID)
}

// FilterHoles applies account-specific division and tag restrictions to a
// query whose hole table is named "hole".
func FilterHoles(db *gorm.DB, userID int) *gorm.DB {
	return db.Where(`NOT EXISTS (SELECT 1 FROM user_hidden_division uhd WHERE uhd.user_id = ? AND uhd.division_id = hole.division_id)`, userID).
		Where(`NOT EXISTS (SELECT 1 FROM hole_tags ht WHERE ht.hole_id = hole.id AND ht.tag_id IN (?))`, HiddenTagIDs(db, userID))
}

func VisibleHoleIDs(db *gorm.DB, userID int) *gorm.DB {
	return FilterHoles(db.Model(&Hole{}), userID).Select("hole.id")
}

// CollectionIDs returns a user's saved hole IDs using the same visibility
// policy for list requests and for responses after a collection is edited.
func CollectionIDs(tx *gorm.DB, userID int) *gorm.DB {
	enabled := tx.Model(&UserContentFilter{}).Select("1").Where("user_id = ?", userID)
	return tx.Where("(NOT EXISTS (?) OR hole_id IN (?))", enabled, VisibleHoleIDs(tx, userID))
}
