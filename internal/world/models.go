package world

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WorldSetting represents a foundational rule, concept, or element of the story's world.
// Enhanced for REQ-1.1 and REQ-1.2.
type WorldSetting struct {
	ID          string `gorm:"primaryKey;type:varchar(36)"`
	Category    string `gorm:"type:varchar(100);not null"`
	Name        string `gorm:"type:varchar(200);not null;index"`
	Description string `gorm:"type:text"`
	LogicRules  string `gorm:"type:text"`

	// Tags for multi-dimensional management (REQ-1.1)
	// Stored as a comma-separated string for SQLite simplicity, but can be split in logic.
	Tags string `gorm:"type:text"`

	// ImageURL for reference (REQ-1.1)
	ImageURL string `gorm:"type:varchar(500)"`

	// ParentID for infinite nested classification (REQ-1.2)
	ParentID string `gorm:"type:varchar(36);index"`

	// UsageCount for Iceberg Theory tracking (REQ-1.4)
	UsageCount int `gorm:"type:int;default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// HistoricalEvent records significant past occurrences (REQ-1.3).
type HistoricalEvent struct {
	ID          string `gorm:"primaryKey;type:varchar(36)"`
	Title       string `gorm:"type:varchar(255);not null"`
	EventTime   string `gorm:"type:varchar(100)"`
	ImpactScope string `gorm:"type:text"`

	// InvolvedCharacters stored as comma-separated IDs
	InvolvedCharacters string `gorm:"type:text"`

	Cause        string `gorm:"type:text"`
	Effect       string `gorm:"type:text"`
	IsIcebergTip bool   `gorm:"type:boolean;default:false"`

	CreatedAt time.Time
}

// WorldTemplate represents a pre-defined professional world-building element (REQ-1.4).
type WorldTemplate struct {
	ID             string `gorm:"primaryKey;type:varchar(36)"`
	Category       string `gorm:"type:varchar(100);not null"`
	Name           string `gorm:"type:varchar(200);not null"`
	Description    string `gorm:"type:text"`
	SuggestedLogic string `gorm:"type:text"`
}

func (ws *WorldSetting) BeforeCreate(tx *gorm.DB) (err error) {
	if ws.ID == "" {
		ws.ID = uuid.New().String()
	}
	return
}

func (he *HistoricalEvent) BeforeCreate(tx *gorm.DB) (err error) {
	if he.ID == "" {
		he.ID = uuid.New().String()
	}
	return
}
