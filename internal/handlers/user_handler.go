package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pratham15541/go-crud/internal/models"
	"github.com/pratham15541/go-crud/internal/services"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		h.sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.sendSuccessResponse(w, "User created successfully", user.ToResponse(), http.StatusCreated)
}

// GetUser handles GET /users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		if err.Error() == "user not found" {
			h.sendErrorResponse(w, "User not found", http.StatusNotFound)
		} else {
			h.sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.sendSuccessResponse(w, "User retrieved successfully", user.ToResponse(), http.StatusOK)
}

// GetUsers handles GET /users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	users, total, err := h.userService.GetUsers(page, limit)
	if err != nil {
		h.sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	userResponses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	// Create paginated response
	response := map[string]interface{}{
		"users": userResponses,
		"pagination": map[string]interface{}{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	}

	h.sendSuccessResponse(w, "Users retrieved successfully", response, http.StatusOK)
}

// UpdateUser handles PUT /users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	user, err := h.userService.UpdateUser(id, &req)
	if err != nil {
		if err.Error() == "user not found" {
			h.sendErrorResponse(w, "User not found", http.StatusNotFound)
		} else {
			h.sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	h.sendSuccessResponse(w, "User updated successfully", user.ToResponse(), http.StatusOK)
}

// DeleteUser handles DELETE /users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.userService.DeleteUser(id)
	if err != nil {
		if err.Error() == "user not found" {
			h.sendErrorResponse(w, "User not found", http.StatusNotFound)
		} else {
			h.sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.sendSuccessResponse(w, "User deleted successfully", nil, http.StatusOK)
}

// sendErrorResponse sends an error response
func (h *UserHandler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResp := models.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    statusCode,
	}

	json.NewEncoder(w).Encode(errorResp)
}

// sendSuccessResponse sends a success response
func (h *UserHandler) sendSuccessResponse(w http.ResponseWriter, message string, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	successResp := models.SuccessResponse{
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(successResp)
}