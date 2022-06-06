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
	HoleID      int      `json:"hole_id,omitempty"`
	UserID      int      `json:"-,omitempty"`
	Content     string   `json:"content,omitempty"`
	Anonyname   string   `json:"anonyname,omitempty" gorm:"size:16"`
	Storey      int      `json:"storey,omitempty" gorm:"index"`                    // The sequence of the root nodes
	ReplyTo     int      `json:"reply_to,omitempty"`                               // Floor id that it replies to (must be in the same hole)
	Mention     []Floor  `json:"mention,omitempty" gorm:"many2many:floor_mention"` // Many to many mentions (in different holes)
	Like        int      `json:"like,omitempty" gorm:"index"`                      // like - dislike
	LikeData    IntArray `json:"-,omitempty"`                                      // user ids
	DislikeData IntArray `json:"-,omitempty"`                                      // user ids
	Deleted     bool     `json:"deleted,omitempty"`
	Fold        string   `json:"fold,omitempty"`        // fold reason
	SpecialTag  string   `json:"special_tag,omitempty"` // Additional info
}

//goland:noinspection GoNameStartsWithPackageName
type FloorHistory struct {
	BaseModel
	Content string `json:"content,omitempty"`
	FloorID int    `json:"floor_id,omitempty"`
	UserID  int    `json:"user_id,omitempty"`
}
