package audit

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/hxw-cloud/StoryLoom/models"
	audit_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/audit"
)

// HandleAuditScene processes the GET /audit/scene/{sceneId} request.
// It invokes the LogicEngine to perform a validation pass and maps the results
// to the external Swagger response model.
func HandleAuditScene(params audit_ops.GetAuditSceneParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()
	engine := NewLogicEngine()

	// Perform the validation logic
	auditRes, err := engine.ValidateScene(ctx, params.SceneID.String())
	if err != nil {
		// If scene not found or other error occurs, return 404 for now
		// In a refined API, we would distinguish between 404 and 500.
		return audit_ops.NewGetAuditSceneNotFound()
	}

	// Map internal result to API model
	response := &models.AuditResult{
		IsValid: &auditRes.IsValid,
		Issues:  auditRes.Issues,
	}

	return audit_ops.NewGetAuditSceneOK().WithPayload(response)
}
