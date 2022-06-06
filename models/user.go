package models

type User struct {
	BaseModel
	Favorites []Hole    `json:"favorites,omitempty" gorm:"many2many:user_favorites"`
	Nickname  string    `json:"nickname,omitempty" gorm:"-:migration"`
	Config    StringMap `json:"config,omitempty" gorm:"-:migration"`
	IsAdmin   bool      `json:"is_admin,omitempty" gorm:"-:migration"`
}
