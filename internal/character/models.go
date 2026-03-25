package character

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Character represents a person, entity, or actor within the story's universe.
// In StoryLoom, characters are deeply integrated into the structural logic of the novel.
// They are tied to specific Point of View (POV) constraints, scenes, and relationships.
// The "Digital Editor" will use this metadata (like POV type) to validate whether a scene
// written from this character's perspective contains information they logically shouldn't know.
type Character struct {
	// ID is the unique UUID for the character, ensuring global uniqueness
	// across offline syncs or collaborative environments.
	ID string `gorm:"primaryKey;type:varchar(36)"`

	// Name is the primary designation or appellation used for the character.
	// We index this field to provide fast autocomplete suggestions when authors
	// mention characters in the plot cards or scene text.
	Name string `gorm:"type:varchar(200);not null;index"`

	// Role designates the structural narrative purpose of the character
	// (e.g., Protagonist, Antagonist, Supporting, Mentor).
	// This helps the logic engine evaluate character arcs and screen time.
	Role string `gorm:"type:varchar(100);not null"`

	// POVType determines the default narrative perspective if this character
	// is the focal point of a scene. Examples: "First Person", "Third Person Limited".
	// The AI Digital Editor uses this to flag sensory or informational violations
	// (e.g., a "First Person" character cannot see what happens behind a closed door).
	POVType string `gorm:"type:varchar(100)"`

	// CreatedAt timestamps the creation of the character, useful for tracking timeline.
	CreatedAt time.Time

	// UpdatedAt timestamps the last modification for caching and synchronization.
	UpdatedAt time.Time
}

// BeforeCreate is a GORM lifecycle hook that triggers right before a new Character
// record is inserted into the SQLite database.
//
// Why this approach? We manage UUID generation at the model layer rather than the
// handler or database layer. This ensures that even if we create characters programmatically
// via scripts or data imports (bypassing the API), they will still receive valid UUIDs.
func (c *Character) BeforeCreate(tx *gorm.DB) (err error) {
	// We only generate a UUID if one has not already been provided.
	// This flexibility allows for deterministic ID assignment during migrations or testing.
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}
