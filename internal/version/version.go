package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/TBXark/gbvm/internal/env"
)

func FetchLatest(modName string) (string, error) {
	url := fmt.Sprintf("%s/%s/@latest", env.GoProxy, strings.ToLower(modName))
	resp, err := http.Get(url)
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
	parts1 := strings.Split(trimVersion(v1), ".")
	parts2 := strings.Split(trimVersion(v2), ".")
	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}
	for i := 0; i < maxLen; i++ {
		num1 := 0
		if i < len(parts1) {
			num1, _ = strconv.Atoi(parts1[i])
		}
		num2 := 0
		if i < len(parts2) {
			num2, _ = strconv.Atoi(parts2[i])
		}
		if num1 > num2 {
			return 1
		}
		if num1 < num2 {
			return -1
		}
	}
	return 0
}

func trimVersion(version string) string {
	if version == env.DevelVersion {
		return "0"
	}
	ver := strings.Split(strings.TrimPrefix(version, "v"), "-")
	if len(ver) > 1 && ver[0] == "0.0.0" {
		if ts, err := strconv.ParseInt(ver[1], 10, 64); err == nil {
			return fmt.Sprintf("0.0.0.%d", ts)
		}
	}
	return ver[0]
}
