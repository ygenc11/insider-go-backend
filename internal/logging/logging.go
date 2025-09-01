package logging

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var logFile *os.File

// Init, global slog logger'ı bir dosyaya yazacak şekilde yapılandırır.
// Çevre değişkenleri (ENV):
// - LOG_LEVEL: debug|info|warn|error (varsayılan: info)
// - LOG_FORMAT: json|text (varsayılan: json)
// - LOG_DIR: logların yazılacağı dizin (varsayılan: ./logs)
// - LOG_FILE: dosya adı; boşsa app-YYYY-MM-DD.log kullanılır
// - LOG_ADD_SOURCE: kaynak bilgisini eklemek için true|false (varsayılan: false)
func Init() {
	level := parseLevel(getenv("LOG_LEVEL", "info"))
	format := strings.ToLower(getenv("LOG_FORMAT", "json"))
	dir := getenv("LOG_DIR", "./logs")
	filename := getenv("LOG_FILE", "")
	addSource := strings.EqualFold(getenv("LOG_ADD_SOURCE", "false"), "true")

	// Log dizininin var olduğundan emin ol
	if err := os.MkdirAll(dir, 0o755); err != nil {
		// Oluşturma başarısız olursa mevcut dizine geri dön
		fmt.Fprintf(os.Stderr, "failed to create log dir %s: %v\n", dir, err)
		dir = "."
	}

	if filename == "" {
		filename = fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02"))
	}
	path := filepath.Join(dir, filename)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open log file %s: %v, falling back to stderr\n", path, err)
	} else {
		logFile = f
	}

	var out *os.File = os.Stderr
	if logFile != nil {
		out = logFile
	}

	opts := &slog.HandlerOptions{Level: level, AddSource: addSource}
	var handler slog.Handler
	switch format {
	case "text":
		handler = slog.NewTextHandler(out, opts)
	default:
		handler = slog.NewJSONHandler(out, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// Log dosyasını kapatır eğer açıksa
func Close() {
	if logFile != nil {
		_ = logFile.Sync()
		_ = logFile.Close()
		logFile = nil
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func parseLevel(s string) slog.Leveler {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
