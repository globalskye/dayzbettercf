package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"dayzsmartcf/backend/internal/auth"
	"dayzsmartcf/backend/internal/config"
	"dayzsmartcf/backend/internal/cftools"
	"dayzsmartcf/backend/internal/handlers"
	"dayzsmartcf/backend/internal/player"
)

type Server struct {
	cfg           *config.Config
	router        chi.Router
	cftoolsClient *cftools.Client
	repo          *player.Repository
	syncSvc       *player.SyncService
	authRepo      *auth.Repo
}

func New(cfg *config.Config, cf *cftools.Client, repo *player.Repository, syncSvc *player.SyncService, authRepo *auth.Repo) *Server {
	s := &Server{
		cfg:           cfg,
		cftoolsClient: cf,
		repo:          repo,
		syncSvc:       syncSvc,
		authRepo:      authRepo,
	}
	s.setupRouter(repo, syncSvc)
	return s
}

func (s *Server) setupRouter(repo *player.Repository, syncSvc *player.SyncService) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173", "http://localhost:3000", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", handlers.Health)
	r.Get("/api/v1/hello", handlers.Hello)

	// Auth — public
	r.Post("/api/v1/auth/login", handlers.AuthLogin(s.authRepo, s.cfg.JWTSecret))

	// Protected API
	requireAuth := auth.RequireAuth(s.cfg.JWTSecret, s.authRepo)
	requireAdmin := auth.RequireRole(auth.RoleAdmin)

	r.Group(func(r chi.Router) {
		r.Use(requireAuth)
		r.Use(auth.LogRequests(s.authRepo))
		r.Get("/api/v1/auth/me", handlers.AuthMe())

		r.Get("/api/v1/cftools/status", handlers.CFtoolsStatus(s.cftoolsClient))
		r.Get("/api/v1/cftools/states", handlers.CFToolsStates(s.cftoolsClient))

		r.Route("/api/v1/players", func(r chi.Router) {
			r.Get("/", handlers.PlayersList(repo))
			r.Get("/search", handlers.PlayersSearchLocal(repo))
			r.Get("/search-cf", handlers.PlayersSearch(syncSvc, repo))
			r.Get("/cftools-search", handlers.PlayersSearchCFtools(s.cftoolsClient))
			r.Post("/sync-batch", handlers.PlayersSyncBatch(syncSvc))
			r.Get("/{id}", handlers.PlayersGet(repo))
			r.Get("/{id}/history", handlers.PlayerHistory(repo))
			r.Post("/{id}/sync", handlers.PlayersSyncOne(syncSvc, repo))
		})
		r.Route("/api/v1/tracked", func(r chi.Router) {
			r.Get("/", handlers.TrackedList(repo, syncSvc))
			r.Post("/add/{cftoolsId}", handlers.TrackedAdd(repo, syncSvc))
			r.Delete("/remove/{cftoolsId}", handlers.TrackedRemove(repo))
			r.Get("/{cftoolsId}/history", handlers.TrackedHistory(repo))
		})
		r.Route("/api/v1/groups", func(r chi.Router) {
			r.Get("/", handlers.GroupsList(repo, syncSvc))
			r.Post("/create/{name}", handlers.GroupsCreate(repo))
			r.Get("/{id}", handlers.GroupsGet(repo, syncSvc))
			r.Delete("/{id}", handlers.GroupsDelete(repo))
			r.Post("/{id}/add/{cftoolsId}", handlers.GroupsAddMember(repo, syncSvc))
			r.Patch("/{id}/members/{cftoolsId}", handlers.GroupsUpdateMemberAlias(repo, syncSvc))
			r.Delete("/{id}/remove/{cftoolsId}", handlers.GroupsRemoveMember(repo))
		})

		// Settings + Admin — admin only
		r.Group(func(r chi.Router) {
			r.Use(requireAdmin)
			r.Post("/api/v1/cftools/login", handlers.CFtoolsLogin(s.cftoolsClient))
			r.Get("/api/v1/settings/auth", handlers.AuthSettingsGet(s.cftoolsClient))
			r.Get("/api/v1/settings/auth/check", handlers.AuthSettingsCheck(s.cftoolsClient))
			r.Post("/api/v1/settings/auth", handlers.AuthSettingsUpdate(s.cftoolsClient, s.cfg))
			r.Post("/api/v1/settings/db/wipe", handlers.DBWipe(repo))
			r.Get("/api/v1/admin/users", handlers.AdminListUsers(s.authRepo))
			r.Post("/api/v1/admin/users", handlers.AdminCreateUser(s.authRepo))
			r.Patch("/api/v1/admin/users/{id}", handlers.AdminUpdateUser(s.authRepo))
			r.Delete("/api/v1/admin/users/{id}", handlers.AdminDeleteUser(s.authRepo))
			r.Get("/api/v1/admin/users/{id}/logs", handlers.AdminGetRequestLogs(s.authRepo))
		})
	})

	s.router = r
}

func (s *Server) Router() http.Handler {
	return s.router
}
