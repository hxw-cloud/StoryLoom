package timeline

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/hxw-cloud/StoryLoom/models"
	"github.com/hxw-cloud/StoryLoom/pkg/db"
	timeline_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/timeline"
)

// setupTestEnvironment configures a pristine in-memory database for executing
// timeline handler tests.
//
// Why this approach? TDD necessitates reliable, isolated test cases. An in-memory
// SQLite database combined with a TRUNCATE/DELETE operation ensures no test pollutes
// the global state, avoiding flaky, cascading test failures.
func setupTestEnvironment() {
	db.InitTestDB(&TimelineEvent{})
	db.DB.Exec("DELETE FROM timeline_events")
}

// TestHandlePostTimelineEvents verifies that the POST handler accepts the Swagger API
// payload, translates it into the GORM struct, persists it, and returns the properly
// populated created entity.
func TestHandlePostTimelineEvents(t *testing.T) {
	setupTestEnvironment()

	// 1. Arrange: Construct the input payload
	// Swagger-generated required fields are pointers, so we use swag.String/Int32
	sceneID := strfmt.UUID("scene-999")
	input := models.TimelineEventInput{
		Title:              swag.String("The First Dragon Appears"),
		Description:        "A massive red dragon attacks the northern watchtower.",
		ChronologicalOrder: swag.Int32(100),
		SceneID:            sceneID,
	}

	params := timeline_ops.PostTimelineEventsParams{
		Body: &input,
	}

	// 2. Act: Execute the handler
	responder := HandlePostTimelineEvents(params)

	// 3. Assert: Validate response structure and data mapping
	createdResponder, ok := responder.(*timeline_ops.PostTimelineEventsCreated)
	if !ok {
		t.Fatalf("Expected PostTimelineEventsCreated responder, got %T", responder)
	}

	payload := createdResponder.Payload
	if payload == nil {
		t.Fatal("Expected payload to not be nil")
	}

	// Confirm the BeforeCreate hook successfully injected a new UUID.
	if payload.ID == "" {
		t.Error("Expected ID to be auto-populated by GORM hook")
	}

	// Validate field assignments.
	if payload.Title != "The First Dragon Appears" {
		t.Errorf("Expected Title 'The First Dragon Appears', got '%s'", payload.Title)
	}
	if payload.ChronologicalOrder != 100 {
		t.Errorf("Expected ChronologicalOrder 100, got '%d'", payload.ChronologicalOrder)
	}
}

// TestHandleGetTimelineEvents ensures that the GET handler successfully extracts
// chronological event records from the database and constructs the Swagger output array.
func TestHandleGetTimelineEvents(t *testing.T) {
	setupTestEnvironment()

	// 1. Arrange: Seed the in-memory database
	db.DB.Create(&TimelineEvent{
		Title:              "The Sundering",
		Description:        "The continent splits in two.",
		ChronologicalOrder: 1,
	})
	db.DB.Create(&TimelineEvent{
		Title:              "The Foundation of the Empire",
		Description:        "The first emperor is crowned.",
		ChronologicalOrder: 50,
	})

	params := timeline_ops.GetTimelineEventsParams{}

	// 2. Act: Execute the handler
	responder := HandleGetTimelineEvents(params)

	// 3. Assert: Verify response integrity
	okResponder, ok := responder.(*timeline_ops.GetTimelineEventsOK)
	if !ok {
		t.Fatalf("Expected GetTimelineEventsOK responder, got %T", responder)
	}

	payload := okResponder.Payload

	// We expect exactly 2 records.
	if len(payload) != 2 {
		t.Fatalf("Expected 2 records, got %d", len(payload))
	}

	// Verify chronological sorting (assuming the handler implements order logic)
	// The DB seed was 1 then 50. The handler should return them ordered by ChronologicalOrder asc.
	if payload[0].Title != "The Sundering" {
		t.Errorf("Expected first record Title to be 'The Sundering', got '%s'", payload[0].Title)
	}
}
