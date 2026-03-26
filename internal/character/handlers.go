package character

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/hxw-cloud/StoryLoom/models"
	"github.com/hxw-cloud/StoryLoom/pkg/db"
	character_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/character"
)

// HandleGetCharacters processes the GET /characters request with search and camp filtering.
func HandleGetCharacters(params character_ops.GetCharactersParams) middleware.Responder {
	var chars []Character
	query := db.DB

	if params.Camp != nil && *params.Camp != "" {
		query = query.Where("camp = ?", *params.Camp)
	}

	if params.Search != nil && *params.Search != "" {
		s := "%" + *params.Search + "%"
		query = query.Where("name LIKE ? OR background LIKE ?", s, s)
	}

	result := query.Find(&chars)
	if result.Error != nil {
		return character_ops.NewGetCharactersOK().WithPayload([]*models.Character{})
	}

	payload := make([]*models.Character, 0, len(chars))
	for _, c := range chars {
		payload = append(payload, mapToAPI(&c))
	}

	return character_ops.NewGetCharactersOK().WithPayload(payload)
}

// HandlePostCharacters creates a new character.
func HandlePostCharacters(params character_ops.PostCharactersParams) middleware.Responder {
	input := params.Body
	newChar := Character{
		Name:            swag.StringValue(input.Name),
		Age:             int(input.Age),
		Gender:          input.Gender,
		Role:            swag.StringValue(input.Role),
		Camp:            input.Camp,
		Appearance:      input.Appearance,
		Background:      input.Background,
		POVType:         input.PovType,
		ImageURL:        input.ImageURL,
		Want:            input.Want,
		Need:            input.Need,
		PersonaTemplate: input.PersonaTemplate,
	}

	result := db.DB.Create(&newChar)
	if result.Error != nil {
		return character_ops.NewPostCharactersCreated()
	}

	return character_ops.NewPostCharactersCreated().WithPayload(mapToAPI(&newChar))
}

// HandleGetCharactersID returns character details.
func HandleGetCharactersID(params character_ops.GetCharactersIDParams) middleware.Responder {
	var c Character
	if err := db.DB.First(&c, "id = ?", params.ID.String()).Error; err != nil {
		return character_ops.NewGetCharactersIDNotFound()
	}
	return character_ops.NewGetCharactersIDOK().WithPayload(mapToAPI(&c))
}

// HandlePutCharactersID updates character.
func HandlePutCharactersID(params character_ops.PutCharactersIDParams) middleware.Responder {
	var c Character
	if err := db.DB.First(&c, "id = ?", params.ID.String()).Error; err != nil {
		return character_ops.NewPutCharactersIDNotFound()
	}

	input := params.Body
	c.Name = swag.StringValue(input.Name)
	c.Age = int(input.Age)
	c.Gender = input.Gender
	c.Role = swag.StringValue(input.Role)
	c.Camp = input.Camp
	c.Appearance = input.Appearance
	c.Background = input.Background
	c.POVType = input.PovType
	c.ImageURL = input.ImageURL
	c.Want = input.Want
	c.Need = input.Need
	c.PersonaTemplate = input.PersonaTemplate

	db.DB.Save(&c)
	return character_ops.NewPutCharactersIDOK()
}

// HandleGetRelationships returns all relationships for the network map.
func HandleGetRelationships(params character_ops.GetCharactersRelationshipsParams) middleware.Responder {
	var rels []Relationship
	db.DB.Find(&rels)

	payload := make([]*models.Relationship, 0, len(rels))
	for _, r := range rels {
		payload = append(payload, &models.Relationship{
			SourceID:    strfmt.UUID(r.SourceID),
			TargetID:    strfmt.UUID(r.TargetID),
			Type:        r.Type,
			Description: r.Description,
		})
	}
	return character_ops.NewGetCharactersRelationshipsOK().WithPayload(payload)
}

// HandlePostRelationship creates or updates a relationship.
func HandlePostRelationship(params character_ops.PostCharactersRelationshipsParams) middleware.Responder {
	input := params.Body
	rel := Relationship{
		SourceID:    input.SourceID.String(),
		TargetID:    input.TargetID.String(),
		Type:        swag.StringValue(input.Type),
		Description: input.Description,
	}

	db.DB.Save(&rel) // Use Save for upsert behavior
	return character_ops.NewPostCharactersRelationshipsCreated()
}

// HandleGetCharacterArcs returns growth tracking.
func HandleGetCharacterArcs(params character_ops.GetCharactersIDArcsParams) middleware.Responder {
	var arcs []CharacterArc
	db.DB.Where("character_id = ?", params.ID.String()).Find(&arcs)

	payload := make([]*models.CharacterArc, 0, len(arcs))
	for _, a := range arcs {
		payload = append(payload, &models.CharacterArc{
			CharacterID:    strfmt.UUID(a.CharacterID),
			PlotCardID:     strfmt.UUID(a.PlotCardID),
			StateChange:    a.StateChange,
			InternalGrowth: int32(a.InternalGrowth),
		})
	}
	return character_ops.NewGetCharactersIDArcsOK().WithPayload(payload)
}

func mapToAPI(c *Character) *models.Character {
	return &models.Character{
		ID:              strfmt.UUID(c.ID),
		Name:            c.Name,
		Age:             int64(c.Age),
		Gender:          c.Gender,
		Role:            c.Role,
		Camp:            c.Camp,
		Appearance:      c.Appearance,
		Background:      c.Background,
		PovType:         c.POVType,
		ImageURL:        c.ImageURL,
		Want:            c.Want,
		Need:            c.Need,
		PersonaTemplate: c.PersonaTemplate,
		CreatedAt:       strfmt.DateTime(c.CreatedAt),
	}
}
