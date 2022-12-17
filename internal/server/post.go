package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/blog_engine/internal/dto"
	"github.com/HardDie/blog_engine/internal/logger"
	"github.com/HardDie/blog_engine/internal/service"
	"github.com/HardDie/blog_engine/internal/utils"
)

type Post struct {
	service service.IPost
}

func NewPost(service service.IPost) *Post {
	return &Post{
		service: service,
	}
}
func (s *Post) RegisterPublicRouter(router *mux.Router) {
	postRouter := router.PathPrefix("").Subrouter()
	postRouter.HandleFunc("/feed", s.Feed).Methods(http.MethodGet)
}
func (s *Post) RegisterPrivateRouter(router *mux.Router, middleware ...mux.MiddlewareFunc) {
	postRouter := router.PathPrefix("").Subrouter()
	postRouter.HandleFunc("", s.Create).Methods(http.MethodPost)
	postRouter.Use(middleware...)
}

func (s *Post) Feed(w http.ResponseWriter, r *http.Request) {
	req := &dto.FeedPostDTO{
		Limit: utils.GetInt32FromQuery(r, "limit", 0),
		Page:  utils.GetInt32FromQuery(r, "page", 0),
		Query: r.URL.Query().Get("query"),
	}

	err := req.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	posts, err := s.service.Feed(req)
	if err != nil {
		logger.Error.Println("Can't get feed:", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	err = utils.ResponseWithMeta(w, posts, nil)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}

func (s *Post) Create(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	req := &dto.CreatePostDTO{}
	err := utils.ParseJsonFromHTTPRequest(r.Body, req)
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Can't parse request", http.StatusBadRequest)
		return
	}

	err = req.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := s.service.Create(req, userID)
	if err != nil {
		logger.Error.Println("Can't create post:", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = utils.Response(w, post)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}
