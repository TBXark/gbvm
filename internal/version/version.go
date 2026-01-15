package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/TBXark/gbvm/internal/env"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

func FetchLatest(modName string) (string, string, error) {
	proxies, err := proxyCandidates(env.GoProxies)
	if err != nil {
		return "", "", err
	}
	escapedPath, err := module.EscapePath(modName)
	if err != nil {
		return "", "", err
	}
	var latestVersion string
	var selectedProxy string
	var lastErr error
	for _, proxy := range proxies {
		version, err := fetchLatestFromProxy(proxy, escapedPath)
		if err != nil {
			lastErr = err
			continue
		}
		if latestVersion == "" || Compare(latestVersion, version) < 0 {
			latestVersion = version
			selectedProxy = proxy
		}
	}
	if latestVersion == "" {
		if lastErr != nil {
			return "", "", lastErr
		}
		return "", "", fmt.Errorf("no valid GOPROXY endpoint")
	}
	return latestVersion, selectedProxy, nil
}

func fetchLatestFromProxy(proxy, escapedPath string) (string, error) {
	url := fmt.Sprintf("%s/%s/@latest", proxy, escapedPath)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var versionInfo struct {
		Version string `json:"Version"`
	}
	err = json.Unmarshal(body, &versionInfo)
	if err != nil {
		return "", err
	}
	return versionInfo.Version, nil
}

func Compare(v1, v2 string) int {
	return semver.Compare(normalizeVersion(v1), normalizeVersion(v2))
}

func normalizeVersion(version string) string {
	if version == env.DevelVersion {
		return "v0.0.0"
	}
	normalized := version
	if !strings.HasPrefix(normalized, "v") {
		normalized = "v" + normalized
	}
	if semver.IsValid(normalized) {
		return normalized
	}
	trimmed := strings.TrimPrefix(version, "v")
	parts := strings.Split(trimmed, "-")
	fallback := "v" + parts[0]
	if semver.IsValid(fallback) {
		return fallback
	}
	return "v0.0.0"
}

func proxyCandidates(proxies []string) ([]string, error) {
	var candidates []string
	seen := make(map[string]struct{})
	for _, part := range proxies {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		if trimmed == "direct" {
			if _, exists := seen[env.OfficialProxy]; !exists {
				candidates = append(candidates, env.OfficialProxy)
				seen[env.OfficialProxy] = struct{}{}
			}
			continue
		}
		if trimmed == "off" {
			return nil, fmt.Errorf("GOPROXY is off")
		}
		if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
			cleaned := strings.TrimRight(trimmed, "/")
			if _, exists := seen[cleaned]; !exists {
				candidates = append(candidates, cleaned)
				seen[cleaned] = struct{}{}
			}
		}
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no valid GOPROXY endpoint")
	}
	return candidates, nil
}
