package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
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

// Both production name columns use the same case-insensitive MySQL collation.
// Keep the rule column indexable there. SQLite needs explicit case folding.
func tagNameMatch(db *gorm.DB, column string) string {
	if db.Dialector.Name() == "mysql" {
		return "r.name = " + column
	}
	return "LOWER(r.name) = LOWER(" + column + ")"
}

// FilterHoles applies account-specific division and tag restrictions to a
// query whose hole table is named "hole".
func FilterHoles(db *gorm.DB, userID int) *gorm.DB {
	return db.Where(`NOT EXISTS (SELECT 1 FROM user_hidden_division uhd WHERE uhd.user_id = ? AND uhd.division_id = hole.division_id)`, userID).
		Where(`NOT EXISTS (SELECT 1 FROM hole_tags ht JOIN tag t ON t.id=ht.tag_id
 JOIN user_hidden_tag r ON r.user_id=? AND `+tagNameMatch(db, "t.name")+` WHERE ht.hole_id=hole.id)`, userID)
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

// FilterUserID returns zero for accounts using the ordinary visibility policy.
func FilterUserID(c *fiber.Ctx) (int, error) {
	enabled, err := ContentFilterEnabled(c)
	if err != nil || !enabled {
		return 0, err
	}
	user, err := GetCurrLoginUser(c)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func VisiblePublicHoleIDs(tx *gorm.DB, uid int) *gorm.DB {
	return VisibleHoleIDs(tx.Session(&gorm.Session{NewDB: true}), uid).Where("hole.hidden = ?", false)
}

func VisibleFloorIDs(tx *gorm.DB, uid int) *gorm.DB {
	return tx.Session(&gorm.Session{NewDB: true}).Model(&Floor{}).Select("floor.id").Where("hole_id IN (?)", VisiblePublicHoleIDs(tx, uid))
}

// SelectableTags includes only names used on at least one visible post.
func SelectableTags(tx *gorm.DB, uid int) *gorm.DB {
	q := tx.Model(&Tag{})
	if uid == 0 {
		return q
	}
	linked := VisiblePublicHoleIDs(tx, uid).Select("1").
		Joins("JOIN hole_tags link ON link.hole_id=hole.id").Where("link.tag_id=tag.id").Limit(1)
	// A scalar LIMIT stops at the first visible post even for popular tags.
	return q.Where("NOT EXISTS (SELECT 1 FROM user_hidden_tag r WHERE r.user_id = ? AND "+tagNameMatch(tx, "tag.name")+")", uid).Where("1 = (?)", linked)
}

func RequireVisibleDivision(tx *gorm.DB, uid, id int) error {
	if uid == 0 {
		return nil
	}
	var count int64
	if err := VisibleDivisions(tx, uid).Model(&Division{}).Where("division.id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return common.NotFound()
	}
	return nil
}

func RequireVisibleHoles(tx *gorm.DB, uid int, ids []int) error {
	if uid == 0 || len(ids) == 0 {
		return nil
	}
	var count int64
	if err := VisiblePublicHoleIDs(tx, uid).Where("hole.id IN ?", ids).Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return common.NotFound()
	}
	return nil
}

func RequireCollectionTargets(c *fiber.Ctx, tx *gorm.DB, ids []int) error {
	uid, err := FilterUserID(c)
	if err != nil {
		return err
	}
	return RequireVisibleHoles(tx, uid, ids)
}
