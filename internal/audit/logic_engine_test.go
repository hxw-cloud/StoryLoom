package audit

import (
	"context"
	"testing"

	"github.com/hxw-cloud/StoryLoom/internal/character"
	"github.com/hxw-cloud/StoryLoom/internal/scene"
	"github.com/hxw-cloud/StoryLoom/internal/world"
	"github.com/hxw-cloud/StoryLoom/pkg/db"
)

// setupAuditTestEnvironment initializes all required models for cross-module auditing.
func setupAuditTestEnvironment() {
	db.InitTestDB(&world.WorldSetting{}, &character.Character{}, &scene.Scene{})
	db.DB.Exec("DELETE FROM world_settings")
	db.DB.Exec("DELETE FROM characters")
	db.DB.Exec("DELETE FROM scenes")
}

// TestValidateScene_NoPOV verifies that a scene without a POV character is flagged.
func TestValidateScene_NoPOV(t *testing.T) {
	setupAuditTestEnvironment()
	engine := NewLogicEngine()

	// Arrange: Create a scene with no POV
	s := scene.Scene{
		ID:    "scene-1",
		Title: "The Empty Room",
	}
	db.DB.Create(&s)

	// Act
	result, err := engine.ValidateScene(context.Background(), s.ID)
	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsValid {
		t.Error("Expected scene to be invalid due to missing POV")
	}
	found := false
	for _, issue := range result.Issues {
		if issue == "No POV character assigned to the scene." {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected 'No POV' issue not found in: %v", result.Issues)
	}
}

// TestValidateScene_FirstPersonGoal ensures First Person characters have goals.
func TestValidateScene_FirstPersonGoal(t *testing.T) {
	setupAuditTestEnvironment()
	engine := NewLogicEngine()

	// Arrange
	char := character.Character{
		ID:      "char-1",
		Name:    "John",
		POVType: "First Person",
	}
	db.DB.Create(&char)

	s := scene.Scene{
		ID:             "scene-2",
		Title:          "The Soliloquy",
		POVCharacterID: char.ID,
		Goal:           "", // Missing goal
	}
	db.DB.Create(&s)

	// Act
	result, err := engine.ValidateScene(context.Background(), s.ID)
	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsValid {
		t.Error("Expected scene to be invalid due to missing First Person goal")
	}
}
