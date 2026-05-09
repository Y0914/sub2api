//go:build !embed

// Package web provides frontend asset servers for the application.
package web

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// PublicSettingsProvider is an interface to fetch public settings
// This stub is needed for compilation when frontend is not embedded
type PublicSettingsProvider interface {
	GetPublicSettingsForInjection(ctx context.Context) (any, error)
}

// FrontendServer is a stub for non-embed builds
type FrontendServer struct{}

// NewFrontendServer returns an error when frontend is not embedded and no external frontend dir is provided.
func NewFrontendServer(settingsProvider PublicSettingsProvider, externalDir string) (*FrontendServer, error) {
	if strings.TrimSpace(externalDir) != "" {
		return NewExternalFrontendServer(settingsProvider, externalDir)
	}
	return nil, errors.New("frontend not embedded")
}

// NewExternalFrontendServer validates the external frontend directory for non-embed builds.
func NewExternalFrontendServer(settingsProvider PublicSettingsProvider, frontendDir string) (*FrontendServer, error) {
	if strings.TrimSpace(frontendDir) == "" {
		return nil, errors.New("external frontend dir is required")
	}
	info, err := os.Stat(filepath.Join(strings.TrimSpace(frontendDir), "index.html"))
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, errors.New("external frontend index.html is a directory")
	}
	return &FrontendServer{}, nil
}

// InvalidateCache is a no-op for non-embed builds
func (s *FrontendServer) InvalidateCache() {}

// Middleware returns a handler that returns 404 for non-embed builds
func (s *FrontendServer) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusNotFound, "Frontend server unavailable in !embed build. Build with -tags embed or use embedded release.")
		c.Abort()
	}
}

func ServeEmbeddedFrontend() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusNotFound, "Frontend not embedded. Build with -tags embed to include frontend.")
		c.Abort()
	}
}

func HasEmbeddedFrontend() bool {
	return false
}
