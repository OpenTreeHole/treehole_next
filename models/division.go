package models

import (
	"time"

	"treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Division struct {
	/// saved fields
	ID        int       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"time_created" gorm:"not null"`
	UpdatedAt time.Time `json:"time_updated" gorm:"not null"`

	/// base info
	Name        string `json:"name" gorm:"unique;size:10"`
	Description string `json:"description" gorm:"size:64"`
	Hidden      bool   `json:"hidden" gorm:"not null;default:false"`

	// pinned holes in given order
	Pinned []int `json:"-" gorm:"serializer:json;size:100;not null;default:\"[]\""`

	/// association fields, should add foreign key

	// return pinned hole to frontend
	Holes Holes `json:"pinned"`

	/// generated field
	DivisionID int `json:"division_id" gorm:"-:all"`

	ShowInHomePage bool `json:"show_in_home_page" gorm:"not null;default:true"`
}

func HomepageDivisionIDs(tx *gorm.DB, excludeDivisionIDs *[]int) (divisionIDs []int, err error) {
	if excludeDivisionIDs != nil {
		err = tx.Model(&Division{}).Select("id").Where("show_in_home_page = ? AND id NOT IN ?", true, excludeDivisionIDs).Find(&divisionIDs).Error
		return
	}

	err = tx.Model(&Division{}).Select("id").Where("show_in_home_page = ?", true).Find(&divisionIDs).Error
	return
}

func (division *Division) GetID() int {
	return division.ID
}

type Divisions []*Division

// UserHiddenDivision hides a division from a user's division list. It does not
// restrict access to the division or its content through other endpoints.
type UserHiddenDivision struct {
	UserID     int `gorm:"primaryKey;autoIncrement:false"`
	DivisionID int `gorm:"primaryKey;autoIncrement:false"`
}

// VisibleDivisions returns the existing public list with user-specific entries
// excluded. An absent configuration leaves the list unchanged.
func VisibleDivisions(db *gorm.DB, userID int) *gorm.DB {
	hiddenIDs := db.Model(&UserHiddenDivision{}).Select("division_id").Where("user_id = ?", userID)
	return db.Where("hidden = ? AND id NOT IN (?)", false, hiddenIDs)
}

func (divisions Divisions) Preprocess(c *fiber.Ctx) error {
	for _, division := range divisions {
		err := division.Preprocess(c)
		if err != nil {
			return err
		}
	}
	return nil //utils.SetCache("divisions", divisions, 0)
}

func (division *Division) Preprocess(c *fiber.Ctx) error {
	var pinned = division.Pinned
	division.Holes = make(Holes, 0, 10)
	if len(pinned) == 0 {
		return nil
	}
	DB.Find(&division.Holes, pinned)
	if len(division.Holes) == 0 {
		return nil
	}
	division.Holes = utils.OrderInGivenOrder(division.Holes, pinned)
	// division.Holes = division.Holes.RemoveIf(func(hole *Hole) bool {
	// 	return hole.Hidden
	// })
	return division.Holes.Preprocess(c)
}

func (division *Division) AfterFind(_ *gorm.DB) (err error) {
	division.DivisionID = division.ID
	return nil
}

func (division *Division) AfterCreate(_ *gorm.DB) (err error) {
	division.DivisionID = division.ID
	return nil
}
