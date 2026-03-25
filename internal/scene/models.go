package scene

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Scene represents a single continuous narrative event in the story.
// Scenes are the granular building blocks attached to PlotCards (REQ-5.1).
// The "Digital Editor" validates these scenes against the World rules and the
// POV Character's constraints (e.g., ensuring a First Person POV doesn't describe
// things they can't perceive).
type Scene struct {
	// ID is the unique identifier for the scene.
	ID string `gorm:"primaryKey;type:varchar(36)"`

	// Title is a brief descriptor for the scene (e.g., "The Cafe Meeting").
	Title string `gorm:"type:varchar(200);not null"`

	// PlotCardID links this scene to a broader narrative beat.
	// This establishes the many-to-one relationship between Scenes and PlotCards.
	PlotCardID string `gorm:"type:varchar(36);not null;index"`

	// POVCharacterID designates whose perspective this scene is experienced through.
	// The AI Logic Editor uses this to enforce narrative consistency based on that
	// specific character's known information and sensory limits.
	POVCharacterID string `gorm:"type:varchar(36);index"`

	// Goal defines what the POV character is trying to achieve within this specific scene.
	// Establishing clear goals is a foundational principle of scene structure.
	Goal string `gorm:"type:text"`

	// Conflict describes the obstacle, person, or internal struggle preventing
	// the POV character from easily achieving their Goal.
	Conflict string `gorm:"type:text"`

	// Resolution dictates how the scene concludes (e.g., "Yes, but...", "No, and...").
	// This drives the pacing and transitions into the next scene.
	Resolution string `gorm:"type:text"`

	// CreatedAt timestamps the creation of the scene.
	CreatedAt time.Time

	// UpdatedAt timestamps the last modification.
	UpdatedAt time.Time
}

// BeforeCreate is a GORM hook that executes prior to inserting a new Scene record.
//
// Why this approach? By abstracting UUID generation into the model layer, we ensure
// that every Scene gets a valid identifier immediately, safely decoupled from the HTTP handlers.
func (s *Scene) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return
}
