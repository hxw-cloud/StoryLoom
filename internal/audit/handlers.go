package audit

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/hxw-cloud/StoryLoom/models"
	audit_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/audit"
)

// HandleAuditScene processes the GET /audit/scene/{sceneId} request.
// It invokes the LogicEngine to perform a validation pass and maps the results
// to the external Swagger response model.
func HandleAuditScene(params audit_ops.GetAuditSceneSceneIDParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()
	engine := NewLogicEngine()

	// Perform the validation logic
	auditRes, err := engine.ValidateScene(ctx, params.SceneID.String())
	if err != nil {
		// If scene not found or other error occurs, return 404 for now
		return audit_ops.NewGetAuditSceneSceneIDNotFound()
	}

	// Map internal result to API model
	// The generated model uses a boolean value directly, not a pointer.
	response := &models.AuditResult{
		IsValid: auditRes.IsValid,
		Issues:  auditRes.Issues,
	}

	return audit_ops.NewGetAuditSceneSceneIDOK().WithPayload(response)
}
