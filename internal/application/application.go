package application

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/HardDie/blog_engine/internal/config"
	"github.com/HardDie/blog_engine/internal/db"
	"github.com/HardDie/blog_engine/internal/middleware"
	"github.com/HardDie/blog_engine/internal/migration"
	repositoryInvite "github.com/HardDie/blog_engine/internal/repository/invite"
	repositoryPassword "github.com/HardDie/blog_engine/internal/repository/password"
	repositoryPost "github.com/HardDie/blog_engine/internal/repository/post"
	repositorySession "github.com/HardDie/blog_engine/internal/repository/session"
	repositoryUser "github.com/HardDie/blog_engine/internal/repository/user"
	"github.com/HardDie/blog_engine/internal/server"
	"github.com/HardDie/blog_engine/internal/service"
	serviceAuth "github.com/HardDie/blog_engine/internal/service/auth"
	serviceInvite "github.com/HardDie/blog_engine/internal/service/invite"
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
	userRepository := repositoryUser.New(app.DB)
	passwordRepository := repositoryPassword.New(app.DB)
	sessionRepository := repositorySession.New(app.DB)
	inviteRepository := repositoryInvite.New(app.DB)
	postRepository := repositoryPost.New(app.DB)

	// Init services
	authService := serviceAuth.New(app.Cfg, userRepository, passwordRepository, sessionRepository, inviteRepository)
	inviteService := serviceInvite.New(userRepository, inviteRepository)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)
	timeoutMiddleware := middleware.NewTimeoutRequestMiddleware(time.Duration(app.Cfg.RequestTimeout) * time.Second)

	// Register servers
	authRouter := v1Router.PathPrefix("/auth").Subrouter()
	authServer := server.NewAuth(app.Cfg, authService)
	authServer.RegisterPublicRouter(authRouter)
	authServer.RegisterPrivateRouter(authRouter, timeoutMiddleware.RequestMiddleware, authMiddleware.RequestMiddleware)

	inviteRouter := v1Router.PathPrefix("/invites").Subrouter()
	server.NewInvite(
		inviteService,
	).RegisterPrivateRouter(inviteRouter, timeoutMiddleware.RequestMiddleware, authMiddleware.RequestMiddleware)

	postsRouter := v1Router.PathPrefix("/posts").Subrouter()
	postServer := server.NewPost(
		service.NewPost(
			postRepository,
		),
	)
	postServer.RegisterPublicRouter(postsRouter, timeoutMiddleware.RequestMiddleware)
	postServer.RegisterPrivateRouter(postsRouter, timeoutMiddleware.RequestMiddleware, authMiddleware.RequestMiddleware)

	userRouter := v1Router.PathPrefix("/user").Subrouter()
	userServer := server.NewUser(
		service.NewUser(userRepository, passwordRepository),
	)
	userServer.RegisterPublicRouter(userRouter, timeoutMiddleware.RequestMiddleware)
	userServer.RegisterPrivateRouter(userRouter, timeoutMiddleware.RequestMiddleware, authMiddleware.RequestMiddleware)

	return app, nil
}

func (app *Application) Run() error {
	return http.ListenAndServe(app.Cfg.Port, app.Router)
}
