package category

import (
	"time"

	"fm_expenses_be/internal/entity"
)

// Model is the GORM-tagged persistence row for a category.
type Model struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string    `gorm:"column:user_id;type:uuid;index;not null"`
	Name      string    `gorm:"column:name;not null"`
	Icon      string    `gorm:"column:icon"`
	Type      string    `gorm:"column:type;not null;check:type IN ('income','expense')"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Model) TableName() string { return "categories" }

func (m Model) toEntity() entity.Category {
	return entity.Category{
		ID:        m.ID,
		UserID:    m.UserID,
		Name:      m.Name,
		Icon:      m.Icon,
		Type:      m.Type,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func modelFromEntity(c *entity.Category) *Model {
	return &Model{
		ID:        c.ID,
		UserID:    c.UserID,
		Name:      c.Name,
		Icon:      c.Icon,
		Type:      c.Type,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
