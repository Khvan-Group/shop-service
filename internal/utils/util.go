package utils

import "shop-service/internal/items/models"

func GetUniqueStrings(items []models.ItemBasketHistoryView) []string {
	payedAtSet := make(map[string]struct{})
	for _, item := range items {
		payedAtSet[item.PayedAt] = struct{}{}
	}

	var uniquePayedAt []string
	for payedAt := range payedAtSet {
		uniquePayedAt = append(uniquePayedAt, payedAt)
	}

	return uniquePayedAt
}
