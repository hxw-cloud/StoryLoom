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
func (e *LogicEngine) ValidateScene(ctx context.Context, sceneID string) (*AuditResult, error) {
	var s scene.Scene
	if err := db.DB.First(&s, "id = ?", sceneID).Error; err != nil {
		return nil, fmt.Errorf("scene not found: %w", err)
	}

	result := &AuditResult{IsValid: true, Issues: []string{}}

	// 1. Validate POV existence and consistency
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
	// Rule: Every scene should have an assigned POV character
	if s.POVCharacterID == "" {
		result.Issues = append(result.Issues, "No POV character assigned to the scene.")
		return nil
	}

	var char character.Character
	// Rule: The assigned POV character MUST exist in the database (Reference integrity)
	if err := db.DB.First(&char, "id = ?", s.POVCharacterID).Error; err != nil {
		result.Issues = append(result.Issues, fmt.Sprintf("POV Character ID '%s' does not exist in your character library.", s.POVCharacterID))
		return nil
	}

	// Rule: First Person POV heuristic check
	if char.POVType == "First Person" {
		if s.Goal == "" {
			result.Issues = append(result.Issues, fmt.Sprintf("Character %s (First Person) needs a clearly defined internal goal for this scene.", char.Name))
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

	// Dynamic rule checking based on world settings
	for _, setting := range settings {
		// Example: If a setting is categorized as Magic and has specific logic rules
		if strings.Contains(strings.ToLower(setting.Category), "magic") {
			rules := strings.ToLower(setting.LogicRules)
			// Heuristic: if world rule mentions "cost" or "sacrifice", ensure the scene mentions it.
			if strings.Contains(rules, "cost") || strings.Contains(rules, "limit") {
				content := strings.ToLower(s.Conflict + " " + s.Resolution)
				if !strings.Contains(content, "cost") && !strings.Contains(content, "spent") && !strings.Contains(content, "drain") {
					// We add a warning (non-blocking) for the writer to consider
					result.Issues = append(result.Issues, fmt.Sprintf("World Rule '%s' defines magic limitations. Ensure this scene's conflict reflects the cost of using magic.", setting.Name))
				}
			}
		}
	}

	return nil
}
