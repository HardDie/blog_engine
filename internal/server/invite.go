package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/blog_engine/internal/logger"
	"github.com/HardDie/blog_engine/internal/service"
)

type Invite struct {
	service service.IInvite
}

func NewInvite(service service.IInvite) *Invite {
	return &Invite{
		service: service,
	}
}
func (s *Invite) RegisterRouter(router *mux.Router) {
	router.HandleFunc("/generate", s.Generate).Methods(http.MethodGet)
	router.HandleFunc("/revoke", s.Revoke).Methods(http.MethodDelete)
}

func (s *Invite) Generate(w http.ResponseWriter, r *http.Request) {
	inviteCode, err := s.service.Generate(1)
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
func (s *Invite) Revoke(w http.ResponseWriter, r *http.Request) {
	err := s.service.Revoke(1)
	if err != nil {
		logger.Error.Println("Error revoking invite code:", err.Error())
		http.Error(w, "Can't revoke invite code", http.StatusInternalServerError)
		return
	}
}
