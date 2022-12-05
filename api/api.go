package api

import (
	"fmt"
	"log"
	"net/http"

	"btc_billionaire/postgres"

	"github.com/gorilla/mux"
)

type Handler struct {
	router *mux.Router
	store  postgres.Store
}

// type Handler struct {
// 	store *postgres.Storage
// }

// func New(db *postgres.Storage) *Handler {
// 	return &Handler{
// 		store: db,
// 	}
// }

func NewServer(db postgres.Store) *Handler {
	h := &Handler{
		router: mux.NewRouter(),
		store:  db,
	}
	h.setRoutes()
	return h
}

func (h *Handler) StartServer(host, port string) {
	fmt.Printf("Server started on http://%v:%v\n", host, port)
	log.Fatal(http.ListenAndServe(":"+port, h.router))
}

func (h *Handler) setRoutes() {
	h.Post("/api/btc", setMiddlewareJSON(h.SendBTC))
	h.Post("/api/btc/history", setMiddlewareJSON(h.GetBTCHistory))
}

func (h *Handler) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	h.router.HandleFunc(path, f).Methods("GET")
}

func (h *Handler) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	h.router.HandleFunc(path, f).Methods("POST")
}

func setMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}
