package favourite

type Response struct {
	Message string `json:"message"`
	Data    []int  `json:"data"`
}

type ListFavoriteModel struct {
	Order           string `json:"order" query:"order" validate:"omitempty,oneof=id time_created hole_time_updated" default:"time_created"`
	Plain           bool   `json:"plain" default:"false" query:"plain"`
	FavoriteGroupID int    `json:"favorite_group_id" default:"0" query:"favorite_group_id"`
}

type AddModel struct {
	HoleID          int `json:"hole_id"`
	FavoriteGroupID int `json:"favorite_group_id" default:"0"`
}

type ModifyModel struct {
	HoleIDs         []int `json:"hole_ids"`
	FavoriteGroupID int   `json:"favorite_group_id" default:"0"`
}

type DeleteModel struct {
	HoleID          int `json:"hole_id"`
	FavoriteGroupID int `json:"favorite_group_id" default:"0"` //ambiguous
}

type AddFavoriteGroupModel struct {
	Name string `json:"name" validate:"required,max=64"`
}

type ModifyFavoriteGroupModel struct {
	Name            string `json:"name" validate:"required,max=64"`
	FavoriteGroupID int    `json:"favorite_group_id" validate:"required"`
}

type MoveModel struct {
	HoleIDs             []int `json:"hole_ids"`
	FromFavoriteGroupID int   `json:"from_favorite_group_id" default:"0"`
	ToFavoriteGroupID   int   `json:"to_favorite_group_id" validate:"required"`
}

type ListFavoriteGroupModel struct {
	Order string `json:"order" query:"order" validate:"omitempty,oneof=id time_created time_updated hole_time_updated" default:"time_created"`
	Plain bool   `json:"plain" default:"false" query:"plain"`
}
