package timeline

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TimelineEvent represents a specific, time-bound occurrence within the narrative's universe.
// As per REQ-6.1, this module maintains the global "Master Timeline", crucial for detecting
// logical paradoxes (e.g., a character being in two places at once).
type TimelineEvent struct {
	// ID is the unique identifier for the timeline event.
	ID string `gorm:"primaryKey;type:varchar(36)"`

	// Title provides a quick summary of the event (e.g., "The King's Assassination").
	Title string `gorm:"type:varchar(255);not null"`

	// Description offers comprehensive details about the event's execution and aftermath.
	Description string `gorm:"type:text"`

	// ChronologicalOrder determines the exact sequence of this event relative to others.
	// Unlike Scene order (which is narrative/reading order), this is the in-universe historical order.
	// The AI Editor uses this integer to calculate causality and linear progression.
	ChronologicalOrder int32 `gorm:"type:int;index"`

	// SceneID is an optional foreign key linking this event to a specific narrative Scene.
	// If null, it means the event happens "off-screen" in the world's history.
	SceneID string `gorm:"type:varchar(36);index"`

	// CreatedAt records when the event was added to the timeline.
	CreatedAt time.Time

	// UpdatedAt tracks the last revision for caching and sync mechanisms.
	UpdatedAt time.Time
}

// BeforeCreate is a GORM lifecycle hook that triggers before the event is saved to the database.
//
// Why this approach? Keeping ID generation inside the model layer ensures that no matter how
// the struct is created (API, database seeding script, or test), it always possesses a
// structurally valid, universally unique primary key before hitting the SQLite engine.
func (t *TimelineEvent) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return
}
