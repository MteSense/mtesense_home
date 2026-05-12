package settings

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Appearance struct {
	SiteTitle    string  `json:"siteTitle"`
	BrowserTitle string  `json:"browserTitle"`
	Subtitle     string  `json:"subtitle"`
	FooterHTML   string  `json:"footerHtml"`
	Background   string  `json:"backgroundImage"`
	DefaultTheme string  `json:"defaultTheme"`
	CardOpacity  float64 `json:"cardOpacity"`
	BlurStrength int     `json:"blurStrength"`
}

type Search struct {
	DefaultSearchEngine  string   `json:"defaultSearchEngine"`
	EnabledSearchEngines []string `json:"enabledSearchEngines"`
}

type PublicSettings struct {
	Appearance Appearance `json:"appearance"`
	Search     Search     `json:"search"`
}

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetPublic() (PublicSettings, error) {
	appearance := Appearance{
		SiteTitle:    "MteSense",
		BrowserTitle: "MteSense",
		Subtitle:     "Personal navigation",
		FooterHTML:   "© 2025 - 2026 MteSense Studio. All rights reserved.",
		DefaultTheme: "dark",
		CardOpacity:  0.34,
		BlurStrength: 18,
	}
	search := Search{
		DefaultSearchEngine:  "google",
		EnabledSearchEngines: []string{"google", "bing", "baidu"},
	}
	if err := s.getJSON("appearance", &appearance); err != nil {
		return PublicSettings{}, err
	}
	if appearance.BrowserTitle == "" {
		appearance.BrowserTitle = appearance.SiteTitle
	}
	if err := s.getJSON("search", &search); err != nil {
		return PublicSettings{}, err
	}
	return PublicSettings{Appearance: appearance, Search: normalizeSearch(search)}, nil
}

func (s *Service) SavePublic(settings PublicSettings) (PublicSettings, error) {
	settings.Search = normalizeSearch(settings.Search)
	if settings.Appearance.SiteTitle == "" {
		settings.Appearance.SiteTitle = "MteSense"
	}
	if settings.Appearance.BrowserTitle == "" {
		settings.Appearance.BrowserTitle = settings.Appearance.SiteTitle
	}
	if settings.Appearance.FooterHTML == "" {
		settings.Appearance.FooterHTML = "© 2025 - 2026 MteSense Studio. All rights reserved."
	}
	if settings.Appearance.DefaultTheme != "light" && settings.Appearance.DefaultTheme != "dark" {
		settings.Appearance.DefaultTheme = "dark"
	}
	if settings.Appearance.CardOpacity <= 0 || settings.Appearance.CardOpacity > 1 {
		settings.Appearance.CardOpacity = 0.34
	}
	if settings.Appearance.BlurStrength < 0 || settings.Appearance.BlurStrength > 40 {
		settings.Appearance.BlurStrength = 18
	}
	if err := s.saveJSON("appearance", settings.Appearance); err != nil {
		return PublicSettings{}, err
	}
	if err := s.saveJSON("search", settings.Search); err != nil {
		return PublicSettings{}, err
	}
	return settings, nil
}

func (s *Service) getJSON(key string, target any) error {
	var raw string
	err := s.db.QueryRow("SELECT value_json FROM settings WHERE key = ?", key).Scan(&raw)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("read setting %s: %w", key, err)
	}
	if err := json.Unmarshal([]byte(raw), target); err != nil {
		return fmt.Errorf("decode setting %s: %w", key, err)
	}
	return nil
}

func (s *Service) saveJSON(key string, value any) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("encode setting %s: %w", key, err)
	}
	_, err = s.db.Exec(
		`INSERT INTO settings (key, value_json, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(key) DO UPDATE SET value_json = excluded.value_json, updated_at = CURRENT_TIMESTAMP`,
		key,
		string(raw),
	)
	if err != nil {
		return fmt.Errorf("save setting %s: %w", key, err)
	}
	return nil
}

func normalizeSearch(search Search) Search {
	allowed := map[string]bool{"google": true, "bing": true, "baidu": true}
	enabled := make([]string, 0, len(search.EnabledSearchEngines))
	for _, engine := range search.EnabledSearchEngines {
		if allowed[engine] {
			enabled = append(enabled, engine)
		}
	}
	if len(enabled) == 0 {
		enabled = []string{"google"}
	}
	defaultEngine := search.DefaultSearchEngine
	if !allowed[defaultEngine] {
		defaultEngine = enabled[0]
	}
	found := false
	for _, engine := range enabled {
		if engine == defaultEngine {
			found = true
			break
		}
	}
	if !found {
		defaultEngine = enabled[0]
	}
	return Search{DefaultSearchEngine: defaultEngine, EnabledSearchEngines: enabled}
}
