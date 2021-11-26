package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	lru "github.com/hashicorp/golang-lru"
	"hw9/internal/message_broker"
	"hw9/internal/models"
	"hw9/internal/store"
	"net/http"
	"strconv"
)

type JobResource struct {
	store store.Store
	broker message_broker.MessageBroker
	cache *lru.TwoQueueCache
}

func NewJobResource(store store.Store, broker message_broker.MessageBroker,cache *lru.TwoQueueCache) *JobResource {
	return &JobResource{
		store: store,
		broker: broker,
		cache: cache,
	}
}
func (jr *JobResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", jr.CreateJob)
	r.Get("/", jr.AllJobs)
	r.Get("/{id}", jr.ByID)
	r.Put("/", jr.UpdateJob)
	r.Delete("/{id}", jr.DeleteJob)

	return r
}
func (jr *JobResource) CreateJob(w http.ResponseWriter, r *http.Request) {
	job := new(models.Job)
	if err := json.NewDecoder(r.Body).Decode(job); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := jr.store.Jobs().Create(r.Context(), job); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	// Правильно пройтись по всем буквам и всем словам
	jr.broker.Cache().Purge() // в рамках учебного проекта полностью чистим кэш после создания новой категории

	w.WriteHeader(http.StatusCreated)
}

func (jr *JobResource) AllJobs(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.JobsFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		jobsFromCache, ok := jr.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, jobsFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	jobs, err := jr.store.Jobs().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if searchQuery != "" {
		jr.cache.Add(searchQuery, jobs)
	}
	render.JSON(w, r, jobs)
}

func (jr *JobResource) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	jobFromCache, ok := jr.cache.Get(id)
	if ok {
		render.JSON(w, r, jobFromCache)
		return
	}

	job, err := jr.store.Jobs().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	jr.cache.Add(id, job)
	render.JSON(w, r, job)
}

func (jr *JobResource) UpdateJob(w http.ResponseWriter, r *http.Request) {
	job := new(models.Job)
	if err := json.NewDecoder(r.Body).Decode(job); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	//err := validation.ValidateStruct(
	//	job,
	//	validation.Field(&job.ID, validation.Required),
	//	validation.Field(&job.Name, validation.Required),
	//)
	//if err != nil {
	//	w.WriteHeader(http.StatusUnprocessableEntity)
	//	fmt.Fprintf(w, "Unknown err: %v", err)
	//	return
	//}

	if err := jr.store.Jobs().Update(r.Context(), job); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	jr.broker.Cache().Remove(job.ID)
}

func (jr *JobResource) DeleteJob(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := jr.store.Jobs().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	jr.broker.Cache().Remove(id)
}