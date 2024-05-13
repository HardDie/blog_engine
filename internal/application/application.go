package application

import (
	"net/http"
	"time"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
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
	serviceAuth "github.com/HardDie/blog_engine/internal/service/auth"
	serviceInvite "github.com/HardDie/blog_engine/internal/service/invite"
	servicePost "github.com/HardDie/blog_engine/internal/service/post"
	serviceUser "github.com/HardDie/blog_engine/internal/service/user"
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
	app.Router.Use(
		middleware.CorsMiddleware,
		chiMiddleware.RequestID,
		chiMiddleware.RealIP,
		chiMiddleware.Logger,
		chiMiddleware.Recoverer,
	)
	app.Router.MethodNotAllowedHandler = http.HandlerFunc(notAllowed)

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
	passwordRepository := repositoryPassword.New(app.DB.DB)
	sessionRepository := repositorySession.New(app.DB.DB)
	inviteRepository := repositoryInvite.New(app.DB.DB)
	postRepository := repositoryPost.New(app.DB.DB)

	// Init services
	authService := serviceAuth.New(app.Cfg, userRepository, passwordRepository, sessionRepository, inviteRepository)
	inviteService := serviceInvite.New(userRepository, inviteRepository)
	postService := servicePost.New(postRepository, userRepository)
	userService := serviceUser.New(userRepository, passwordRepository)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)
	timeoutMiddleware := chiMiddleware.Timeout(time.Duration(app.Cfg.RequestTimeout) * time.Second)

	// Register servers
	authRouter := v1Router.PathPrefix("/auth").Subrouter()
	authServer := server.NewAuth(app.Cfg, authService)
	authServer.RegisterPublicRouter(authRouter)
	authServer.RegisterPrivateRouter(authRouter, timeoutMiddleware, authMiddleware.RequestMiddleware)

	inviteRouter := v1Router.PathPrefix("/invites").Subrouter()
	inviteServer := server.NewInvite(inviteService)
	inviteServer.RegisterPrivateRouter(inviteRouter, timeoutMiddleware, authMiddleware.RequestMiddleware)

	postsRouter := v1Router.PathPrefix("/posts").Subrouter()
	postServer := server.NewPost(postService)
	postServer.RegisterPublicRouter(postsRouter, timeoutMiddleware)
	postServer.RegisterPrivateRouter(postsRouter, timeoutMiddleware, authMiddleware.RequestMiddleware)

	userRouter := v1Router.PathPrefix("/user").Subrouter()
	userServer := server.NewUser(userService)
	userServer.RegisterPublicRouter(userRouter, timeoutMiddleware)
	userServer.RegisterPrivateRouter(userRouter, timeoutMiddleware, authMiddleware.RequestMiddleware)

	return app, nil
}

func (app *Application) Run() error {
	return http.ListenAndServe(app.Cfg.Port, app.Router)
}

func notAllowed(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == http.MethodOptions {
		middleware.SetupCors(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}
