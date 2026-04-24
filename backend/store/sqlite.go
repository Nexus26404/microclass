package store

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"microclass-backend/model"
)

type Store struct {
	db *sql.DB
}

func New(dbPath string) (*Store, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	s := &Store{db: db}
	if err := s.init(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) init() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS scripts (
			id TEXT PRIMARY KEY,
			topic TEXT NOT NULL,
			discipline TEXT,
			target_audience TEXT,
			style TEXT,
			estimated_total_duration INTEGER,
			sections TEXT,
			created_at DATETIME,
			updated_at DATETIME
		)
	`)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS settings (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			mock_mode INTEGER NOT NULL DEFAULT 1,
			open_ai_key TEXT,
			open_ai_url TEXT,
			open_ai_model TEXT
		)
	`)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`INSERT OR IGNORE INTO settings (id) VALUES (1)`)
	return err
}

func (s *Store) Save(script *model.Script) error {
	sectionsJSON, _ := json.Marshal(script.Sections)
	_, err := s.db.Exec(`
		INSERT INTO scripts (id, topic, discipline, target_audience, style, estimated_total_duration, sections, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, script.ID, script.Topic, script.Discipline, script.TargetAudience, script.Style, script.EstimatedTotalDuration, string(sectionsJSON), script.CreatedAt, script.UpdatedAt)
	return err
}

func (s *Store) GetAll() ([]model.Script, error) {
	rows, err := s.db.Query(`
		SELECT id, topic, discipline, target_audience, style, estimated_total_duration, sections, created_at, updated_at
		FROM scripts ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scripts []model.Script
	for rows.Next() {
		var sc model.Script
		var sectionsJSON string
		if err := rows.Scan(&sc.ID, &sc.Topic, &sc.Discipline, &sc.TargetAudience, &sc.Style, &sc.EstimatedTotalDuration, &sectionsJSON, &sc.CreatedAt, &sc.UpdatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(sectionsJSON), &sc.Sections)
		scripts = append(scripts, sc)
	}
	return scripts, nil
}

func (s *Store) GetByID(id string) (*model.Script, error) {
	var sc model.Script
	var sectionsJSON string
	err := s.db.QueryRow(`
		SELECT id, topic, discipline, target_audience, style, estimated_total_duration, sections, created_at, updated_at
		FROM scripts WHERE id = ?
	`, id).Scan(&sc.ID, &sc.Topic, &sc.Discipline, &sc.TargetAudience, &sc.Style, &sc.EstimatedTotalDuration, &sectionsJSON, &sc.CreatedAt, &sc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(sectionsJSON), &sc.Sections)
	return &sc, nil
}

func (s *Store) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM scripts WHERE id = ?", id)
	return err
}

func (s *Store) GetSettings() (mockMode bool, openAIKey, openAIURL, openAIModel string, err error) {
	err = s.db.QueryRow(`
		SELECT mock_mode, open_ai_key, open_ai_url, open_ai_model FROM settings WHERE id = 1
	`).Scan(&mockMode, &openAIKey, &openAIURL, &openAIModel)
	return
}

func (s *Store) SaveSettings(mockMode bool, openAIKey, openAIURL, openAIModel string) error {
	_, err := s.db.Exec(`
		UPDATE settings SET mock_mode = ?, open_ai_key = ?, open_ai_url = ?, open_ai_model = ? WHERE id = 1
	`, mockMode, openAIKey, openAIURL, openAIModel)
	return err
}

func (s *Store) Close() error {
	return s.db.Close()
}
