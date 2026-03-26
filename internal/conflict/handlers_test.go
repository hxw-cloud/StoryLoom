package conflict

import (
	"testing"

	"github.com/go-openapi/swag"

	"github.com/hxw-cloud/StoryLoom/models"
	"github.com/hxw-cloud/StoryLoom/pkg/db"
	conflict_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/conflict"
)

// setupTestEnvironment initializes an isolated in-memory database for Conflict module testing.
//
// Why this approach? TDD mandates that tests must be idempotent and isolated.
// By clearing the table before every test, we prevent cross-test data contamination.
func setupTestEnvironment() {
	db.InitTestDB(&Conflict{})
	db.DB.Exec("DELETE FROM conflicts")
}

// TestHandlePostConflicts verifies that the POST handler correctly maps API input,
// persists it to the database, and returns the expected Swagger model.
func TestHandlePostConflicts(t *testing.T) {
	setupTestEnvironment()

	// 1. Arrange: Define the input payload using Swagger's pointer-based models.
	input := models.ConflictInput{
		Type:        swag.String("Man vs. Self"),
		Description: "The protagonist struggles with their inner shadow.",
		Intensity:   swag.Int32(7),
		Resolved:    false,
	}

	params := conflict_ops.PostConflictsParams{
		Body: &input,
	}

	// 2. Act: Invoke the handler logic.
	responder := HandlePostConflicts(params)

	// 3. Assert: Validate the response and persisted state.
	createdResponder, ok := responder.(*conflict_ops.PostConflictsCreated)
	if !ok {
		t.Fatalf("Expected PostConflictsCreated responder, got %T", responder)
	}

	payload := createdResponder.Payload
	if payload == nil {
		t.Fatal("Expected payload to not be nil")
	}

	// Verify that UUID and timestamps were generated.
	if payload.ID == "" {
		t.Error("Expected ID to be populated by GORM hook")
	}

	// Verify field mapping accuracy.
	if payload.Type != "Man vs. Self" {
		t.Errorf("Expected Type 'Man vs. Self', got '%s'", payload.Type)
	}
}

// TestHandleGetConflicts ensures that the GET handler retrieves all active conflict
// tracks and transforms them into the API response array.
func TestHandleGetConflicts(t *testing.T) {
	setupTestEnvironment()

	// 1. Arrange: Seed the database with sample conflicts.
	db.DB.Create(&Conflict{
		Type:      "Man vs. Nature",
		Intensity: 9,
	})
	db.DB.Create(&Conflict{
		Type:      "Man vs. Machine",
		Intensity: 4,
	})

	params := conflict_ops.GetConflictsParams{}

	// 2. Act: Invoke the GET handler.
	responder := HandleGetConflicts(params)

	// 3. Assert: Verify response integrity and count.
	okResponder, ok := responder.(*conflict_ops.GetConflictsOK)
	if !ok {
		t.Fatalf("Expected GetConflictsOK responder, got %T", responder)
	}

	payload := okResponder.Payload
	if len(payload) != 2 {
		t.Fatalf("Expected 2 records, got %d", len(payload))
	}

	if payload[0].Type != "Man vs. Nature" {
		t.Errorf("Expected first record to be 'Man vs. Nature', got '%s'", payload[0].Type)
	}
}
