package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/AppleGamer22/raker/shared/types"
)

type Point struct {
	X float64
	Y float64
}

func (p *Point) Scan(src any) error {
	switch v := src.(type) {
	case string:
		_, err := fmt.Sscanf(v, "(%f,%f)", &p.X, &p.Y)
		return err
	case []byte:
		_, err := fmt.Sscanf(string(v), "(%f,%f)", &p.X, &p.Y)
		return err
	default:
		return fmt.Errorf("cannot scan %T into Point", src)
	}
}

func (p Point) Value() (driver.Value, error) {
	return fmt.Sprintf("(%f,%f)", p.X, p.Y), nil
}

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
