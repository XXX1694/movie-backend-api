package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang/internal/usecase"
	"golang/pkg/modules"

	"github.com/gorilla/mux"
)

type MovieHandler struct {
	usecase usecase.MovieUsecaseInterface
}

func NewMovieHandler(uc usecase.MovieUsecaseInterface) *MovieHandler {
	return &MovieHandler{usecase: uc}
}

// GetMovies godoc
// @Summary Get all movies
// @Tags movies
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} modules.Movie
// @Router /movies [get]
// @Security ApiKeyAuth
func (h *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	limit := 10
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}

	movies, err := h.usecase.GetMovies(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// GetMovieByID godoc
// @Summary Get movie by ID
// @Tags movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} modules.Movie
// @Failure 404 {string} string "not found"
// @Router /movies/{id} [get]
// @Security ApiKeyAuth
func (h *MovieHandler) GetMovieByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	movie, err := h.usecase.GetMovieByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

// CreateMovie godoc
// @Summary Create a new movie
// @Tags movies
// @Accept json
// @Produce json
// @Param movie body modules.Movie true "Movie"
// @Success 201 {object} map[string]int
// @Router /movies [post]
// @Security ApiKeyAuth
func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie modules.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	id, err := h.usecase.CreateMovie(movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// UpdateMovie godoc
// @Summary Update movie
// @Tags movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Param movie body modules.Movie true "Movie"
// @Success 200 {object} map[string]string
// @Failure 404 {string} string "not found"
// @Router /movies/{id} [put]
// @Security ApiKeyAuth
func (h *MovieHandler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var movie modules.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.usecase.UpdateMovie(id, movie); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "updated"})
}

// DeleteMovie godoc
// @Summary Delete movie
// @Tags movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} map[string]int64
// @Failure 404 {string} string "not found"
// @Router /movies/{id} [delete]
// @Security ApiKeyAuth
func (h *MovieHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	rows, err := h.usecase.DeleteMovie(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"rows_affected": rows})
}
