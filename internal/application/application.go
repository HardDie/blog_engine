package application

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/HardDie/blog_engine/internal/config"
	"github.com/HardDie/blog_engine/internal/db"
	"github.com/HardDie/blog_engine/internal/middleware"
	"github.com/HardDie/blog_engine/internal/migration"
	"github.com/HardDie/blog_engine/internal/repository/invite"
	"github.com/HardDie/blog_engine/internal/repository/password"
	"github.com/HardDie/blog_engine/internal/repository/post"
	"github.com/HardDie/blog_engine/internal/repository/session"
	"github.com/HardDie/blog_engine/internal/repository/user"
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
	userRepository := user.New(app.DB)
	passwordRepository := password.New(app.DB)
	sessionRepository := session.New(app.DB)
	inviteRepository := invite.New(app.DB)
	postRepository := post.New(app.DB)

	// Init services
	authService := service.NewAuth(app.Cfg, userRepository, passwordRepository, sessionRepository, inviteRepository)

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
		service.NewInvite(
			userRepository,
			inviteRepository,
		),
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
