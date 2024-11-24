package utils

import (
	"golangsevillabar/domain"
	"slices"

	funk "github.com/thoas/go-funk"
)

func RemoveOrderedItemsThatAppearInMarkedServedItems(orderedItems []domain.OrderedItem, markingServedItems []int) []domain.OrderedItem {
	m := make(map[int]bool)
	for _, v := range markingServedItems {
		m[v] = true
	}

	result := make([]domain.OrderedItem, 0)
	for _, v := range orderedItems {
		if !m[v.MenuItem] {
			result = append(result, v)
		}
	}

	return result
}

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
