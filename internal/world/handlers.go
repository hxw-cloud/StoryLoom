package world

import (
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/hxw-cloud/StoryLoom/models"
	"github.com/hxw-cloud/StoryLoom/pkg/db"
	world_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/world"
)

// HandleGetSettings processes the GET request for retrieving world settings with search and filtering.
func HandleGetSettings(params world_ops.GetWorldSettingsParams) middleware.Responder {
	var settings []WorldSetting
	query := db.DB

	// Filtering by Category (REQ-1.1)
	if params.Category != nil && *params.Category != "" {
		query = query.Where("category = ?", *params.Category)
	}

	// Filtering by Tag (REQ-1.1)
	if params.Tag != nil && *params.Tag != "" {
		query = query.Where("tags LIKE ?", "%"+*params.Tag+"%")
	}

	// Full-text Search (REQ-1.1)
	if params.Search != nil && *params.Search != "" {
		s := "%" + *params.Search + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", s, s)
	}

	result := query.Find(&settings)
	if result.Error != nil {
		return world_ops.NewGetWorldSettingsOK().WithPayload([]*models.WorldSetting{})
	}

	payload := make([]*models.WorldSetting, 0, len(settings))
	for _, s := range settings {
		tags := strings.Split(s.Tags, ",")
		if s.Tags == "" {
			tags = []string{}
		}

		apiSetting := &models.WorldSetting{
			ID:          strfmt.UUID(s.ID),
			Category:    s.Category,
			Name:        s.Name,
			Description: s.Description,
			LogicRules:  s.LogicRules,
			Tags:        tags,
			ImageURL:    s.ImageURL,
			ParentID:    strfmt.UUID(s.ParentID),
			UsageCount:  int64(s.UsageCount),
			CreatedAt:   strfmt.DateTime(s.CreatedAt),
		}
		payload = append(payload, apiSetting)
	}

	return world_ops.NewGetWorldSettingsOK().WithPayload(payload)
}

// HandlePostSettings processes the POST request for creating a new world setting.
func HandlePostSettings(params world_ops.PostWorldSettingsParams) middleware.Responder {
	input := params.Body
	tags := strings.Join(input.Tags, ",")

	newSetting := WorldSetting{
		Category:    swag.StringValue(input.Category),
		Name:        swag.StringValue(input.Name),
		Description: input.Description,
		LogicRules:  input.LogicRules,
		Tags:        tags,
		ImageURL:    input.ImageURL,
		ParentID:    string(input.ParentID),
	}

	result := db.DB.Create(&newSetting)
	if result.Error != nil {
		return world_ops.NewPostWorldSettingsCreated()
	}

	response := &models.WorldSetting{
		ID:          strfmt.UUID(newSetting.ID),
		Category:    newSetting.Category,
		Name:        newSetting.Name,
		Description: newSetting.Description,
		LogicRules:  newSetting.LogicRules,
		Tags:        input.Tags,
		ImageURL:    newSetting.ImageURL,
		ParentID:    strfmt.UUID(newSetting.ParentID),
		CreatedAt:   strfmt.DateTime(newSetting.CreatedAt),
	}

	return world_ops.NewPostWorldSettingsCreated().WithPayload(response)
}

// HandleGetSettingsID returns a single setting.
func HandleGetSettingsID(params world_ops.GetWorldSettingsIDParams) middleware.Responder {
	var s WorldSetting
	if err := db.DB.First(&s, "id = ?", params.ID.String()).Error; err != nil {
		return world_ops.NewGetWorldSettingsIDNotFound()
	}

	tags := strings.Split(s.Tags, ",")
	if s.Tags == "" {
		tags = []string{}
	}

	return world_ops.NewGetWorldSettingsIDOK().WithPayload(&models.WorldSetting{
		ID:          strfmt.UUID(s.ID),
		Category:    s.Category,
		Name:        s.Name,
		Description: s.Description,
		LogicRules:  s.LogicRules,
		Tags:        tags,
		ImageURL:    s.ImageURL,
		ParentID:    strfmt.UUID(s.ParentID),
		UsageCount:  int64(s.UsageCount),
		CreatedAt:   strfmt.DateTime(s.CreatedAt),
	})
}

// HandlePutSettingsID updates an existing setting.
func HandlePutSettingsID(params world_ops.PutWorldSettingsIDParams) middleware.Responder {
	input := params.Body
	tags := strings.Join(input.Tags, ",")

	var s WorldSetting
	if err := db.DB.First(&s, "id = ?", params.ID.String()).Error; err != nil {
		return world_ops.NewPutWorldSettingsIDNotFound()
	}

	s.Category = swag.StringValue(input.Category)
	s.Name = swag.StringValue(input.Name)
	s.Description = input.Description
	s.LogicRules = input.LogicRules
	s.Tags = tags
	s.ImageURL = input.ImageURL
	s.ParentID = string(input.ParentID)

	db.DB.Save(&s)
	return world_ops.NewPutWorldSettingsIDOK()
}

// HandleGetHistory returns historical events.
func HandleGetHistory(params world_ops.GetWorldHistoryParams) middleware.Responder {
	var events []HistoricalEvent
	db.DB.Find(&events)

	payload := make([]*models.HistoricalEvent, 0, len(events))
	for _, e := range events {
		chars := strings.Split(e.InvolvedCharacters, ",")
		if e.InvolvedCharacters == "" {
			chars = []string{}
		}

		payload = append(payload, &models.HistoricalEvent{
			ID:                 strfmt.UUID(e.ID),
			Title:              e.Title,
			EventTime:          e.EventTime,
			ImpactScope:        e.ImpactScope,
			InvolvedCharacters: chars,
			Cause:              e.Cause,
			Effect:             e.Effect,
			IsIcebergTip:       e.IsIcebergTip,
			CreatedAt:          strfmt.DateTime(e.CreatedAt),
		})
	}
	return world_ops.NewGetWorldHistoryOK().WithPayload(payload)
}

// HandlePostHistory creates a historical event.
func HandlePostHistory(params world_ops.PostWorldHistoryParams) middleware.Responder {
	input := params.Body
	chars := strings.Join(input.InvolvedCharacters, ",")

	e := HistoricalEvent{
		Title:              swag.StringValue(input.Title),
		EventTime:          swag.StringValue(input.EventTime),
		ImpactScope:        input.ImpactScope,
		InvolvedCharacters: chars,
		Cause:              input.Cause,
		Effect:             input.Effect,
		IsIcebergTip:       input.IsIcebergTip,
	}

	db.DB.Create(&e)
	return world_ops.NewPostWorldHistoryCreated().WithPayload(&models.HistoricalEvent{
		ID:                 strfmt.UUID(e.ID),
		Title:              e.Title,
		EventTime:          e.EventTime,
		ImpactScope:        e.ImpactScope,
		InvolvedCharacters: input.InvolvedCharacters,
		Cause:              e.Cause,
		Effect:             e.Effect,
		IsIcebergTip:       e.IsIcebergTip,
		CreatedAt:          strfmt.DateTime(e.CreatedAt),
	})
}

// HandleGetAudit returns usage stats (REQ-1.4).
func HandleGetAudit(params world_ops.GetWorldAuditParams) middleware.Responder {
	var settings []WorldSetting
	db.DB.Find(&settings)

	intensityMap := make(map[string]int64)
	usedCount := 0
	for _, s := range settings {
		intensityMap[s.Name] = int64(s.UsageCount)
		if s.UsageCount > 0 {
			usedCount++
		}
	}

	ratio := 0.0
	if len(settings) > 0 {
		ratio = float64(usedCount) / float64(len(settings))
	}

	return world_ops.NewGetWorldAuditOK().WithPayload(&models.WorldAuditData{
		IntensityMap: intensityMap,
		IcebergRatio: ratio,
	})
}

// HandleGetTemplates returns a list of professional world-building templates.
func HandleGetTemplates(params world_ops.GetWorldTemplatesParams) middleware.Responder {
	templates := []*models.WorldTemplate{
		{
			ID:             "t1",
			Category:       "Magic System",
			Name:           "Hard Magic: Resource Depletion",
			Description:    "Magic is fueled by a non-renewable physical resource (e.g., gems, crystals).",
			SuggestedLogic: "Cost: Physical resource consumed. Limit: Finite supply in the world.",
		},
		{
			ID:             "t2",
			Category:       "Geography",
			Name:           "Vertical City-State",
			Description:    "A society built entirely within a massive canyon or on a singular peak.",
			SuggestedLogic: "Constraint: Travel between levels is highly controlled. Resource: Water flows top-to-bottom.",
		},
		{
			ID:             "t3",
			Category:       "Race",
			Name:           "The Long-Lived Elders",
			Description:    "A race that lives for centuries but has extremely low birth rates.",
			SuggestedLogic: "Psychology: Highly conservative/risk-averse. Demographics: Children are extremely rare.",
		},
	}

	return world_ops.NewGetWorldTemplatesOK().WithPayload(templates)
}
