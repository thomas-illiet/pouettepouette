package editor

import (
	"common/log"
	"fmt"
	"io"
	"net/http"
	"strings"
	"supervisor/pkg/config"
	"time"
)

// runEditorReadinessProbe ensures the configured editor is available before marking it ready.
func runEditorReadinessProbe(cfg *config.Config) {
	defer log.WithField("ide", cfg.Editor.Name).Info("editor is ready")

	switch cfg.Editor.ReadinessProbe.Type {
	case config.ReadinessProcessProbe:
		// No readiness check needed
		return

	case config.ReadinessHTTPProbe:
		url := buildProbeURL(cfg)
		for range time.Tick(250 * time.Millisecond) {
			body, err := editorStatusRequest(url)
			if err != nil {
				log.WithError(err).Debug("editor readiness probe failed")
				continue
			}
			log.WithField("body", string(body)).Debug("editor readiness response received")
			break
		}
	}
}

// buildProbeURL constructs the full readiness probe URL based on the editor's configuration.
func buildProbeURL(cfg *config.Config) string {
	schema := defaultIfEmpty(cfg.Editor.ReadinessProbe.HTTPProbe.Schema, "http")
	host := defaultIfEmpty(cfg.Editor.ReadinessProbe.HTTPProbe.Host, "localhost")
	port := defaultIfZero(cfg.Editor.ReadinessProbe.HTTPProbe.Port, 3000)
	path := strings.TrimPrefix(cfg.Editor.ReadinessProbe.HTTPProbe.Path, "/")

	return fmt.Sprintf("%s://%s:%d/%s", schema, host, port, path)
}

// editorStatusRequest sends an HTTP GET request to the given readiness probe URL.
func editorStatusRequest(url string) ([]byte, error) {
	client := http.Client{Timeout: 1 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %v", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
