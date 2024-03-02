package controllers

import (
	"github.com/go-chi/chi/v5"
	"github.com/kameikay/service-orchestration/internal/infra/web/handlers"
)

type Controller struct {
	router  chi.Router
	Handler *handlers.Handler
}

func NewController(
	router chi.Router,
	Handler *handlers.Handler,
) *Controller {
	return &Controller{
		router:  router,
		Handler: Handler,
	}
}

func (wc *Controller) Route() {
	wc.router.Route("/", func(r chi.Router) {
		r.Get("/", wc.Handler.GetTemperatures)
	})
}
