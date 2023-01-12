package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/blog_engine/internal/logger"
	"github.com/HardDie/blog_engine/internal/service"
	"github.com/HardDie/blog_engine/internal/utils"
)

type Invite struct {
	service service.IInvite
}

func NewInvite(service service.IInvite) *Invite {
	return &Invite{
		service: service,
	}
}
func (s *Invite) RegisterPrivateRouter(router *mux.Router, middleware ...mux.MiddlewareFunc) {
	inviteRouter := router.PathPrefix("").Subrouter()
	inviteRouter.HandleFunc("/generate", s.Generate).Methods(http.MethodGet)
	inviteRouter.HandleFunc("/revoke", s.Revoke).Methods(http.MethodDelete)
	inviteRouter.Use(middleware...)
}

/*
 * Private
 */

// swagger:parameters InviteGenerateRequest
type InviteGenerateRequest struct {
}

// swagger:response InviteGenerateResponse
type InviteGenerateResponse struct {
}

// swagger:route GET /api/v1/invites/generate Invite InviteGenerateRequest
//
// # Generate a new invitation code
//
//	Responses:
//	  200: InviteGenerateResponse
func (s *Invite) Generate(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	inviteCode, err := s.service.Generate(userID)
	if err != nil {
		logger.Error.Println("Error generating invite code:", err.Error())
		http.Error(w, "Can't generate invite code", http.StatusInternalServerError)
		return
	}
	_, err = fmt.Fprintf(w, inviteCode)
	if err != nil {
		logger.Error.Println("Error sending response:", err.Error())
	}
}

// swagger:parameters InviteRevokeRequest
type InviteRevokeRequest struct {
}

// swagger:response InviteRevokeResponse
type InviteRevokeResponse struct {
}

// swagger:route DELETE /api/v1/invites/revoke Invite InviteRevokeRequest
//
// # Revoke the generated invitation code
//
//	Responses:
//	  200: InviteRevokeResponse
func (s *Invite) Revoke(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	err := s.service.Revoke(userID)
	if err != nil {
		logger.Error.Println("Error revoking invite code:", err.Error())
		http.Error(w, "Can't revoke invite code", http.StatusInternalServerError)
		return
	}
}
