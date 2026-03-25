// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	"github.com/hxw-cloud/StoryLoom/internal/character"
	"github.com/hxw-cloud/StoryLoom/internal/plot"
	"github.com/hxw-cloud/StoryLoom/internal/scene"
	"github.com/hxw-cloud/StoryLoom/internal/world"
	"github.com/hxw-cloud/StoryLoom/pkg/db"
	"github.com/hxw-cloud/StoryLoom/restapi/operations"
	character_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/character"
	plot_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/plot"
	scene_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/scene"
	world_ops "github.com/hxw-cloud/StoryLoom/restapi/operations/world"
)

//go:generate swagger generate server --target ..\..\StoryLoom --name Storyloom --spec ..\api\swagger.yaml --principal any

func configureFlags(api *operations.StoryloomAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
	_ = api
}

func configureAPI(api *operations.StoryloomAPI) http.Handler {
	// Initialize the global SQLite database connection before configuring handlers.
	// This ensures that all endpoint handlers have access to the active database session.
	// "data.db" is a local SQLite file; for production, this should ideally be driven by config.
	db.InitDB("data.db")

	// Explicitly auto-migrate domain models to break the import cycle between db and domain packages.
	// This synchronizes the database schema with the GORM models on application startup.
	err := db.DB.AutoMigrate(
		&world.WorldSetting{},
		&character.Character{},
		&plot.PlotCard{},
		&scene.Scene{},
	)
	if err != nil {
		panic("Failed to migrate database schema: " + err.Error())
	}

	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...any)
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Register World-Building Handlers
	api.WorldGetWorldSettingsHandler = world_ops.GetWorldSettingsHandlerFunc(world.HandleGetSettings)
	api.WorldPostWorldSettingsHandler = world_ops.PostWorldSettingsHandlerFunc(world.HandlePostSettings)

	// Register Character Handlers
	api.CharacterGetCharactersHandler = character_ops.GetCharactersHandlerFunc(character.HandleGetCharacters)
	api.CharacterPostCharactersHandler = character_ops.PostCharactersHandlerFunc(character.HandlePostCharacters)

	// Register Plot Handlers
	api.PlotGetPlotsHandler = plot_ops.GetPlotsHandlerFunc(plot.HandleGetPlots)
	api.PlotPostPlotsHandler = plot_ops.PostPlotsHandlerFunc(plot.HandlePostPlots)

	// Register Scene Handlers
	api.SceneGetScenesHandler = scene_ops.GetScenesHandlerFunc(scene.HandleGetScenes)
	api.ScenePostScenesHandler = scene_ops.PostScenesHandlerFunc(scene.HandlePostScenes)

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
	_ = tlsConfig
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(server *http.Server, scheme, addr string) {
	_ = server
	_ = scheme
	_ = addr
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
