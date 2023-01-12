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

type User struct {
	service service.IUser
}

func NewUser(service service.IUser) *User {
	return &User{
		service: service,
	}
}
func (s *User) RegisterPublicRouter(router *mux.Router) {
	userRouter := router.PathPrefix("").Subrouter()
	userRouter.HandleFunc("/{id:[0-9]+}", s.Get).Methods(http.MethodGet)
}
func (s *User) RegisterPrivateRouter(router *mux.Router, middleware ...mux.MiddlewareFunc) {
	userRouter := router.PathPrefix("").Subrouter()
	userRouter.HandleFunc("/password", s.Password).Methods(http.MethodPut)
	userRouter.HandleFunc("/profile", s.Profile).Methods(http.MethodPut)
	userRouter.Use(middleware...)
}

/*
 * Public
 */

// swagger:parameters UserGetRequest
type UserGetRequest struct {
	// In: path
	ID int32 `json:"id"`
}

// swagger:response UserGetResponse
type UserGetResponse struct {
	// In: body
	Body struct {
		Data *entity.User `json:"data"`
	}
}

// swagger:route GET /api/v1/user/{id} User UserGetRequest
//
// # Getting information about a user by ID
//
//	Responses:
//	  200: UserGetResponse
func (s *User) Get(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetInt32FromPath(r, "id")
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Bad id in path", http.StatusBadRequest)
		return
	}

	user, err := s.service.Get(userID)
	if err != nil {
		logger.Error.Println("Can't get post:", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = utils.Response(w, user)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}

/*
 * Private
 */

// swagger:parameters UserPasswordRequest
type UserPasswordRequest struct {
	// In: body
	Body struct {
		dto.UpdatePasswordDTO
	}
}

// swagger:response UserPasswordResponse
type UserPasswordResponse struct {
}

// swagger:route PUT /api/v1/user/password User UserPasswordRequest
//
// # Updating the password for a user
//
//	Responses:
//	  200: UserPasswordResponse
func (s *User) Password(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	req := &dto.UpdatePasswordDTO{}
	err := utils.ParseJsonFromHTTPRequest(r.Body, req)
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Can't parse request", http.StatusBadRequest)
		return
	}

	err = GetValidator().Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.service.Password(req, userID)
	if err != nil {
		logger.Error.Println("Can't update password:", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

// swagger:parameters UserProfileRequest
type UserProfileRequest struct {
	// In: body
	Body struct {
		dto.UpdateProfileDTO
	}
}

// swagger:response UserProfileResponse
type UserProfileResponse struct {
}

// swagger:route PUT /api/v1/user/profile User UserProfileRequest
//
// # Updating user information
//
//	Responses:
//	  200: UserProfileResponse
func (s *User) Profile(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	req := &dto.UpdateProfileDTO{}
	err := utils.ParseJsonFromHTTPRequest(r.Body, req)
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Can't parse request", http.StatusBadRequest)
		return
	}

	err = GetValidator().Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.service.Profile(req, userID)
	if err != nil {
		logger.Error.Println("Can't update user profile:", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = utils.Response(w, user)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}
