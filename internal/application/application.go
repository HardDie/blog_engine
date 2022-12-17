package application

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/blog_engine/internal/config"
	"github.com/HardDie/blog_engine/internal/db"
	"github.com/HardDie/blog_engine/internal/middleware"
	"github.com/HardDie/blog_engine/internal/migration"
	"github.com/HardDie/blog_engine/internal/repository"
	"github.com/HardDie/blog_engine/internal/server"
	"github.com/HardDie/blog_engine/internal/service"
)

type Application struct {
	Cfg    *config.Config
	DB     *db.DB
	Router *mux.Router
}

func Get() (*Application, error) {
	app := &Application{
		Cfg:    config.Get(),
		Router: mux.NewRouter(),
	}

	// Init DB
	newDB, err := db.Get(app.Cfg.DBPath)
	if err != nil {
		return nil, err
	}
	app.DB = newDB

	// Init migrations
	err = migration.NewMigrate(app.DB).Up()
	if err != nil {
		return nil, err
	}

	// Prepare router
	apiRouter := app.Router.PathPrefix("/api").Subrouter()
	v1Router := apiRouter.PathPrefix("/v1").Subrouter()

	// Init repositories
	userRepository := repository.NewUser(app.DB)
	passwordRepository := repository.NewPassword(app.DB)
	sessionRepository := repository.NewSession(app.DB)
	inviteRepository := repository.NewInvite(app.DB)
	postRepository := repository.NewPost(app.DB)

	// Init services
	authService := service.NewAuth(userRepository, passwordRepository, sessionRepository, inviteRepository)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Register servers
	authRouter := v1Router.PathPrefix("/auth").Subrouter()
	server.NewAuth(authService).
		RegisterPublicRouter(authRouter)

	inviteRouter := v1Router.PathPrefix("/invites").Subrouter()
	server.NewInvite(
		service.NewInvite(
			userRepository,
			inviteRepository,
		),
	).RegisterPrivateRouter(inviteRouter, authMiddleware.RequestMiddleware)

	postsRouter := v1Router.PathPrefix("/posts").Subrouter()
	postServer := server.NewPost(
		service.NewPost(
			postRepository,
		),
	)
	postServer.RegisterPublicRouter(postsRouter)
	postServer.RegisterPrivateRouter(postsRouter, authMiddleware.RequestMiddleware)

	return app, nil
}

func (app *Application) Run() error {
	return http.ListenAndServe(app.Cfg.Port, app.Router)
}
