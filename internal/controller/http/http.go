// Package http implements controller by HTTP handler.
package http

import (
	"messenger/internal/usecase"
	"net/http"

	"github.com/go-chi/chi"
)

// HTTPControllerInterfece defines all handler methods.
type HTTPControllerInterface interface {
	Signup(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	GetChats(w http.ResponseWriter, r *http.Request)
	GetMessages(w http.ResponseWriter, r *http.Request)
	CreateChat(w http.ResponseWriter, r *http.Request)
	SendMessage(w http.ResponseWriter, r *http.Request)
	DeleteChat(w http.ResponseWriter, r *http.Request)
	DeleteAccount(w http.ResponseWriter, r *http.Request)
}

// HTTPController is a default implementation of HTTPControllerInterface.
type HTTPController struct {
	r  *chi.Mux
	uc *usecase.Usecase
}

// NewHTTPController initializes HTTPController and returns pointer on it.
func NewHTTPController(uc *usecase.Usecase) *HTTPController {
	c := HTTPController{uc: uc}

	c.r = chi.NewRouter()

	return &c
}

// DefineEndpoint is a method for definding all endpoints in chi router.
func (c *HTTPController) DefindEndpoints() {
	c.r.Get("/users/{id}/login", c.Login)
	c.r.Get("/chats", c.GetChats)
	c.r.Get("/chats/{id}/messages", c.GetMessages)
	c.r.Post("/users", c.Signup)
	c.r.Post("/chats", c.CreateChat)
	c.r.Post("/chats/{id}/messages", c.SendMessage)
	c.r.Delete("/users/{id}", c.DeleteAccount)
	c.r.Delete("/chats/{id}", c.DeleteChat)
}

// UseMiddleware adds new middleware usage to the controller.
func (c *HTTPController) UseMiddleware(f func(http.Handler) http.Handler) { c.r.Use(f) }

// ListenAndServe starts listening on the received address.
func (c *HTTPController) ListenAndServe(addr string) { http.ListenAndServe(addr, c.r) }
