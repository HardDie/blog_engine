package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/blog_engine/internal/dto"
	"github.com/HardDie/blog_engine/internal/entity"
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
	postRouter.HandleFunc("/{id:[0-9]+}", s.Edit).Methods(http.MethodPut)
	postRouter.HandleFunc("", s.List).Methods(http.MethodGet)
	postRouter.Use(middleware...)
}

/*
 * Public
 */

// swagger:parameters PostFeedRequest
type PostFeedRequest struct {
	// In: query
	dto.FeedPostDTO
}

// swagger:response PostFeedResponse
type PostFeedResponse struct {
	// In: body
	Body struct {
		Data []*entity.Post `json:"data"`
	}
}

// swagger:route GET /api/v1/posts/feed Post PostFeedRequest
//
// Get feed
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: https
//
//	Responses:
//	  200: PostFeedResponse
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

/*
 * Private
 */

// swagger:parameters PostCreateRequest
type PostCreateRequest struct {
	// In: body
	Body struct {
		dto.CreatePostDTO
	}
}

// swagger:response PostCreateResponse
type PostCreateResponse struct {
	// In: body
	Body struct {
		Data *entity.Post `json:"data"`
	}
}

// swagger:route POST /api/v1/posts Post PostCreateRequest
//
// Create post form
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: https
//
//	Responses:
//	  200: PostCreateResponse
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

// swagger:parameters PostEditRequest
type PostEditRequest struct {
	// In: path
	ID int32 `json:"id"`
	// In: body
	Body struct {
		dto.EditPostDTO
	}
}

// swagger:response PostEditResponse
type PostEditResponse struct {
	// In: body
	Body struct {
		Data *entity.Post `json:"data"`
	}
}

// swagger:route PUT /api/v1/posts/{id} Post PostEditRequest
//
// Edit post form
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: https
//
//	Responses:
//	  200: PostEditResponse
func (s *Post) Edit(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	req := &dto.EditPostDTO{}
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

	req.ID, err = utils.GetInt32FromPath(r, "id")
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Bad id in path", http.StatusBadRequest)
		return
	}

	post, err := s.service.Edit(req, userID)
	if err != nil {
		logger.Error.Println("Can't edit post:", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = utils.Response(w, post)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}

// swagger:parameters PostListRequest
type PostListRequest struct {
	// In: query
	dto.ListPostDTO
}

// swagger:response PostListResponse
type PostListResponse struct {
	// In: body
	Body struct {
		Data []*entity.Post `json:"data"`
	}
}

// swagger:route GET /api/v1/posts Post PostListRequest
//
// Get list of posts for current user
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: https
//
//	Responses:
//	  200: PostListResponse
func (s *Post) List(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	req := &dto.ListPostDTO{
		Limit: utils.GetInt32FromQuery(r, "limit", 0),
		Page:  utils.GetInt32FromQuery(r, "page", 0),
		Query: r.URL.Query().Get("query"),
	}

	err := req.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	posts, err := s.service.List(req, userID)
	if err != nil {
		logger.Error.Println("Can't get list of posts:", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	err = utils.ResponseWithMeta(w, posts, nil)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}
