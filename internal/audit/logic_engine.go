package audit

import (
	"context"
	"fmt"
	"strings"

	"github.com/hxw-cloud/StoryLoom/internal/character"
	"github.com/hxw-cloud/StoryLoom/internal/scene"
	"github.com/hxw-cloud/StoryLoom/internal/world"
	"github.com/hxw-cloud/StoryLoom/pkg/db"
)

// LogicEngine is the core "Digital Editor" service that evaluates narrative consistency.
// It analyzes scenes against character motivations and world-building rules.
type LogicEngine struct{}

// NewLogicEngine initializes a new instance of the Digital Editor engine.
func NewLogicEngine() *LogicEngine {
	return &LogicEngine{}
}

// AuditResult encapsulates the findings of a logic validation pass.
type AuditResult struct {
	IsValid bool
	Issues  []string
}

// ValidateScene performs a comprehensive consistency check on a specific scene.
// It cross-references the scene's content and metadata with character POV rules
// and world-building constraints.
func (e *LogicEngine) ValidateScene(ctx context.Context, sceneID string) (*AuditResult, error) {
	var s scene.Scene
	if err := db.DB.First(&s, "id = ?", sceneID).Error; err != nil {
		return nil, fmt.Errorf("scene not found: %w", err)
	}

	result := &AuditResult{IsValid: true, Issues: []string{}}

	// 1. Validate POV consistency
	if err := e.checkPOVConsistency(&s, result); err != nil {
		return nil, err
	}

	// 2. Validate World Rules (Logic Settings)
	if err := e.checkWorldRules(&s, result); err != nil {
		return nil, err
	}

	if len(result.Issues) > 0 {
		result.IsValid = false
	}

	return result, nil
}

// checkPOVConsistency ensures the scene respects the limitations of the chosen POV character.
func (e *LogicEngine) checkPOVConsistency(s *scene.Scene, result *AuditResult) error {
	if s.POVCharacterID == "" {
		result.Issues = append(result.Issues, "No POV character assigned to the scene.")
		return nil
	}

	var char character.Character
	if err := db.DB.First(&char, "id = ?", s.POVCharacterID).Error; err != nil {
		return fmt.Errorf("POV character not found: %w", err)
	}

	// Rule: First Person POV should focus on internal experience and immediate perception.
	if char.POVType == "First Person" {
		// This is a heuristic check for the "Digital Editor" prototype.
		// In a real scenario, this would involve LLM analysis of the scene text.
		if s.Goal == "" {
			result.Issues = append(result.Issues, fmt.Sprintf("Character %s (First Person) needs a clearly defined internal goal.", char.Name))
		}
	}

	return nil
}

// checkWorldRules evaluates if the scene's resolution or conflict violates world-building logic.
func (e *LogicEngine) checkWorldRules(s *scene.Scene, result *AuditResult) error {
	var settings []world.WorldSetting
	if err := db.DB.Find(&settings).Error; err != nil {
		return fmt.Errorf("failed to fetch world settings: %w", err)
	}

	// Prototype Rule: If a world setting mentions "Magic" and scene conflict doesn't,
	// warn the writer if this is a high-fantasy context (simple keyword matching for now).
	for _, setting := range settings {
		if strings.Contains(strings.ToLower(setting.Category), "magic") {
			if strings.Contains(strings.ToLower(setting.LogicRules), "cost") && !strings.Contains(strings.ToLower(s.Conflict), "cost") {
				// result.Issues = append(result.Issues, fmt.Sprintf("World Rule '%s' requires a cost for magic, but none is mentioned in this scene's conflict.", setting.Name))
			}
		}
	}

	return nil
}
