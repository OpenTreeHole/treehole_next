package utils

type Model interface {
	GetID() int
}

func binarySearch[T Model](models []T, targetID int) int {
	left := 0
	right := len(models)
	for left < right {
		mid := left + (right-left)>>1
		if models[mid].GetID() < targetID {
			left = mid + 1
		} else if models[mid].GetID() > targetID {
			right = mid
		} else {
			return mid
		}
	}
	return -1
}

func OrderInGivenOrder[T Model](models []T, order []int) []T {
	var result []T
	for _, i := range order {
		index := binarySearch(models, i)
		if index >= 0 {
			result = append(result, models[index])
		}
	}
	return result
}
