package server

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/blog_engine/internal/dto"
	"github.com/HardDie/blog_engine/internal/entity"
	"github.com/HardDie/blog_engine/internal/logger"
	servicePost "github.com/HardDie/blog_engine/internal/service/post"
	"github.com/HardDie/blog_engine/internal/utils"
)

type Post struct {
	postService servicePost.IPost
}

func NewPost(post servicePost.IPost) *Post {
	return &Post{
		postService: post,
	}
}
func (s *Post) RegisterPublicRouter(router *mux.Router, middleware ...mux.MiddlewareFunc) {
	postRouter := router.PathPrefix("").Subrouter()
	postRouter.HandleFunc("/feed", s.Feed).Methods(http.MethodGet)
	postRouter.HandleFunc("/{id:[0-9]+}", s.PublicGet).Methods(http.MethodGet)
	postRouter.Use(middleware...)
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
// # Get feed
//
//	Responses:
//	  200: PostFeedResponse
func (s *Post) Feed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &dto.FeedPostDTO{
		Limit: utils.GetInt32FromQuery(r, "limit", 0),
		Page:  utils.GetInt32FromQuery(r, "page", 0),
		Query: r.URL.Query().Get("query"),
	}

	err := GetValidator().Struct(req)
	if err != nil {
		utils.WriteJSONHTTPResponse(w, http.StatusBadRequest, JSONResponse{
			Error: "Validation",
			Data:  err.Error(),
		})
		return
	}

	posts, total, err := s.postService.Feed(ctx, req)
	if err != nil {
		logger.Error.Printf("Feed() Feed: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	meta := &utils.Meta{
		Total: total,
		Limit: req.Limit,
		Page:  req.Page,
	}
	err = utils.ResponseWithMeta(w, posts, meta)
	if err != nil {
		logger.Error.Printf("Feed() ResponseWithMeta: %s", err.Error())
	}
}

// swagger:parameters PostPublicGetRequest
type PostPublicGetRequest struct {
	// In: path
	ID int32 `json:"id"`
}

// swagger:response PostPublicGetResponse
type PostPublicGetResponse struct {
	// In: body
	Body struct {
		Data *entity.Post `json:"data"`
	}
}

// swagger:route GET /api/v1/posts/{id} Post PostPublicGetRequest
//
// # Get public post
//
//	Responses:
//	  200: PostPublicGetResponse
func (s *Post) PublicGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID, err := utils.GetInt32FromPath(r, "id")
	if err != nil {
		logger.Error.Printf("PublicGet() GetInt32FromPath: %s", err.Error())
		utils.WriteJSONHTTPResponse(w, http.StatusBadRequest, JSONResponse{
			Error: "Bad id in path",
		})
		return
	}
	req := &dto.PublicGetDTO{
		ID: postID,
	}

	err = GetValidator().Struct(req)
	if err != nil {
		utils.WriteJSONHTTPResponse(w, http.StatusBadRequest, JSONResponse{
			Error: "Validation",
			Data:  err.Error(),
		})
		return
	}

	post, err := s.postService.PublicGet(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, servicePost.ErrorPostNotFound):
			err = utils.WriteJSONHTTPResponse(w, http.StatusNotFound, JSONResponse{
				Error: "Post not found",
			})
			return
		}
		logger.Error.Printf("PublicGet() PublicGet: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJSONHTTPResponse(w, http.StatusOK, JSONResponse{
		Data: post,
	})
	if err != nil {
		logger.Error.Printf("PublicGet() WriteJSONHTTPResponse: %s", err.Error())
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
// # Post creation form
//
//	Responses:
//	  201: PostCreateResponse
func (s *Post) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	req := &dto.CreatePostDTO{}
	err := utils.ParseJsonFromHTTPRequest(r.Body, req)
	if err != nil {
		logger.Error.Printf("Create() ParseJsonFromHTTPRequest: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = GetValidator().Struct(req)
	if err != nil {
		utils.WriteJSONHTTPResponse(w, http.StatusBadRequest, JSONResponse{
			Error: "Validation",
			Data:  err.Error(),
		})
		return
	}

	post, err := s.postService.Create(ctx, req, userID)
	if err != nil {
		logger.Error.Printf("Create() Create: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJSONHTTPResponse(w, http.StatusCreated, JSONResponse{
		Data: post,
	})
	if err != nil {
		logger.Error.Printf("Create() WriteJSONHTTPResponse: %s", err.Error())
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
// # Edit post form
//
//	Responses:
//	  200: PostEditResponse
func (s *Post) Edit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	req := &dto.EditPostDTO{}
	err := utils.ParseJsonFromHTTPRequest(r.Body, req)
	if err != nil {
		logger.Error.Printf("Edit() ParseJsonFromHTTPRequest: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	req.ID, err = utils.GetInt32FromPath(r, "id")
	if err != nil {
		logger.Error.Printf("Edit() GetInt32FromPath: %s", err.Error())
		utils.WriteJSONHTTPResponse(w, http.StatusBadRequest, JSONResponse{
			Error: "Bad id in path",
		})
		return
	}

	err = GetValidator().Struct(req)
	if err != nil {
		utils.WriteJSONHTTPResponse(w, http.StatusBadRequest, JSONResponse{
			Error: "Validation",
			Data:  err.Error(),
		})
		return
	}

	post, err := s.postService.Edit(ctx, req, userID)
	if err != nil {
		logger.Error.Printf("Edit() Edit: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJSONHTTPResponse(w, http.StatusOK, JSONResponse{
		Data: post,
	})
	if err != nil {
		logger.Error.Printf("Edit() WriteJSONHTTPResponse: %s", err.Error())
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
// # Get a list of posts for the current user
//
//	Responses:
//	  200: PostListResponse
func (s *Post) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	req := &dto.ListPostDTO{
		Limit: utils.GetInt32FromQuery(r, "limit", 0),
		Page:  utils.GetInt32FromQuery(r, "page", 0),
		Query: r.URL.Query().Get("query"),
	}

	err := GetValidator().Struct(req)
	if err != nil {
		utils.WriteJSONHTTPResponse(w, http.StatusBadRequest, JSONResponse{
			Error: "Validation",
			Data:  err.Error(),
		})
		return
	}

	posts, total, err := s.postService.List(ctx, req, userID)
	if err != nil {
		logger.Error.Printf("List() List: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	meta := &utils.Meta{
		Total: total,
		Limit: req.Limit,
		Page:  req.Page,
	}
	err = utils.ResponseWithMeta(w, posts, meta)
	if err != nil {
		logger.Error.Printf("List() ResponseWithMeta: %s", err.Error())
	}
}
