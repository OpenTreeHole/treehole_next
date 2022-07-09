package models

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"treehole_next/utils"
)

// Floor has a tree structure, example:
//	id: 1, reply_to: 0, storey: 1
//		id: 2, reply_to: 1, storey: 2
//	id: 3, reply_to: 0, storey: 3
//		id: 4, reply_to: 3, storey: 4
//			id: 6, reply_to: 4, storey: 5
//		id: 5, reply_to: 3, storey: 6
//	id: 7, reply_to: 0, storey: 7
type Floor struct {
	BaseModel
	HoleID     int     `json:"hole_id"`
	UserID     int     `json:"-"`
	Content    string  `json:"content"`                                // not empty
	Anonyname  string  `json:"anonyname" gorm:"size:32"`               // random username, not empty
	Storey     int     `json:"storey" gorm:"index"`                    // The sequence of floors in a hole
	ReplyTo    int     `json:"reply_to"`                               // Floor id that it replies to (must be in the same hole)
	Mention    []Floor `json:"mention" gorm:"many2many:floor_mention"` // Many to many mentions (in different holes)
	Like       int     `json:"like" gorm:"index"`                      // like - dislike
	Liked      bool    `json:"liked" gorm:"-:all"`                     // whether the user has liked the floor, dynamically generated
	IsMe       bool    `json:"is_me" gorm:"-:all"`                     // whether the user is the author of the floor, dynamically generated
	Deleted    bool    `json:"deleted"`                                // whether the floor is deleted
	Fold       string  `json:"fold"`                                   // fold reason
	SpecialTag string  `json:"special_tag"`                            // Additional info
}

type Floors []Floor

type AnonynameMapping struct {
	HoleID    int    `json:"hole_id" gorm:"primarykey"`
	UserID    int    `json:"user_id" gorm:"primarykey"`
	Anonyname string `json:"anonyname" gorm:"index;size:32"`
}

//goland:noinspection GoNameStartsWithPackageName
type FloorHistory struct {
	BaseModel
	Content string `json:"content"`
	Reason  string `json:"reason"`
	FloorID int    `json:"floor_id"`
	UserID  int    `json:"user_id"` // The one who modified the floor
}

func (floor *Floor) Preprocess() error {
	// Load mentions
	if floor.Mention == nil {
		floor.Mention = []Floor{}
	}
	return nil
}

func (floors Floors) Preprocess() error {
	for i := 0; i < len(floors); i++ {
		err := floors[i].Preprocess()
		if err != nil {
			return err
		}
	}
	return nil
}

func (floor *Floor) LoadMention() error {
	var Mention Floors
	err := DB.Model(floor).Association("Mention").Find(&Mention)
	if err != nil {
		return err
	}
	if Mention != nil {
		floor.Mention = Mention
	}
	return nil
}

func (floor Floor) MakeQuerySet(
	limit int, offset int,
	holeID int, orderBy string,
	ifDesc bool) (tx *gorm.DB) {
	return DB.
		Limit(limit).Offset(offset).
		Where("hole_id = ?", holeID).
		Order(clause.OrderByColumn{Column: clause.Column{Name: orderBy}, Desc: ifDesc})
}

func (floor *Floor) Create(c *fiber.Ctx, db ...*gorm.DB) error {
	var tx *gorm.DB
	if len(db) > 0 {
		tx = db[0]
	} else {
		tx = DB
	}

	userID, err := GetUserID(c)
	if err != nil {
		return err
	}

	// get anonymous name
	var mapping AnonynameMapping

	err = tx.Transaction(func(tx *gorm.DB) error {
		result := tx.
			Where("hole_id = ?", floor.HoleID).
			Where("user_id = ?", userID).
			Take(&mapping)

		if result.Error != nil {
			// no mapping exists, generate anonyname
			var names []string
			result = tx.Clauses(clause.Locking{
				Strength: "UPDATE",
			}).Raw(`
				SELECT anonyname FROM anonyname_mapping 
				WHERE hole_id = ? 
				ORDER BY anonyname`, floor.HoleID,
			).Scan(&names)
			if result.Error != nil {
				return result.Error
			}

			floor.Anonyname = utils.GenerateName(names)
			result = tx.Create(&AnonynameMapping{
				UserID:    userID,
				HoleID:    floor.HoleID,
				Anonyname: floor.Anonyname,
			})
			if result.Error != nil {
				return result.Error
			}
		} else {
			floor.Anonyname = mapping.Anonyname
		}
		return nil
	})
	if err != nil {
		return err
	}

	result := tx.Create(floor)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (floor *Floor) Backup(c *fiber.Ctx, reason string) error {
	userID, err := GetUserID(c)
	if err != nil {
		return err
	}

	history := FloorHistory{
		Content: floor.Content,
		Reason:  reason,
		FloorID: floor.ID,
		UserID:  userID,
	}
	return DB.Create(&history).Error
}
