package character

import (
	"testing"

	"github.com/go-openapi/swag"

	"github.com/hxw-cloud/StoryLoom/models"
	"github.com/hxw-cloud/StoryLoom/pkg/db"
	character_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/character"
)

// setupTestEnvironment initializes the in-memory SQLite database specifically for
// the Character module tests.
//
// Why this approach? TDD dictates that tests must be completely isolated. By using
// an in-memory database and explicitly purging the `characters` table before every test,
// we eliminate state leakage. If one test fails, it won't be because of residual
// data left over by a previous test.
func setupTestEnvironment() {
	// Initialize the shared memory test database and run schema migrations
	db.InitTestDB(&Character{})
	// Clear the table to ensure a pristine state for the current test run
	db.DB.Exec("DELETE FROM characters")
}

// TestHandlePostCharacters verifies that creating a character correctly
// processes the Swagger input model, persists the data via GORM, and returns
// the populated output model.
func TestHandlePostCharacters(t *testing.T) {
	setupTestEnvironment()

	// 1. Arrange: Define the input payload. We use Swagger's pointer utilities
	// because the generated Swagger model for 'required' fields uses pointers.
	input := models.CharacterInput{
		Name:    swag.String("Elara Vance"),
		Role:    swag.String("Protagonist"),
		PovType: "Third Person Limited",
	}

	params := character_ops.PostCharactersParams{
		Body: &input,
	}

	// 2. Act: Execute the handler
	responder := HandlePostCharacters(params)

	// 3. Assert: Verify the response and database state
	createdResponder, ok := responder.(*character_ops.PostCharactersCreated)
	if !ok {
		t.Fatalf("Expected PostCharactersCreated responder, got %T", responder)
	}

	payload := createdResponder.Payload
	if payload == nil {
		t.Fatal("Expected payload to not be nil")
	}

	// Check if the BeforeCreate hook successfully assigned a UUID
	if payload.ID == "" {
		t.Error("Expected ID to be populated by the GORM BeforeCreate hook")
	}

	// Verify that the data fields map correctly from input to output
	if payload.Name != "Elara Vance" {
		t.Errorf("Expected Name 'Elara Vance', got '%s'", payload.Name)
	}
	if payload.Role != "Protagonist" {
		t.Errorf("Expected Role 'Protagonist', got '%s'", payload.Role)
	}
	if payload.PovType != "Third Person Limited" {
		t.Errorf("Expected POV 'Third Person Limited', got '%s'", payload.PovType)
	}
}

// TestHandleGetCharacters verifies that the GET endpoint successfully queries
// all records from the SQLite database and transforms them accurately into
// the external Swagger response array.
func TestHandleGetCharacters(t *testing.T) {
	setupTestEnvironment()

	// 1. Arrange: Seed the database directly using GORM to simulate existing data.
	db.DB.Create(&Character{
		Name:    "Kaelen",
		Role:    "Antagonist",
		POVType: "First Person",
	})
	db.DB.Create(&Character{
		Name:    "Lyra",
		Role:    "Supporting",
		POVType: "Third Person Objective",
	})

	params := character_ops.GetCharactersParams{}

	// 2. Act: Execute the handler
	responder := HandleGetCharacters(params)

	// 3. Assert: Verify the response array
	okResponder, ok := responder.(*character_ops.GetCharactersOK)
	if !ok {
		t.Fatalf("Expected GetCharactersOK responder, got %T", responder)
	}

	payload := okResponder.Payload
	// We expect exactly 2 records because we seeded 2, and the DB was cleared in setup.
	if len(payload) != 2 {
		t.Fatalf("Expected 2 records, got %d", len(payload))
	}

	// Verify that the data was correctly mapped to the API models
	if payload[0].Name != "Kaelen" {
		t.Errorf("Expected first record to be 'Kaelen', got '%s'", payload[0].Name)
	}
}
