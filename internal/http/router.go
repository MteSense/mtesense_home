package httpapi

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"html"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"mtesense_home/internal/auth"
	"mtesense_home/internal/config"
	"mtesense_home/internal/nav"
	"mtesense_home/internal/settings"
	"mtesense_home/internal/storage"
)

type Server struct {
	cfg      config.Config
	auth     *auth.Service
	nav      *nav.Service
	settings *settings.Service
	storage  *storage.Service
}

func NewRouter(cfg config.Config, database *sql.DB) http.Handler {
	server := &Server{
		cfg:      cfg,
		auth:     auth.NewService(database, cfg.JWTSecret),
		nav:      nav.NewService(database),
		settings: settings.NewService(database),
		storage:  storage.NewService(cfg.UploadDir),
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors)

	r.Route("/api/v1", func(api chi.Router) {
		api.Use(noindex)
		api.Get("/health", server.health)
		api.Post("/auth/login", server.login)
		api.Get("/navigation", server.publicNavigation)
		api.Get("/settings", server.publicSettings)

		api.Group(func(protected chi.Router) {
			protected.Use(server.requireAuth)
			protected.Get("/me", server.me)

			protected.Route("/admin", func(admin chi.Router) {
				admin.Get("/navigation", server.adminNavigation)
				admin.Post("/groups", server.createGroup)
				admin.Put("/groups/{id}", server.updateGroup)
				admin.Delete("/groups/{id}", server.deleteGroup)
				admin.Post("/links", server.createLink)
				admin.Put("/links/{id}", server.updateLink)
				admin.Delete("/links/{id}", server.deleteLink)
				admin.Put("/settings", server.saveSettings)
				admin.Post("/uploads", server.upload)
			})
		})
	})

	r.Get("/robots.txt", server.robots)
	r.Get("/sitemap.xml", server.sitemap)
	r.Handle("/uploads/*", noindex(http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.UploadDir)))))
	r.NotFound(server.spaHandler(filepath.Join("web", "app", "dist")))
	return r
}

func (s *Server) EnsureAdmin(username, password string) error {
	return s.auth.EnsureAdmin(username, password)
}

func (s *Server) EnsureStorage() error {
	return s.storage.Ensure()
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func noindex(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Robots-Tag", "noindex, nofollow")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r.Header.Get("Authorization"))
		if token == "" {
			writeError(w, http.StatusUnauthorized, "missing token")
			return
		}
		user, err := s.auth.ParseToken(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}
		next.ServeHTTP(w, r.WithContext(auth.ContextWithUser(r.Context(), user)))
	})
}

func bearerToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	token, user, err := s.auth.Login(payload.Username, payload.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": token, "user": user})
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing user")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (s *Server) publicNavigation(w http.ResponseWriter, r *http.Request) {
	data, err := s.nav.PublicNavigation()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (s *Server) adminNavigation(w http.ResponseWriter, r *http.Request) {
	data, err := s.nav.AdminNavigation()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (s *Server) publicSettings(w http.ResponseWriter, r *http.Request) {
	data, err := s.settings.GetPublic()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (s *Server) saveSettings(w http.ResponseWriter, r *http.Request) {
	var payload settings.PublicSettings
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	saved, err := s.settings.SavePublic(payload)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, saved)
}

func (s *Server) createGroup(w http.ResponseWriter, r *http.Request) {
	var payload nav.Group
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if strings.TrimSpace(payload.Title) == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	created, err := s.nav.CreateGroup(payload)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) updateGroup(w http.ResponseWriter, r *http.Request) {
	id, ok := paramID(w, r)
	if !ok {
		return
	}
	var payload nav.Group
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if strings.TrimSpace(payload.Title) == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	updated, err := s.nav.UpdateGroup(id, payload)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (s *Server) deleteGroup(w http.ResponseWriter, r *http.Request) {
	id, ok := paramID(w, r)
	if !ok {
		return
	}
	if err := s.nav.DeleteGroup(id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}

func (s *Server) createLink(w http.ResponseWriter, r *http.Request) {
	var payload nav.Link
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := validateLink(payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := s.nav.CreateLink(payload)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) updateLink(w http.ResponseWriter, r *http.Request) {
	id, ok := paramID(w, r)
	if !ok {
		return
	}
	var payload nav.Link
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := validateLink(payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	updated, err := s.nav.UpdateLink(id, payload)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (s *Server) deleteLink(w http.ResponseWriter, r *http.Request) {
	id, ok := paramID(w, r)
	if !ok {
		return
	}
	if err := s.nav.DeleteLink(id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}

func (s *Server) upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(6 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "invalid multipart form")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file field is required")
		return
	}
	defer file.Close()
	saved, err := s.storage.Save(file, header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, saved)
}

func (s *Server) robots(w http.ResponseWriter, r *http.Request) {
	baseURL := s.publicBaseURL(r)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("User-agent: *\nAllow: /\n\nSitemap: " + baseURL + "/sitemap.xml\n"))
}

func (s *Server) sitemap(w http.ResponseWriter, r *http.Request) {
	type urlEntry struct {
		Loc string `xml:"loc"`
	}
	type urlSet struct {
		XMLName xml.Name   `xml:"urlset"`
		XMLNS   string     `xml:"xmlns,attr"`
		URLs    []urlEntry `xml:"url"`
	}

	body, err := xml.MarshalIndent(urlSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  []urlEntry{{Loc: s.publicBaseURL(r) + "/"}},
	}, "", "  ")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "build sitemap")
		return
	}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	_, _ = w.Write([]byte(xml.Header))
	_, _ = w.Write(body)
	_, _ = w.Write([]byte("\n"))
}

func paramID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	raw := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return id, true
}

func validateLink(link nav.Link) error {
	if link.GroupID <= 0 {
		return errors.New("groupId is required")
	}
	if strings.TrimSpace(link.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(link.URL) == "" {
		return errors.New("url is required")
	}
	if link.IconType == "" {
		link.IconType = "emoji"
	}
	return nil
}

func (s *Server) spaHandler(dist string) http.HandlerFunc {
	fileServer := http.FileServer(http.Dir(dist))
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		indexPath := filepath.Join(dist, "index.html")
		if _, err := os.Stat(indexPath); err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("MteSense Home API is running. Build the frontend with `npm run build` in web/app to serve the UI."))
			return
		}
		requested := filepath.Join(dist, filepath.Clean(r.URL.Path))
		if stat, err := os.Stat(requested); err == nil && !stat.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}
		if _, err := fs.Stat(os.DirFS(dist), "index.html"); err == nil {
			s.serveIndex(w, r, indexPath)
			return
		}
		writeError(w, http.StatusNotFound, "not found")
	}
}

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request, indexPath string) {
	if strings.HasPrefix(r.URL.Path, "/admin") {
		w.Header().Set("X-Robots-Tag", "noindex, nofollow")
	}

	indexHTML, err := os.ReadFile(indexPath)
	if err != nil {
		http.ServeFile(w, r, indexPath)
		return
	}

	settings, err := s.settings.GetPublic()
	if err != nil {
		http.ServeFile(w, r, indexPath)
		return
	}

	title := strings.TrimSpace(settings.Appearance.BrowserTitle)
	if title == "" {
		title = strings.TrimSpace(settings.Appearance.SiteTitle)
	}
	if title == "" {
		title = "MteSense"
	}
	description := strings.TrimSpace(settings.Appearance.Subtitle)
	if description == "" {
		description = "Personal navigation"
	}

	body := string(indexHTML)
	body = replaceTitle(body, title)
	body = injectHeadTags(body, buildSEOTags(title, description, s.publicBaseURL(r), strings.HasPrefix(r.URL.Path, "/admin")))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(body))
}

func replaceTitle(document string, title string) string {
	start := strings.Index(document, "<title>")
	end := strings.Index(document, "</title>")
	if start == -1 || end == -1 || end < start {
		return document
	}
	return document[:start] + "<title>" + html.EscapeString(title) + "</title>" + document[end+len("</title>"):]
}

func injectHeadTags(document string, tags string) string {
	headEnd := strings.Index(document, "</head>")
	if headEnd == -1 {
		return document
	}
	return document[:headEnd] + tags + document[headEnd:]
}

func buildSEOTags(title, description, baseURL string, noindex bool) string {
	canonical := baseURL + "/"
	var b strings.Builder
	b.WriteString("    <meta name=\"description\" content=\"")
	b.WriteString(html.EscapeString(description))
	b.WriteString("\" />\n")
	if noindex {
		b.WriteString("    <meta name=\"robots\" content=\"noindex,nofollow\" />\n")
	} else {
		b.WriteString("    <link rel=\"canonical\" href=\"")
		b.WriteString(html.EscapeString(canonical))
		b.WriteString("\" />\n")
	}
	b.WriteString("    <meta property=\"og:type\" content=\"website\" />\n")
	b.WriteString("    <meta property=\"og:title\" content=\"")
	b.WriteString(html.EscapeString(title))
	b.WriteString("\" />\n")
	b.WriteString("    <meta property=\"og:description\" content=\"")
	b.WriteString(html.EscapeString(description))
	b.WriteString("\" />\n")
	b.WriteString("    <meta property=\"og:url\" content=\"")
	b.WriteString(html.EscapeString(canonical))
	b.WriteString("\" />\n")
	b.WriteString("    <meta name=\"twitter:card\" content=\"summary\" />\n")
	return b.String()
}

func (s *Server) publicBaseURL(r *http.Request) string {
	if configured := strings.TrimRight(strings.TrimSpace(s.cfg.PublicSiteURL), "/"); configured != "" {
		return configured
	}
	proto := forwardedValue(r.Header.Get("X-Forwarded-Proto"))
	if proto == "" {
		if r.TLS != nil {
			proto = "https"
		} else {
			proto = "http"
		}
	}
	host := forwardedValue(r.Header.Get("X-Forwarded-Host"))
	if host == "" {
		host = r.Host
	}
	return proto + "://" + strings.TrimRight(host, "/")
}

func forwardedValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	parts := strings.Split(value, ",")
	return strings.TrimSpace(parts[0])
}
