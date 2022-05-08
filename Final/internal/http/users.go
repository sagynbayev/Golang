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

type UserResource struct {
	store store.Store
	broker message_broker.MessageBroker
	cache *lru.TwoQueueCache
}

func NewUserResource(store store.Store, broker message_broker.MessageBroker,cache *lru.TwoQueueCache) *UserResource {
	return &UserResource{
		store: store,
		broker: broker,
		cache: cache,
	}
}
func (jr *UserResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", jr.CreateUser)
	r.Get("/", jr.AllUsers)
	r.Get("/{id}", jr.ByID)
	r.Put("/", jr.UpdateUser)
	r.Delete("/{id}", jr.DeleteUser)

	return r
}
func (jr *UserResource) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := jr.store.Users().Create(r.Context(), user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	// Правильно пройтись по всем буквам и всем словам
	jr.broker.Cache().Purge() // в рамках учебного проекта полностью чистим кэш после создания новой категории

	w.WriteHeader(http.StatusCreated)
}

func (jr *UserResource) AllUsers(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.UsersFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		usersFromCache, ok := jr.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, usersFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	users, err := jr.store.Users().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if searchQuery != "" {
		jr.cache.Add(searchQuery, users)
	}
	render.JSON(w, r, users)
}

func (jr *UserResource) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	userFromCache, ok := jr.cache.Get(id)
	if ok {
		render.JSON(w, r, userFromCache)
		return
	}

	user, err := jr.store.Users().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	jr.cache.Add(id, user)
	render.JSON(w, r, user)
}

func (jr *UserResource) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
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

	if err := jr.store.Users().Update(r.Context(), user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	jr.broker.Cache().Remove(user.ID)
}

func (jr *UserResource) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := jr.store.Users().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	jr.broker.Cache().Remove(id)
}