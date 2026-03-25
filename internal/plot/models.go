package plot

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PlotCard represents a distinct narrative beat or structural element within the story.
// It acts as the core building block for the dynamic outlining system described in REQ-3.1.
// Writers use these cards to sequence events, track pacing, and measure conflict intensity.
// The "Digital Editor" utilizes the order and intensity fields to generate a pacing map.
type PlotCard struct {
	// ID is a universally unique identifier for the PlotCard.
	// We use UUIDs to prevent collisions during offline synchronization
	// and to ensure seamless collaboration between multiple writers.
	ID string `gorm:"primaryKey;type:varchar(36)"`

	// Title is the brief, descriptive name of the narrative beat (e.g., "The Inciting Incident").
	Title string `gorm:"type:varchar(200);not null"`

	// Description provides the detailed content of what occurs during this plot beat.
	// It allows writers to sketch out the scene before writing the actual manuscript.
	Description string `gorm:"type:text"`

	// ConflictIntensity measures the dramatic tension of the plot card on a scale of 1 to 5.
	// This numeric value is heavily utilized by the "Digital Editor" (REQ-3.3) to dynamically
	// draw pacing charts and warn the writer if the story lacks high-tension climaxes or
	// suffers from prolonged periods of low intensity.
	ConflictIntensity int32 `gorm:"type:int"`

	// OrderIndex strictly dictates the chronological sequence of this plot card
	// within the global structure of the narrative. It enables the drag-and-drop
	// reordering functionality in the user interface.
	OrderIndex int32 `gorm:"type:int;index"`

	// CreatedAt records the timestamp of creation for auditing purposes.
	CreatedAt time.Time

	// UpdatedAt tracks the last modification time, crucial for syncing changes across clients.
	UpdatedAt time.Time
}

// BeforeCreate is a GORM lifecycle hook that triggers before a PlotCard is inserted
// into the SQLite database.
//
// Why this approach? Handling UUID generation internally within the model layer
// guarantees that every PlotCard receives a valid identifier, regardless of whether
// it was created via the REST API, a script, or a database migration tool.
func (p *PlotCard) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}
