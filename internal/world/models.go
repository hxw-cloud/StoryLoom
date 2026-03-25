package world

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WorldSetting represents a foundational rule, concept, or element of the story's world.
// In the context of StoryLoom, this is not just descriptive text, but a "Logic Setting"
// that the Digital Editor will use to validate the consistency of the narrative.
// For example, if a setting dictates "Magic requires verbal incantations," the Editor
// can flag scenes where a gagged character casts a spell.
type WorldSetting struct {
	// ID is the unique identifier for the setting, utilizing UUID for global uniqueness
	// across distributed systems or offline sync scenarios.
	ID string `gorm:"primaryKey;type:varchar(36)"`

	// Category classifies the setting into a distinct domain (e.g., Geography, Magic System).
	// This grouping allows writers and the AI Editor to filter and validate rules contextually.
	// For instance, combat scenes might trigger validation checks specifically against the
	// "Magic System" or "Technology" categories.
	Category string `gorm:"type:varchar(100);not null"`

	// Name is the human-readable title of the setting (e.g., "The Law of Equivalent Exchange").
	// It is indexed to allow for fast autocomplete and search operations in the UI.
	Name string `gorm:"type:varchar(200);not null;index"`

	// Description contains the detailed, human-readable explanation of the setting.
	// This provides the narrative flavor and historical context for the writer.
	Description string `gorm:"type:text"`

	// LogicRules contains the machine-parseable or explicit logical constraints that
	// this setting imposes on the world. This field is critical for the AI "Digital Editor"
	// to perform automated consistency checks against the plot and character actions.
	LogicRules string `gorm:"type:text"`

	// CreatedAt timestamps the creation of the setting, useful for auditing and sorting.
	CreatedAt time.Time

	// UpdatedAt timestamps the last modification, useful for cache invalidation and syncing.
	UpdatedAt time.Time
}

// BeforeCreate is a GORM hook that executes prior to inserting a new record into the database.
// We use this hook to automatically generate a UUID for the WorldSetting if one hasn't
// been provided. This ensures that our primary keys are always valid, distributed-friendly
// identifiers without requiring the handler or business logic layer to manage ID generation.
func (ws *WorldSetting) BeforeCreate(tx *gorm.DB) (err error) {
	// Check if the ID is empty to prevent overwriting IDs that might have been
	// explicitly set (e.g., during a data import or sync operation).
	if ws.ID == "" {
		// Generate a new Version 4 UUID.
		ws.ID = uuid.New().String()
	}
	return
}
