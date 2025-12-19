package core

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"time"

	"Prak_10/internal/http/middleware"
)

type userRepo interface {
	CheckPassword(email, pass string) (User, error)
}
type jwtSigner interface {
	Sign(userID int64, email, role string, ttl time.Duration) (string, error)
	Parse(token string) (jwt.MapClaims, error)
}

type Service struct {
	repo userRepo
	jwt  jwtSigner

	// для refresh токенов
	refreshBlacklist map[string]int64 // map[token]expUnix
}

func NewService(r userRepo, j jwtSigner) *Service {
	return &Service{
		repo:             r,
		jwt:              j,
		refreshBlacklist: make(map[string]int64),
	}
}

// структура ответа при логине
type loginResp struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func (s *Service) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var in struct{ Email, Password string }
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Email == "" || in.Password == "" {
		httpError(w, 400, "invalid_credentials")
		return
	}
	u, err := s.repo.CheckPassword(in.Email, in.Password)
	if err != nil {
		httpError(w, 401, "unauthorized")
		return
	}
	access, err := s.jwt.Sign(u.ID, u.Email, u.Role, 15*time.Minute) // 15 мин
	if err != nil {
		httpError(w, 500, "token_error")
		return
	}
	refresh, err := s.jwt.Sign(u.ID, u.Email, u.Role, 7*24*time.Hour) // 7 дней
	if err != nil {
		httpError(w, 500, "token_error")
		return
	}

	jsonOK(w, loginResp{Access: access, Refresh: refresh})
}

func (s *Service) MeHandler(w http.ResponseWriter, r *http.Request) {
	// клеймы положим в контекст в AuthN-мидлваре
	claims := r.Context().Value(middleware.CtxClaimsKey).(jwt.MapClaims)
	jsonOK(w, map[string]any{
		"id": claims["sub"], "email": claims["email"], "role": claims["role"],
	})
}

func (s *Service) AdminStats(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, map[string]any{"users": 2, "version": "1.0"})
}

// утилиты и ключ для контекста — экспортируем из middleware
// type ctxClaims struct{}
func jsonOK(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
func httpError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (s *Service) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var in struct{ Refresh string }
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Refresh == "" {
		httpError(w, 400, "invalid_request")
		return
	}

	// проверка blacklist
	if exp, ok := s.refreshBlacklist[in.Refresh]; ok && exp > time.Now().Unix() {
		httpError(w, 401, "token_revoked")
		return
	}

	claims, err := s.jwt.Parse(in.Refresh)
	if err != nil {
		httpError(w, 401, "invalid_token")
		return
	}

	userID := int64(claims["sub"].(float64))
	email := claims["email"].(string)
	role := claims["role"].(string)

	access, _ := s.jwt.Sign(userID, email, role, 15*time.Minute)
	refresh, _ := s.jwt.Sign(userID, email, role, 7*24*time.Hour)

	s.refreshBlacklist[in.Refresh] = int64(claims["exp"].(float64))

	jsonOK(w, loginResp{Access: access, Refresh: refresh})
}

func (s *Service) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.CtxClaimsKey).(jwt.MapClaims)
	role := claims["role"].(string)
	sub := int64(claims["sub"].(float64))

	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if role == "user" && sub != id {
		httpError(w, 403, "forbidden")
		return
	}

	// моковые данные
	users := map[int64]User{
		1: {ID: 1, Email: "admin@example.com", Role: "admin"},
		2: {ID: 2, Email: "user@example.com", Role: "user"},
	}
	u, ok := users[id]
	if !ok {
		httpError(w, 404, "not_found")
		return
	}

	jsonOK(w, u)
}
