package conflict

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Conflict represents a core narrative engine that drives character actions and plot progression.
// According to REQ-4.1, conflicts can be multi-dimensional (Internal, Interpersonal, External).
// The "Digital Editor" use these tracks to ensure the story maintains constant pressure
// and stays aligned with its central theme.
type Conflict struct {
	// ID is the universally unique identifier for the conflict track.
	ID string `gorm:"primaryKey;type:varchar(36)"`

	// Type classifies the conflict based on narrative theory (e.g., "Man vs. Self", "Man vs. Society").
	// This helps the AI Editor understand the scope and focus of the story's tension.
	Type string `gorm:"type:varchar(100);not null"`

	// Description provides a detailed account of the stakes and the parties involved.
	Description string `gorm:"type:text"`

	// Intensity measures the current level of tension on a scale of 1 to 10.
	// This data is used to generate the "Conflict Map" and pacing audits.
	Intensity int32 `gorm:"type:int"`

	// Resolved indicates whether the conflict has reached its conclusion.
	// Tracking unresolved conflicts helps the writer avoid "forgotten" plot threads.
	Resolved bool `gorm:"type:boolean;default:false"`

	// CreatedAt records when the conflict track was initiated.
	CreatedAt time.Time

	// UpdatedAt tracks the last time the conflict's status or intensity was adjusted.
	UpdatedAt time.Time
}

// BeforeCreate is a GORM lifecycle hook that triggers before a Conflict record is inserted.
//
// Why this approach? Centralizing ID generation in the model layer guarantees consistency
// regardless of how the record is created, ensuring primary keys are always valid UUIDs.
func (c *Conflict) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}
