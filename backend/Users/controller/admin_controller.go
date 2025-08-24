package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gajare/Fish-market/db"
	"github.com/gajare/Fish-market/logger"
	"github.com/sirupsen/logrus"
)

type LogLevelDTO struct {
	Level    string `json:"level"`               // debug|info|warn|error|fatal|panic|trace
	SQLLevel string `json:"sql_level,omitempty"` // silent|error|warn|info
	SlowMS   *int   `json:"slow_ms,omitempty"`   // optional: gorm slow threshold
}

type AdminController struct{}

func NewAdminController() *AdminController { return &AdminController{} }

func (a *AdminController) GetLogLevel(w http.ResponseWriter, r *http.Request) {
	resp := map[string]any{
		"level": logger.GetLevel().String(),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (a *AdminController) SetLogLevel(w http.ResponseWriter, r *http.Request) {
	var dto LogLevelDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if dto.Level != "" {
		l, err := logrus.ParseLevel(strings.ToLower(dto.Level))
		if err != nil {
			http.Error(w, "invalid level", http.StatusBadRequest)
			return
		}
		logger.SetLevel(l)
		logger.With(map[string]any{"level": l.String()}).Info("app_log_level_updated")
	}
	if dto.SQLLevel != "" || dto.SlowMS != nil {
		slow := db.EnvSlowMS()
		if dto.SlowMS != nil {
			slow = *dto.SlowMS
		}
		db.SetSQLLogLevel(strings.ToLower(dto.SQLLevel), slow)
	}
	w.WriteHeader(http.StatusNoContent)
}
