package commands

import (
	"cqrseventsourcingbar/shared"
	"slices"

	funk "github.com/thoas/go-funk"
)

func FindMenuItemsThatAreNotInOrderedItems(orderedItems []shared.MenuItem, markingServedItems []int) []int {
	orderedMenuItems := funk.Map(orderedItems, func(item shared.MenuItem) int { return item.ID }).([]int)
	result := make([]int, 0)
	for _, v := range markingServedItems {
		if !slices.Contains(orderedMenuItems, v) {
			result = append(result, v)
		}
	}
	return result
}
