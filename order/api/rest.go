package api

import (
	"errors"
	"fmt"
	"l0/order/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"l0/cache"
)

type API struct {
	orderStore repository.Querier
	cache      *cache.Cache
}

func (a *API) NewRouter(orderStore repository.Querier, cache *cache.Cache) chi.Router {
	a.orderStore = orderStore
	a.cache = cache

	r := chi.NewRouter()
	r.Get("/order/{id}", a.getOrderByID)

	return r
}

// GET /order/{id} - returns order by id
func (a *API) getOrderByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		SendErrorJSON(w, r, http.StatusBadRequest, errors.New("empty id"), "empty id")
		return
	}

	order, ok, err := a.cache.Get(r.Context(), id)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get order")
		return
	}
	if !ok {
		SendErrorJSON(w, r, http.StatusNotFound, fmt.Errorf("order %s not found", id), "order not found")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, order)
}
