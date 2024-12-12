package commands

import (
	"golangsevillabar/domain"
	"slices"

	funk "github.com/thoas/go-funk"
)

func FindMenuItemsThatAreNotInOrderedItems(orderedItems []domain.OrderedItem, markingServedItems []int) []int {
	orderedMenuItems := funk.Map(orderedItems, func(item domain.OrderedItem) int { return item.MenuItem }).([]int)
	result := make([]int, 0)
	for _, v := range markingServedItems {
		if !slices.Contains(orderedMenuItems, v) {
			result = append(result, v)
		}
	}
	return result
}
