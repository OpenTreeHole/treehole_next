package models

type DivisionResponse struct {
	BaseModel
	Name        string `json:"name"`
	Description string `json:"description"`
	Pinned      []Hole `json:"pinned"`
}

type Division struct {
	BaseModel
	Name        string   `json:"name" gorm:"unique" `
	Description string   `json:"description"`
	Pinned      IntArray `json:"-"     ` // pinned holes in given order
	Holes       []Hole   `json:"pinned"     `
}

func (division *Division) Preprocess() error {
	var pinned = []int(division.Pinned)
	if len(pinned) == 0 {
		division.Holes = []Hole{}
		return nil
	}
	var holes []Hole
	DB.Find(&holes, pinned)
	orderedHoles := make([]Hole, 0, len(holes))
	for _, order := range pinned {
		// binary search the index
		index := func(target int) int {
			left := 0
			right := len(holes)
			for left < right {
				mid := left + (right-left)>>1
				if holes[mid].ID < target {
					left = mid + 1
				} else if holes[mid].ID > target {
					right = mid
				} else {
					return mid
				}
			}
			return -1
		}(order)
		if index >= 0 {
			orderedHoles = append(orderedHoles, holes[index])
		}
	}
	division.Holes = orderedHoles
	return nil
}
