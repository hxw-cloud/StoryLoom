package world

import (
	"testing"

	"github.com/go-openapi/swag"

	"github.com/hxw-cloud/StoryLoom/models"
	"github.com/hxw-cloud/StoryLoom/pkg/db"
	world_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/world"
)

// setupTestEnvironment initializes the in-memory database and runs auto-migrations
// for the WorldSetting model. It returns a cleanup function (though memory DBs
// clean themselves, it's good practice for clearing global state).
//
// Why this approach? It ensures every test suite starts with a fresh, correctly
// structured database, isolating test logic and preventing state leakage.
func setupTestEnvironment() {
	db.InitTestDB(&WorldSetting{})
	// Ensure the table is empty before each test to prevent state leakage
	// because we are using a shared in-memory SQLite connection.
	db.DB.Exec("DELETE FROM world_settings")
}

// TestHandlePostSettings verifies that the POST handler correctly takes API input,
// persists it to the database via GORM, and returns the expected Swagger model
// with auto-generated fields like ID and CreatedAt.
func TestHandlePostSettings(t *testing.T) {
	setupTestEnvironment()

	// Define the test input payload.
	input := models.WorldSettingInput{
		Category:    swag.String("Magic System"),
		Name:        swag.String("Equivalent Exchange"),
		Description: "To obtain, something of equal value must be lost.",
		LogicRules:  "Mass must be conserved.",
	}

	params := world_ops.PostWorldSettingsParams{
		Body: &input,
	}

	// Execute the handler
	responder := HandlePostSettings(params)

	// Assert the response type
	createdResponder, ok := responder.(*world_ops.PostWorldSettingsCreated)
	if !ok {
		t.Fatalf("Expected PostWorldSettingsCreated responder, got %T", responder)
	}

	payload := createdResponder.Payload
	if payload == nil {
		t.Fatal("Expected payload to not be nil")
	}

	// Validate that auto-generated fields are populated
	if payload.ID == "" {
		t.Error("Expected ID to be populated by the GORM BeforeCreate hook")
	}

	// Validate that input fields were correctly mapped and saved
	if payload.Category != "Magic System" {
		t.Errorf("Expected Category 'Magic System', got '%s'", payload.Category)
	}
	if payload.Name != "Equivalent Exchange" {
		t.Errorf("Expected Name 'Equivalent Exchange', got '%s'", payload.Name)
	}
}

// TestHandleGetSettings verifies that the GET handler retrieves all records
// from the database and maps them correctly back to the Swagger API format.
func TestHandleGetSettings(t *testing.T) {
	setupTestEnvironment()

	// Seed the test database directly using GORM
	db.DB.Create(&WorldSetting{
		Category:   "Geography",
		Name:       "The Shattered Isles",
		LogicRules: "Travel by sea is impossible without an Airship.",
	})
	db.DB.Create(&WorldSetting{
		Category:   "Technology",
		Name:       "Steam Cores",
		LogicRules: "Requires water and a heat source.",
	})

	params := world_ops.GetWorldSettingsParams{}

	// Execute the handler
	responder := HandleGetSettings(params)

	// Assert the response type
	okResponder, ok := responder.(*world_ops.GetWorldSettingsOK)
	if !ok {
		t.Fatalf("Expected GetWorldSettingsOK responder, got %T", responder)
	}

	payload := okResponder.Payload
	if len(payload) != 2 {
		t.Fatalf("Expected 2 records, got %d", len(payload))
	}

	// Validate that the mapping was successful
	if payload[0].Category != "Geography" {
		t.Errorf("Expected first record to be 'Geography', got '%s'", payload[0].Category)
	}
}
