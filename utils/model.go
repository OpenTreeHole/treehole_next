package utils

type IDModel[T any] interface {
	*T
	GetID() int
}

func binarySearch[T any, PT IDModel[T]](models []PT, targetID int) int {
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

func OrderInGivenOrder[T any, PT IDModel[T]](models []PT, order []int) (result []PT) {
	for _, i := range order {
		index := binarySearch(models, i)
		if index >= 0 {
			result = append(result, models[index])
		}
	}
	return result
}

func Models2IDSlice[T any, PT IDModel[T]](models []PT) (result []int) {
	result = make([]int, len(models))
	for i := range models {
		result[i] = models[i].GetID()
	}
	return result
}
