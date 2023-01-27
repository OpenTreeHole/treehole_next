package models

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strings"
)

type UserFavorites struct {
	UserID int `json:"user_id" gorm:"primaryKey"`
	HoleID int `json:"hole_id" gorm:"primaryKey"`
}

func UserCreateFavourite(tx *gorm.DB, c *fiber.Ctx, clear bool, userID int, holeIDs []int) error {
	if clear {
		DB.Exec("DELETE FROM user_favorites WHERE user_id = ?", userID)
	}

	if len(holeIDs) == 0 {
		return nil
	}

	var builder strings.Builder

	if DBType == DBTypeSqlite {
		builder.WriteString("INSERT INTO")
	} else {
		builder.WriteString("INSERT IGNORE INTO")
	}
	builder.WriteString(" user_favorites (user_id, hole_id) VALUES ")
	for i, holeID := range holeIDs {
		builder.WriteString(fmt.Sprintf("(%d, %d)", userID, holeID))
		if i != len(holeIDs)-1 {
			builder.WriteString(", ")
		}
	}

	if DBType == DBTypeSqlite {
		builder.WriteString(" ON CONFLICT DO NOTHING")
	}
	result := tx.Exec(builder.String())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 0 {
		c.Status(201)
	}
	return nil
}

func UserDeleteFavorite(userID int, holeIDs []int) error {
	sql := "DELETE FROM user_favorites WHERE user_id = ? AND hole_id IN ?"
	result := DB.Exec(sql, userID, holeIDs)
	return result.Error
}

func UserGetFavoriteData(userID int) ([]int, error) {
	data := make([]int, 0, 10)
	err := DB.Raw("SELECT hole_id FROM user_favorites WHERE user_id = ?", userID).Scan(&data).Error
	return data, err
}
