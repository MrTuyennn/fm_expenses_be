// Package request holds v1 API request body DTOs.
package request

// Category is the request body for creating/updating a category.
type Category struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
	Type string `json:"type"`
}
