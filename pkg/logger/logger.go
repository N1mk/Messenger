package logger

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const defaultLogFilePath = "logs/app.log"
const defaultOldLogDirPath = "logs/old"

type Logger struct {
	mainFilePath   string
	oldLogsDirPath string
	logger         *slog.Logger
	file           *os.File
}

func NewLogger(logFilePath string, oldLogsDirPath string) (*Logger, error) {
	if logFilePath == "" {
		logFilePath = defaultLogFilePath
	}

	if oldLogsDirPath == "" {
		oldLogsDirPath = defaultOldLogDirPath
	}

	dl := &Logger{mainFilePath: logFilePath, oldLogsDirPath: oldLogsDirPath}

	dl.rotate()

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		slog.Error("Cannot open log file")
		return nil, err
	}
	dl.file = file

	dl.logger = slog.New(slog.NewTextHandler(file, nil))

	return dl, nil
}

func (l *Logger) rotate() error {
	if err := os.MkdirAll(l.oldLogsDirPath, 0755); err != nil {
		return fmt.Errorf("old logs directory creation error: %w", err)
	}

	files, err := os.ReadDir(l.oldLogsDirPath)
	if err != nil {
		return fmt.Errorf("old logs directory read error: %w", err)
	}

	type fileInfo struct {
		entry   os.DirEntry
		modTime time.Time
	}
	var filteredFiles []fileInfo

	if err := os.Rename(l.mainFilePath, fmt.Sprintf("%s/%v.log", l.oldLogsDirPath, time.Now().Format("2006-01-02_15-04-05"))); err != nil {
		return fmt.Errorf("log file moving error: %w", err)
	}

	file, err := os.Create(l.mainFilePath)
	if err != nil {
		return fmt.Errorf("log file creation error: %w", err)
	}
	file.Close()

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		info, err := f.Info()
		if err != nil {
			return fmt.Errorf("old log file reading error: %w", err)
		}
		filteredFiles = append(filteredFiles, fileInfo{entry: f, modTime: info.ModTime()})
	}
	filesCount := len(filteredFiles) + 1

	if filesCount <= 3 {
		return nil
	}

	sort.Slice(filteredFiles, func(i, j int) bool {
		return filteredFiles[i].modTime.Before(filteredFiles[j].modTime)
	})

	toRemoveCount := filesCount - 3

	for i := 0; i < toRemoveCount; i++ {
		fileToDelete := filteredFiles[i].entry
		fullPath := filepath.Join(l.oldLogsDirPath, fileToDelete.Name())

		err := os.Remove(fullPath)
		if err != nil {
			return fmt.Errorf("old log file removing error: %w", err)
		}
	}

	return nil
}

func (l *Logger) LogInfo(s string) {
	l.logger.Info(s)
}

func (l *Logger) LogError(s string) {
	l.logger.Error(s)
}

func (l *Logger) LogFatal(s string) {
	l.logger.Error(fmt.Sprintf("FATAL: %s", s))
}

func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.LogInfo(fmt.Sprintf("Got %s request on localhost:8080%s from %s", r.Method, r.URL.Path, r.RemoteAddr))
		next.ServeHTTP(w, r)
	})
}
