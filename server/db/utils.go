package db

import "github.com/AppleGamer22/raker/shared/types"

func (user *User) SelectedCategories(categories []string) map[string]bool {
	result := make(map[string]bool)
	for _, category := range user.Categories {
		result[category] = true
	}
	for _, category := range categories {
		if _, ok := result[category]; ok {
			result[category] = false
		}
	}
	for category, checked := range result {
		result[category] = !checked
	}
	return result
}

func SelectedMediaTypes(mediaTypes []string) map[string]bool {
	result := make(map[string]bool)
	result[types.Instagram] = true
	result[types.Highlight] = true
	result[types.Story] = true
	result[types.VSCO] = true
	result[types.TikTok] = true
	if len(mediaTypes) > 0 {
		for _, mediaType := range mediaTypes {
			if _, ok := result[mediaType]; ok && types.ValidMediaType(mediaType) {
				result[mediaType] = false
			}
		}
		for mediaType, checked := range result {
			result[mediaType] = !checked
		}
	}
	return result
}
