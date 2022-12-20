package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/blog_engine/internal/dto"
	"github.com/HardDie/blog_engine/internal/logger"
	"github.com/HardDie/blog_engine/internal/service"
	"github.com/HardDie/blog_engine/internal/utils"
)

type Auth struct {
	service service.IAuth
}

func NewAuth(service service.IAuth) *Auth {
	return &Auth{
		service: service,
	}
}
func (s *Auth) RegisterPublicRouter(router *mux.Router) {
	authRouter := router.PathPrefix("").Subrouter()
	authRouter.HandleFunc("/register", s.Register).Methods(http.MethodPost)
	authRouter.HandleFunc("/login", s.Login).Methods(http.MethodPost)
}

// swagger:parameters AuthRegisterRequest
type AuthRegisterRequest struct {
	// In: body
	Body struct {
		dto.RegisterDTO
	}
}

// swagger:response AuthRegisterResponse
type AuthRegisterResponse struct {
}

// swagger:route POST /api/v1/auth/register Auth AuthRegisterRequest
//
// Registration form
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
//	  200: AuthRegisterResponse
func (s *Auth) Register(w http.ResponseWriter, r *http.Request) {
	req := &dto.RegisterDTO{}
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

	user, err := s.service.Register(req)
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	session, err := s.service.GenerateCookie(*user.ID)
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	utils.SetSessionCookie(session, w)
}

// swagger:parameters AuthLoginRequest
type AuthLoginRequest struct {
	// In: body
	Body struct {
		dto.LoginDTO
	}
}

// swagger:response AuthLoginResponse
type AuthLoginResponse struct {
}

// swagger:route POST /api/v1/auth/login Auth AuthLoginRequest
//
// Login form
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
//	  200: AuthLoginResponse
func (s *Auth) Login(w http.ResponseWriter, r *http.Request) {
	req := &dto.LoginDTO{}
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

	user, err := s.service.Login(req)
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	session, err := s.service.GenerateCookie(*user.ID)
	if err != nil {
		logger.Error.Printf(err.Error())
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	utils.SetSessionCookie(session, w)
}
