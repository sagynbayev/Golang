package http

import (
	"context"
	"encoding/json"
	"fmt"
	"hw8/internal/models"
	"hw8/internal/store"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	store       store.Store

	Address string
}

func NewServer(ctx context.Context, address string, store store.Store) *Server {
	return &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
		store:       store,

		Address: address,
	}
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	r.Post("/categories", func(w http.ResponseWriter, r *http.Request) {
		category := new(models.Category)
		if err := json.NewDecoder(r.Body).Decode(category); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Categories().Create(r.Context(), category); err != nil {
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})
	r.Get("/categories", func(w http.ResponseWriter, r *http.Request) {
		categories, err := s.store.Categories().All(r.Context())
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, categories)
	})
	r.Get("/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		category, err := s.store.Categories().ByID(r.Context(), id)
		if err != nil {
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}

		render.JSON(w, r, category)
	})
	r.Put("/categories", func(w http.ResponseWriter, r *http.Request) {
		category := new(models.Category)
		if err := json.NewDecoder(r.Body).Decode(category); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Categories().Update(r.Context(), category); err != nil{
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})
	r.Delete("/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Categories().Delete(r.Context(), id); err != nil {
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})
	r.Post("/jobs", func(w http.ResponseWriter, r *http.Request) {
		job := new(models.Job)
		if err := json.NewDecoder(r.Body).Decode(job); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Jobs().Create(r.Context(), job); err != nil {
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})
	r.Get("/jobs", func(w http.ResponseWriter, r *http.Request) {
		jobs, err := s.store.Jobs().All(r.Context())
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, jobs)
	})
	r.Get("/jobs/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		job, err := s.store.Jobs().ByID(r.Context(), id)
		if err != nil {
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}

		render.JSON(w, r, job)
	})
	r.Put("/jobs", func(w http.ResponseWriter, r *http.Request) {
		job := new(models.Job)
		if err := json.NewDecoder(r.Body).Decode(job); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Jobs().Update(r.Context(), job); err != nil{
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})
	r.Delete("/jobs/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Jobs().Delete(r.Context(), id); err != nil {
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})
	return r
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}
	go s.ListenCtxForGT(srv)

	log.Println("[HTTP] Server running on", s.Address)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done() // блокируемся, пока контекст приложения не отменен

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("[HTTP] Got err while shutting down^ %v", err)
	}

	log.Println("[HTTP] Proccessed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	// блок до записи или закрытия канала
	<-s.idleConnsCh
}
