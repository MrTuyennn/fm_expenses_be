package response

import (
	"time"

	"fm_expenses_be/internal/entity"
)

// Category is the wire format for a category.
type Category struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FromCategory maps a domain entity to its wire format.
func FromCategory(c entity.Category) Category {
	return Category{
		ID:        c.ID,
		UserID:    c.UserID,
		Name:      c.Name,
		Icon:      c.Icon,
		Type:      c.Type,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// FromCategoryList maps a slice of domain entities to their wire format.
func FromCategoryList(items []entity.Category) []Category {
	out := make([]Category, 0, len(items))
	for _, c := range items {
		out = append(out, FromCategory(c))
	}
	return out
}
