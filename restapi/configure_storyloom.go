// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/rs/cors"

	audit_internal "github.com/hxw-cloud/StoryLoom/internal/audit"
	"github.com/hxw-cloud/StoryLoom/internal/character"
	"github.com/hxw-cloud/StoryLoom/internal/conflict"
	"github.com/hxw-cloud/StoryLoom/internal/plot"
	"github.com/hxw-cloud/StoryLoom/internal/scene"
	"github.com/hxw-cloud/StoryLoom/internal/timeline"
	"github.com/hxw-cloud/StoryLoom/internal/world"
	"github.com/hxw-cloud/StoryLoom/pkg/db"
	"github.com/hxw-cloud/StoryLoom/restapi/operations"
	audit_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/audit"
	character_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/character"
	conflict_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/conflict"
	plot_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/plot"
	scene_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/scene"
	timeline_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/timeline"
	world_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/world"
)

//go:generate swagger generate server --target ..\..\StoryLoom --name Storyloom --spec ..\api\swagger.yaml --principal any

func configureFlags(api *operations.StoryloomAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
	_ = api
}

func configureAPI(api *operations.StoryloomAPI) http.Handler {
	// Initialize the global SQLite database connection before configuring handlers.
	db.InitDB("data.db")

	// Explicitly auto-migrate domain models.
	err := db.DB.AutoMigrate(
		&world.WorldSetting{},
		&world.WorldTemplate{},
		&world.HistoricalEvent{},
		&character.Character{},
		&character.Relationship{},
		&character.CharacterArc{},
		&plot.PlotCard{},
		&scene.Scene{},
		&timeline.TimelineEvent{},
		&conflict.Conflict{},
	)
	if err != nil {
		panic("Failed to migrate database schema: " + err.Error())
	}

	// configure the api here
	api.ServeError = errors.ServeError

	api.UseSwaggerUI()

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	// Register World-Building Handlers
	api.WorldGetWorldSettingsHandler = world_ops.GetWorldSettingsHandlerFunc(world.HandleGetSettings)
	api.WorldPostWorldSettingsHandler = world_ops.PostWorldSettingsHandlerFunc(world.HandlePostSettings)
	api.WorldGetWorldSettingsIDHandler = world_ops.GetWorldSettingsIDHandlerFunc(world.HandleGetSettingsID)
	api.WorldPutWorldSettingsIDHandler = world_ops.PutWorldSettingsIDHandlerFunc(world.HandlePutSettingsID)
	api.WorldGetWorldHistoryHandler = world_ops.GetWorldHistoryHandlerFunc(world.HandleGetHistory)
	api.WorldPostWorldHistoryHandler = world_ops.PostWorldHistoryHandlerFunc(world.HandlePostHistory)
	api.WorldGetWorldAuditHandler = world_ops.GetWorldAuditHandlerFunc(world.HandleGetAudit)
	api.WorldGetWorldTemplatesHandler = world_ops.GetWorldTemplatesHandlerFunc(world.HandleGetTemplates)

	// Register Character Handlers
	api.CharacterGetCharactersHandler = character_ops.GetCharactersHandlerFunc(character.HandleGetCharacters)
	api.CharacterPostCharactersHandler = character_ops.PostCharactersHandlerFunc(character.HandlePostCharacters)
	api.CharacterGetCharactersIDHandler = character_ops.GetCharactersIDHandlerFunc(character.HandleGetCharactersID)
	api.CharacterPutCharactersIDHandler = character_ops.PutCharactersIDHandlerFunc(character.HandlePutCharactersID)
	api.CharacterGetCharactersRelationshipsHandler = character_ops.GetCharactersRelationshipsHandlerFunc(character.HandleGetRelationships)
	api.CharacterPostCharactersRelationshipsHandler = character_ops.PostCharactersRelationshipsHandlerFunc(character.HandlePostRelationship)
	api.CharacterGetCharactersIDArcsHandler = character_ops.GetCharactersIDArcsHandlerFunc(character.HandleGetCharacterArcs)

	// Register Plot Handlers
	api.PlotGetPlotsHandler = plot_ops.GetPlotsHandlerFunc(plot.HandleGetPlots)
	api.PlotPostPlotsHandler = plot_ops.PostPlotsHandlerFunc(plot.HandlePostPlots)

	// Register Scene Handlers
	api.SceneGetScenesHandler = scene_ops.GetScenesHandlerFunc(scene.HandleGetScenes)
	api.ScenePostScenesHandler = scene_ops.PostScenesHandlerFunc(scene.HandlePostScenes)

	// Register Timeline Handlers
	api.TimelineGetTimelineEventsHandler = timeline_ops.GetTimelineEventsHandlerFunc(timeline.HandleGetTimelineEvents)
	api.TimelinePostTimelineEventsHandler = timeline_ops.PostTimelineEventsHandlerFunc(timeline.HandlePostTimelineEvents)

	// Register Audit Handlers
	api.AuditGetAuditSceneSceneIDHandler = audit_ops.GetAuditSceneSceneIDHandlerFunc(audit_internal.HandleAuditScene)

	// Register Conflict Handlers
	api.ConflictGetConflictsHandler = conflict_ops.GetConflictsHandlerFunc(conflict.HandleGetConflicts)
	api.ConflictPostConflictsHandler = conflict_ops.PostConflictsHandlerFunc(conflict.HandlePostConflicts)

	api.PreServerShutdown = func() {}
	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	_ = tlsConfig
}

func configureServer(server *http.Server, scheme, addr string) {
	_ = server
	_ = scheme
	_ = addr
}

func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		Debug:          false,
	}).Handler(handler)
}
