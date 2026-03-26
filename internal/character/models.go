package character

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Character represents a person, entity, or actor within the story's universe.
// Enhanced for REQ-2.1 to REQ-2.6.
type Character struct {
	ID string `gorm:"primaryKey;type:varchar(36)"`

	// Basic Info (REQ-2.1)
	Name       string `gorm:"type:varchar(200);not null;index"`
	Age        int    `gorm:"type:int"`
	Gender     string `gorm:"type:varchar(50)"`
	Role       string `gorm:"type:varchar(100);not null"`
	Camp       string `gorm:"type:varchar(100);index"` // Faction/Organization
	Appearance string `gorm:"type:text"`
	Background string `gorm:"type:text"`
	ImageURL   string `gorm:"type:varchar(500)"`

	// POV Setting (REQ-2.1)
	POVType string `gorm:"type:varchar(100)"`

	// Motivation & Conflict (REQ-2.2)
	Want string `gorm:"type:text"` // External goal
	Need string `gorm:"type:text"` // Internal need

	// Persona (REQ-2.5)
	PersonaTemplate string `gorm:"type:text"` // Vocabulary/Voice style

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Relationship tracks connections between characters (REQ-2.4).
type Relationship struct {
	SourceID    string `gorm:"primaryKey;type:varchar(36)"`
	TargetID    string `gorm:"primaryKey;type:varchar(36)"`
	Type        string `gorm:"type:varchar(100);not null"` // Blood, Love, Rival, etc.
	Description string `gorm:"type:text"`
	UpdatedAt   time.Time
}

// CharacterArc tracks growth across story beats (REQ-2.3).
type CharacterArc struct {
	ID             string `gorm:"primaryKey;type:varchar(36)"`
	CharacterID    string `gorm:"type:varchar(36);index;not null"`
	PlotCardID     string `gorm:"type:varchar(36);index;not null"`
	StateChange    string `gorm:"type:text"`
	InternalGrowth int    `gorm:"type:int"` // -100 to 100
}

func (c *Character) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}

func (ca *CharacterArc) BeforeCreate(tx *gorm.DB) (err error) {
	if ca.ID == "" {
		ca.ID = uuid.New().String()
	}
	return
}
