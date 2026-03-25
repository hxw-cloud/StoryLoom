package plot

import (
	"testing"

	"github.com/go-openapi/swag"

	"github.com/hxw-cloud/StoryLoom/models"
	"github.com/hxw-cloud/StoryLoom/pkg/db"
	plot_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/plot"
)

// setupTestEnvironment configures an isolated in-memory SQLite database specifically
// for executing the PlotCard handler tests.
//
// Why this approach? TDD requires strict isolation. By utilizing an ephemeral in-memory
// database and explicitly truncating the `plot_cards` table before every test run,
// we guarantee that tests do not interfere with one another (no state leakage).
func setupTestEnvironment() {
	db.InitTestDB(&PlotCard{})
	// Truncate the table to prevent data from a previous test affecting the current one.
	db.DB.Exec("DELETE FROM plot_cards")
}

// TestHandlePostPlots ensures that the HTTP POST handler successfully processes
// incoming API payloads, transforms them into the internal domain model, saves them
// via GORM, and returns a correctly formatted Swagger response.
func TestHandlePostPlots(t *testing.T) {
	setupTestEnvironment()

	// 1. Arrange: Construct the input payload exactly as the Swagger API expects it.
	// Generated Swagger models use pointers for required fields to handle nullability.
	input := models.PlotCardInput{
		Title:             swag.String("The Hero's Journey Begins"),
		Description:       "The protagonist receives the call to adventure.",
		ConflictIntensity: swag.Int32(4),
		OrderIndex:        1,
	}

	params := plot_ops.PostPlotsParams{
		Body: &input,
	}

	// 2. Act: Execute the handler logic
	responder := HandlePostPlots(params)

	// 3. Assert: Validate the handler response and the persisted data
	createdResponder, ok := responder.(*plot_ops.PostPlotsCreated)
	if !ok {
		t.Fatalf("Expected PostPlotsCreated responder, got %T", responder)
	}

	payload := createdResponder.Payload
	if payload == nil {
		t.Fatal("Expected payload to not be nil")
	}

	// Verify that the database model hook successfully generated a UUID.
	if payload.ID == "" {
		t.Error("Expected ID to be populated by the GORM BeforeCreate hook")
	}

	// Confirm that the data attributes were mapped correctly.
	if payload.Title != "The Hero's Journey Begins" {
		t.Errorf("Expected Title 'The Hero's Journey Begins', got '%s'", payload.Title)
	}
	if payload.ConflictIntensity != 4 {
		t.Errorf("Expected ConflictIntensity 4, got '%d'", payload.ConflictIntensity)
	}
}

// TestHandleGetPlots guarantees that the HTTP GET handler successfully queries
// the database for all plot cards and maps them accurately to the Swagger API format.
func TestHandleGetPlots(t *testing.T) {
	setupTestEnvironment()

	// 1. Arrange: Seed the database directly to simulate a pre-existing state.
	db.DB.Create(&PlotCard{
		Title:             "Inciting Incident",
		Description:       "The meteor strikes the city.",
		ConflictIntensity: 5,
		OrderIndex:        1,
	})
	db.DB.Create(&PlotCard{
		Title:             "Rising Action",
		Description:       "The survivors gather supplies.",
		ConflictIntensity: 2,
		OrderIndex:        2,
	})

	params := plot_ops.GetPlotsParams{}

	// 2. Act: Invoke the GET handler
	responder := HandleGetPlots(params)

	// 3. Assert: Verify the response size and content mapping
	okResponder, ok := responder.(*plot_ops.GetPlotsOK)
	if !ok {
		t.Fatalf("Expected GetPlotsOK responder, got %T", responder)
	}

	payload := okResponder.Payload

	// We expect exactly 2 plot cards to be returned based on our seeding.
	if len(payload) != 2 {
		t.Fatalf("Expected 2 records, got %d", len(payload))
	}

	// Validate the mapping logic from internal GORM models to Swagger structs.
	if payload[0].Title != "Inciting Incident" {
		t.Errorf("Expected first record Title to be 'Inciting Incident', got '%s'", payload[0].Title)
	}
}
