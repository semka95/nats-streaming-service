package api

import (
	"errors"
	"fmt"
	"html/template"
	"l0/order/repository"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"l0/cache"
)

type API struct {
	orderStore repository.Querier
	cache      *cache.Cache
	tmpl       *template.Template
}

func (a *API) NewRouter(orderStore repository.Querier, c *cache.Cache) chi.Router {
	a.orderStore = orderStore
	a.cache = c

	workDir, _ := os.Getwd()
	a.tmpl = template.Must(template.ParseFiles(filepath.Join(workDir, "static", "templates", "index.gohtml"), filepath.Join(workDir, "static", "templates", "order.gohtml")))

	r := chi.NewRouter()
	r.Get("/order/{id}", a.getOrderByID)
	r.Post("/order", a.renderOrder)
	r.Get("/", a.renderIndex)

	root := http.Dir(filepath.Join(workDir, "static", "assets"))
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})

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

type orderInfo struct {
	Success bool
	Err     bool
	Order   string
	ErrMsg  string
}

// GET / - renders index page
func (a *API) renderIndex(w http.ResponseWriter, r *http.Request) {
	a.tmpl.Execute(w, nil)
}

// POST /order - renders order by id
func (a *API) renderOrder(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("order_uid")

	order, ok, err := a.cache.Get(r.Context(), id)
	if err != nil {
		a.tmpl.Execute(w, orderInfo{Err: true, ErrMsg: err.Error()})
		return
	}
	if !ok {
		a.tmpl.Execute(w, orderInfo{Err: true, ErrMsg: "order was not found"})
		return
	}

	a.tmpl.Execute(w, orderInfo{Success: true, Order: string(order)})
}
