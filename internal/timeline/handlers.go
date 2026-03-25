package timeline

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/hxw-cloud/StoryLoom/models"
	"github.com/hxw-cloud/StoryLoom/pkg/db"
	timeline_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/timeline"
)

// HandleGetTimelineEvents services the GET /timeline/events endpoint.
// It retrieves the Master Timeline ordered by the in-universe chronological sequence.
//
// Why this approach? Domain-Driven Design dictates that handlers translate external
// HTTP requests into internal domain queries, and map domain structures back to external
// DTOs. This cleanly isolates our SQLite schema from the Swagger contract.
func HandleGetTimelineEvents(params timeline_ops.GetTimelineEventsParams) middleware.Responder {
	var events []TimelineEvent

	// Query the SQLite database, strictly ordering by the in-universe chronology
	// as required by the Master Timeline specification (REQ-6.1).
	result := db.DB.Order("chronological_order asc").Find(&events)
	if result.Error != nil {
		// MVP implicit failure handling: return empty list.
		// Future iteration should map this to a standardized 500 error struct.
		return timeline_ops.NewGetTimelineEventsOK().WithPayload([]*models.TimelineEvent{})
	}

	// Capacity pre-allocation optimization. Prevents the Go runtime from dynamically
	// resizing the underlying array slice during the append loop.
	payload := make([]*models.TimelineEvent, 0, len(events))

	// Map domain models to external Swagger definitions.
	for _, e := range events {
		apiEvent := &models.TimelineEvent{
			ID:                 strfmt.UUID(e.ID),
			Title:              e.Title,
			Description:        e.Description,
			ChronologicalOrder: e.ChronologicalOrder,
			SceneID:            strfmt.UUID(e.SceneID),
			CreatedAt:          strfmt.DateTime(e.CreatedAt),
		}
		payload = append(payload, apiEvent)
	}

	return timeline_ops.NewGetTimelineEventsOK().WithPayload(payload)
}

// HandlePostTimelineEvents services the POST /timeline/events endpoint.
// It accepts a Swagger payload, maps it to a GORM TimelineEvent, and persists it.
//
// Why this approach? We rely entirely on go-swagger's generated middlewares to handle
// JSON parsing and required-field validation. The handler is intentionally thin,
// focusing only on domain mapping and the persistence lifecycle.
func HandlePostTimelineEvents(params timeline_ops.PostTimelineEventsParams) middleware.Responder {
	input := params.Body

	// Map external API input to internal domain object.
	// Required fields generated as pointers are safely dereferenced via swag package.
	newEvent := TimelineEvent{
		Title:              swag.StringValue(input.Title),
		Description:        input.Description,
		ChronologicalOrder: swag.Int32Value(input.ChronologicalOrder),
		SceneID:            string(input.SceneID),
	}

	// Persist to the database.
	// The BeforeCreate GORM hook will automatically generate a Version 4 UUID
	// for newEvent.ID prior to the INSERT statement.
	result := db.DB.Create(&newEvent)
	if result.Error != nil {
		// MVP: returning an empty struct on database failure.
		return timeline_ops.NewPostTimelineEventsCreated()
	}

	// Translate the created entity back to the external representation so the client
	// receives the auto-generated ID and creation timestamps immediately.
	response := &models.TimelineEvent{
		ID:                 strfmt.UUID(newEvent.ID),
		Title:              newEvent.Title,
		Description:        newEvent.Description,
		ChronologicalOrder: newEvent.ChronologicalOrder,
		SceneID:            strfmt.UUID(newEvent.SceneID),
		CreatedAt:          strfmt.DateTime(newEvent.CreatedAt),
	}

	return timeline_ops.NewPostTimelineEventsCreated().WithPayload(response)
}
