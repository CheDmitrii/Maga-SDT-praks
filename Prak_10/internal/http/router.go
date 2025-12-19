package http

import (
	"Prak_10/internal/http/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"

	"Prak_10/internal/core"
	"Prak_10/internal/platform/config"
	"Prak_10/internal/platform/jwt"
	"Prak_10/internal/repo"
)

func Build(cfg config.Config) http.Handler {
	r := chi.NewRouter()

	// DI
	userRepo := repo.NewUserMem() // храним заранее захэшированных юзеров (email, bcrypt)
	jwtv := jwt.NewHS256(cfg.JWTSecret, cfg.JWTTTL)
	svc := core.NewService(userRepo, jwtv)

	// Публичные маршруты
	r.Post("/api/v1/login", svc.LoginHandler) // выдаёт JWT по email+password
	r.Post("/api/v1/refresh", svc.RefreshHandler)

	// Защищённые маршруты
	r.Group(func(priv chi.Router) {
		priv.Use(middleware.AuthN(jwtv))                   // аутентификация JWT
		priv.Use(middleware.AuthZRoles("admin", "user"))   // базовая RBAC
		priv.Get("/api/v1/me", svc.MeHandler)              // вернёт профиль из токена
		priv.Get("/api/v1/users/{id}", svc.GetUserHandler) // маршрут ABAC
	})

	// Пример только для админов
	r.Group(func(admin chi.Router) {
		admin.Use(middleware.AuthN(jwtv))
		admin.Use(middleware.AuthZRoles("admin"))
		admin.Get("/api/v1/admin/stats", svc.AdminStats)
	})

	return r
}
