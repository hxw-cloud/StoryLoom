package conflict

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/hxw-cloud/StoryLoom/models"
	"github.com/hxw-cloud/StoryLoom/pkg/db"
	conflict_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/conflict"
)

// HandleGetConflicts processes the GET /conflicts request.
// It queries the database for all conflict tracks and maps them to API models.
//
// Why this approach? Separating internal domain logic from the API contract (DDD)
// allows our database schema to evolve independently of our public interface.
func HandleGetConflicts(params conflict_ops.GetConflictsParams) middleware.Responder {
	var conflicts []Conflict

	// Retrieve all conflicts, sorted by intensity to highlight high-tension areas.
	result := db.DB.Order("intensity desc").Find(&conflicts)
	if result.Error != nil {
		// Implicit failure handling: return empty list on DB error.
		return conflict_ops.NewGetConflictsOK().WithPayload([]*models.Conflict{})
	}

	// Performance Optimization: Pre-allocate capacity for the response slice.
	payload := make([]*models.Conflict, 0, len(conflicts))

	for _, c := range conflicts {
		apiConflict := &models.Conflict{
			ID:          strfmt.UUID(c.ID),
			Type:        c.Type,
			Description: c.Description,
			Intensity:   c.Intensity,
			Resolved:    c.Resolved,
			CreatedAt:   strfmt.DateTime(c.CreatedAt),
		}
		payload = append(payload, apiConflict)
	}

	return conflict_ops.NewGetConflictsOK().WithPayload(payload)
}

// HandlePostConflicts processes the POST /conflicts request.
// It creates a new tracking record for a story conflict and persists it to SQLite.
func HandlePostConflicts(params conflict_ops.PostConflictsParams) middleware.Responder {
	input := params.Body

	// Map external Swagger input to internal GORM domain model.
	// Swag helper functions are used to safely dereference pointers.
	newConflict := Conflict{
		Type:        swag.StringValue(input.Type),
		Description: input.Description,
		Intensity:   swag.Int32Value(input.Intensity),
		Resolved:    input.Resolved,
	}

	// Persist the new conflict to the database via GORM.
	// The BeforeCreate hook will automatically handle UUID generation.
	result := db.DB.Create(&newConflict)
	if result.Error != nil {
		// Return 201 Created with empty payload if DB fails (to be improved).
		return conflict_ops.NewPostConflictsCreated()
	}

	// Map the internal model back to the external representation for the client.
	response := &models.Conflict{
		ID:          strfmt.UUID(newConflict.ID),
		Type:        newConflict.Type,
		Description: newConflict.Description,
		Intensity:   newConflict.Intensity,
		Resolved:    newConflict.Resolved,
		CreatedAt:   strfmt.DateTime(newConflict.CreatedAt),
	}

	return conflict_ops.NewPostConflictsCreated().WithPayload(response)
}
