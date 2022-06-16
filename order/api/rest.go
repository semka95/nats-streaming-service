package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"l0/order/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type API struct {
	orderStore repository.Querier
}

func (a *API) NewRouter(orderStore repository.Querier) chi.Router {
	a.orderStore = orderStore

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

	order, err := a.orderStore.GetOrderByID(r.Context(), json.RawMessage(id))
	if errors.Is(err, sql.ErrNoRows) {
		SendErrorJSON(w, r, http.StatusNotFound, err, "payment not found")
		return
	}
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get payment")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, order)
}
