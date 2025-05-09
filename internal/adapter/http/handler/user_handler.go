package handler

import (
	"encoding/json"
	"net"
	"net/http"
	"net/mail"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"ip_detector/internal/app/service"
	"ip_detector/internal/domain/model"
	"ip_detector/internal/logger"
)

type registerRequest struct {
	Name     string `json:"name"     example:"John Doe"`
	Email    string `json:"email"    example:"john@example.com"`
	IP       string `json:"ip"       example:"8.8.8.8"`
	Password string `json:"password" example:"secret123"`
}

type loginRequest struct {
	Email    string `json:"email"    example:"john@example.com"`
	Password string `json:"password" example:"secret123"`
}

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// ---------------- Register ----------------

// RegisterUser godoc
// @Summary      User Registration
// @Description  Creates a new user, determines country by IP, hashes the password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      registerRequest  true  "User Registration Data"
// @Success      201      {object}  model.User
// @Failure      400,500  {string}  string
// @Router       /register [post]
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	log := logger.Log.Sugar()
	log.Infow("register request received")

	var input struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		IP       string `json:"ip" validate:"required,ip"`
		Password string `json:"password" validate:"required,min=6"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Warnw("invalid JSON", "error", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		log.Warnw("validation failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if net.ParseIP(input.IP) == nil {
		log.Warnw("invalid IP format", "ip", input.IP)
		http.Error(w, "invalid IP address", http.StatusBadRequest)
		return
	}

	if _, err := mail.ParseAddress(input.Email); err != nil {
		log.Warnw("invalid email format", "email", input.Email)
		http.Error(w, "invalid email address", http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorw("password hash error", "error", err)
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user := model.User{
		Name:         input.Name,
		Email:        input.Email,
		IP:           input.IP,
		PasswordHash: string(passwordHash),
	}

	if err := h.service.CreateUser(r.Context(), &user); err != nil {
		log.Errorw("create user failed", "email", user.Email, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infow("user registered", "id", user.ID, "email", user.Email)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

// ---------------- Login ----------------

// Login godoc
// @Summary      User Login
// @Description  Verifies user credentials and returns a JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      loginRequest  true  "User Login Data"
// @Success      200      {object}  map[string]string "token"
// @Failure      400,401,500  {string}  string
// @Router       /login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	log := logger.Log.Sugar()
	log.Infow("login request received")

	var credentials struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Warnw("invalid JSON", "error", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(credentials); err != nil {
		log.Warnw("validation failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserByEmail(r.Context(), credentials.Email)
	if err != nil {
		log.Errorw("DB error on login", "email", credentials.Email, "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		log.Warnw("user not found", "email", credentials.Email)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
		log.Warnw("password mismatch", "email", credentials.Email)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.service.GenerateJWT(user.Email)
	if err != nil {
		log.Errorw("token generation failed", "email", user.Email, "error", err)
		http.Error(w, "failed to generate JWT", http.StatusInternalServerError)
		return
	}

	log.Infow("login successful", "email", user.Email)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// ---------------- GetUsers ----------------

// GetUsers godoc
// @Summary      List Users
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   model.User
// @Failure      401,500  {string}  string
// @Router       /users [get]
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	log := logger.Log.Sugar()
	log.Infow("get users request")

	users, err := h.service.GetAllUsers(r.Context())
	if err != nil {
		log.Errorw("failed to fetch users", "error", err)
		http.Error(w, "failed to fetch users", http.StatusInternalServerError)
		return
	}

	log.Infow("users fetched", "count", len(users))
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(users)
}

// ---------------- GetUserByID ----------------

// GetUserByID godoc
// @Summary      Get User by ID
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  model.User
// @Failure      401,404,500  {string}  string
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log := logger.Log.Sugar()
	log.Infow("get user by id request", "id", id)

	user, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		log.Errorw("failed to fetch user", "id", id, "error", err)
		http.Error(w, "failed to fetch user", http.StatusInternalServerError)
		return
	}
	if user == nil {
		log.Infow("user not found", "id", id)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	log.Infow("user fetched", "id", id)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}
