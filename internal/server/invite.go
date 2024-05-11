package server

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/blog_engine/internal/logger"
	serviceInvite "github.com/HardDie/blog_engine/internal/service/invite"
	"github.com/HardDie/blog_engine/internal/utils"
)

type Invite struct {
	inviteService serviceInvite.IInvite
}

func NewInvite(invite serviceInvite.IInvite) *Invite {
	return &Invite{
		inviteService: invite,
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
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	inviteCode, err := s.inviteService.Generate(ctx, userID)
	if err != nil {
		logger.Error.Printf("Invite() Generate: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJSONHTTPResponse(w, http.StatusOK, JSONResponse{
		Data: inviteCode,
	})
	if err != nil {
		logger.Error.Printf("Invite() WriteJSONHTTPResponse: %s", err.Error())
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
	ctx := r.Context()
	userID := utils.GetUserIDFromContext(ctx)

	err := s.inviteService.Revoke(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, serviceInvite.ErrorInviteNotFound):
			utils.WriteJSONHTTPResponse(w, http.StatusBadRequest, JSONResponse{
				Error: "Invite not found",
			})
			return
		}
		logger.Error.Printf("Revoke() Revoke: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
