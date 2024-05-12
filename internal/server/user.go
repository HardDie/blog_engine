package server

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/blog_engine/internal/dto"
	"github.com/HardDie/blog_engine/internal/entity"
	"github.com/HardDie/blog_engine/internal/logger"
	serviceUser "github.com/HardDie/blog_engine/internal/service/user"
	"github.com/HardDie/blog_engine/internal/utils"
)

type User struct {
	user serviceUser.IUser
}

func NewUser(user serviceUser.IUser) *User {
	return &User{
		user: user,
	}
}
func (s *User) RegisterPublicRouter(router *mux.Router, middleware ...mux.MiddlewareFunc) {
	userRouter := router.PathPrefix("").Subrouter()
	userRouter.HandleFunc("/{id:[0-9]+}", s.Get).Methods(http.MethodGet)
	userRouter.Use(middleware...)
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
	ctx := r.Context()

	userID, err := utils.GetInt64FromPath(r, "id")
	if err != nil {
		logger.Error.Printf("User.Get() GetInt32FromPath: %s", err.Error())
		utils.WriteJSONHTTPResponse(w, http.StatusBadRequest, JSONResponse{
			Error: "Bad id in path",
		})
		return
	}
	req := dto.GetUserDTO{
		ID: userID,
	}

	err = GetValidator().Struct(req)
	if err != nil {
		utils.WriteJSONHTTPResponse(w, http.StatusBadRequest, JSONResponse{
			Error: "Validation",
			Data:  err.Error(),
		})
		return
	}

	user, err := s.user.Get(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, serviceUser.ErrorUserNotFound):
			utils.WriteJSONHTTPResponse(w, http.StatusBadRequest, JSONResponse{
				Error: "User not found",
			})
			return
		}
		logger.Error.Printf("User.Get() Get: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJSONHTTPResponse(w, http.StatusOK, JSONResponse{
		Data: user,
	})
	if err != nil {
		logger.Error.Printf("User.Get() WriteJSONHTTPResponse: %s", err.Error())
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
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	req := &dto.UpdatePasswordDTO{}
	err := utils.ParseJsonFromHTTPRequest(r.Body, req)
	if err != nil {
		logger.Error.Printf("User.Password() ParseJsonFromHTTPRequest: %s", err.Error())
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

	err = s.user.Password(ctx, req, userID)
	if err != nil {
		switch {
		case errors.Is(err, serviceUser.ErrorInvalidPassword):
			utils.WriteJSONHTTPResponse(w, http.StatusBadRequest, JSONResponse{
				Error: "Invalid old password",
			})
			return
		}
		logger.Error.Printf("User.Password() Password: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	req := &dto.UpdateProfileDTO{}
	err := utils.ParseJsonFromHTTPRequest(r.Body, req)
	if err != nil {
		logger.Error.Printf("User.Profile() ParseJsonFromHTTPRequest: %s", err.Error())
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

	user, err := s.user.Profile(ctx, req, userID)
	if err != nil {
		logger.Error.Printf("User.Profile() Profile: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJSONHTTPResponse(w, http.StatusOK, JSONResponse{
		Data: user,
	})
	if err != nil {
		logger.Error.Printf("User.Profile() WriteJSONHTTPResponse: %s", err.Error())
	}
}
