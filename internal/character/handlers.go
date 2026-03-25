package character

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/hxw-cloud/StoryLoom/models"
	"github.com/hxw-cloud/StoryLoom/pkg/db"
	character_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/character"
)

// HandleGetCharacters processes the GET /characters request.
// It queries the SQLite database for all Character records and maps them
// to the external Swagger definitions to be returned as JSON.
//
// Why this approach? Translating between the internal `Character` struct and the
// external `models.Character` struct enforces Domain-Driven Design (DDD). It isolates
// our database schema from our API contract, allowing either to change independently
// without breaking the other.
func HandleGetCharacters(params character_ops.GetCharactersParams) middleware.Responder {
	var chars []Character

	// Fetch all characters from the database.
	// TODO: For production, implement pagination (limit/offset) to prevent massive memory spikes
	// if a user creates thousands of characters for an epic fantasy series.
	result := db.DB.Find(&chars)
	if result.Error != nil {
		// Implicit error handling: return an empty list if query fails.
		// In a mature API, this should return a structured 500 Internal Server Error model.
		return character_ops.NewGetCharactersOK().WithPayload([]*models.Character{})
	}

	// Pre-allocate the payload slice with the exact capacity needed.
	// This is a Go performance optimization that prevents the underlying array
	// from needing to resize and copy memory during the append loop.
	payload := make([]*models.Character, 0, len(chars))

	for _, c := range chars {
		apiChar := &models.Character{
			ID:        strfmt.UUID(c.ID),
			Name:      c.Name,
			Role:      c.Role,
			PovType:   c.POVType,
			CreatedAt: strfmt.DateTime(c.CreatedAt),
		}
		payload = append(payload, apiChar)
	}

	return character_ops.NewGetCharactersOK().WithPayload(payload)
}

// HandlePostCharacters processes the POST /characters request.
// It accepts a Swagger model, translates it into our internal GORM model,
// saves it to SQLite, and returns the newly created record.
//
// Why this approach? We rely on go-swagger's generated middleware to enforce
// JSON structure and basic required field validation. The handler's only responsibility
// is domain translation and persistence.
func HandlePostCharacters(params character_ops.PostCharactersParams) middleware.Responder {
	input := params.Body

	// Map external API input to internal domain model.
	// Required string pointers are safely dereferenced using swag.StringValue.
	// Optional fields are passed directly.
	newChar := Character{
		Name:    swag.StringValue(input.Name),
		Role:    swag.StringValue(input.Role),
		POVType: input.PovType,
	}

	// Save the record to the database.
	// The DB will invoke the BeforeCreate hook on the Character model to automatically
	// generate and assign a V4 UUID to newChar.ID before insertion.
	result := db.DB.Create(&newChar)
	if result.Error != nil {
		// Logically, if the DB fails to save, we shouldn't return a 201 Created.
		// For the MVP, we will return an empty created payload, but this must be updated
		// to an error response once standardized error models are added to swagger.yaml.
		return character_ops.NewPostCharactersCreated()
	}

	// Map the internal model back to the external representation so the client
	// receives the auto-generated ID and CreatedAt timestamps.
	response := &models.Character{
		ID:        strfmt.UUID(newChar.ID),
		Name:      newChar.Name,
		Role:      newChar.Role,
		PovType:   newChar.POVType,
		CreatedAt: strfmt.DateTime(newChar.CreatedAt),
	}

	return character_ops.NewPostCharactersCreated().WithPayload(response)
}
