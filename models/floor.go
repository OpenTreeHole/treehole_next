package models

// Floor has a tree structure, example:
//	id: 1, reply_to: 0, storey: 1
//		id: 2, reply_to: 1, storey: 1
//	id: 3, reply_to: 0, storey: 2
//		id: 4, reply_to: 3, storey: 2
//			id: 6, reply_to: 4, storey: 2
//		id: 5, reply_to: 3, storey: 2
//	id: 7, reply_to: 0, storey: 3
type Floor struct {
	BaseModel
	HoleID      int      `json:"hole_id"`
	UserID      int      `json:"-"`
	Content     string   `json:"content"`
	Anonyname   string   `json:"anonyname" gorm:"size:16"`
	Storey      int      `json:"storey" gorm:"index"`                    // The sequence of the root nodes
	ReplyTo     int      `json:"reply_to"`                               // Floor id that it replies to (must be in the same hole)
	Mention     []Floor  `json:"mention" gorm:"many2many:floor_mention"` // Many to many mentions (in different holes)
	Like        int      `json:"like" gorm:"index"`                      // like - dislike
	LikeData    IntArray `json:"-"`                                      // user ids
	DislikeData IntArray `json:"-"`                                      // user ids
	Deleted     bool     `json:"deleted"`
	Fold        string   `json:"fold"`        // fold reason
	SpecialTag  string   `json:"special_tag"` // Additional info
}

//goland:noinspection GoNameStartsWithPackageName
type FloorHistory struct {
	BaseModel
	Content string `json:"content"`
	Reason  string `json:"reason"`
	FloorID int    `json:"floor_id"`
	UserID  int    `json:"user_id"` // The one who modified the floor
}
