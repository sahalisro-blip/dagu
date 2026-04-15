// Copyright (C) 2026 Yota Hamada
// SPDX-License-Identifier: GPL-3.0-or-later

package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

var imagePathPattern = regexp.MustCompile(`(?i)(/[^\s"'<>|]+?\.(?:jpg|jpeg|png|gif))`)

type imagesResponse struct {
	Images []string `json:"images"`
}

func (a *API) handleGetImages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	logContent := strings.Join(query["log"], "\n")

	if logContent == "" {
		if logFile := strings.TrimSpace(query.Get("logFile")); logFile != "" {
			data, err := os.ReadFile(logFile)
			if err == nil {
				logContent = string(data)
			}
		}
	}

	sourcePrefix := query.Get("sourcePrefix")
	if sourcePrefix == "" {
		sourcePrefix = "/data1"
	}
	servePrefix := query.Get("servePrefix")
	if servePrefix == "" {
		servePrefix = "/images"
	}

	images := extractImagePaths(logContent, sourcePrefix, servePrefix)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(imagesResponse{Images: images})
}

func extractImagePaths(logContent, sourcePrefix, servePrefix string) []string {
	matches := imagePathPattern.FindAllString(logContent, -1)
	if len(matches) == 0 {
		return []string{}
	}

	seen := make(map[string]struct{}, len(matches))
	images := make([]string, 0, len(matches))
	for _, raw := range matches {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}
		mapped := mapImagePath(trimmed, sourcePrefix, servePrefix)
		if _, ok := seen[mapped]; ok {
			continue
		}
		seen[mapped] = struct{}{}
		images = append(images, mapped)
	}
	return images
}

func mapImagePath(rawPath, sourcePrefix, servePrefix string) string {
	cleaned := path.Clean(strings.TrimSpace(rawPath))
	if !strings.HasPrefix(cleaned, "/") {
		cleaned = "/" + cleaned
	}
	src := normalizePathPrefix(sourcePrefix)
	dst := normalizePathPrefix(servePrefix)

	if src != "/" && (cleaned == src || strings.HasPrefix(cleaned, src+"/")) {
		rel := strings.TrimPrefix(cleaned, src)
		return path.Clean(dst + "/" + strings.TrimPrefix(rel, "/"))
	}
	return cleaned
}

func normalizePathPrefix(prefix string) string {
	if strings.TrimSpace(prefix) == "" {
		return "/"
	}
	cleaned := path.Clean("/" + strings.TrimSpace(prefix))
	if cleaned == "." {
		return "/"
	}
	return cleaned
}
