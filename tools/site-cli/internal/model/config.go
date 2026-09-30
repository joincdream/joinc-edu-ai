package model

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// SiteConfig config.yaml 전역 사이트 설정 매핑 모델
type SiteConfig struct {
	Title         string `yaml:"title"`
	Subtitle      string `yaml:"subtitle"`
	BaseURL       string `yaml:"baseURL"`
	Theme         string `yaml:"theme"`
	Source        string `yaml:"source"`
	Pages         string `yaml:"pages"`
	Output        string `yaml:"output"`
	CNAME         string `yaml:"cname"`
	RedirectsFile string `yaml:"redirects"`
}

// DefaultSiteConfig 기본 사이트 설정값 반환
func DefaultSiteConfig() SiteConfig {
	return SiteConfig{
		Title:         "AI Info",
		Subtitle:      "Enterprise AI & Software Engineering Tech Blog",
		BaseURL:       "/",
		Theme:         "default-light",
		Source:        "posts",
		Pages:         "pages",
		Output:        "dist",
		CNAME:         "www.joinc.co.kr",
		RedirectsFile: "redirect.yaml",
	}
}

// LoadSiteConfig 지정된 경로의 config.yaml을 로드합니다. 파일이 없으면 기본값을 반환합니다.
func LoadSiteConfig(configPath string) (SiteConfig, error) {
	cfg := DefaultSiteConfig()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// ResolveThemeDir 테마 이름을 실제 파일시스템 디렉터리 경로로 해석합니다.
// 예: "default-light" -> "templates/default-light", "templates/default" -> "templates/default"
func ResolveThemeDir(themeName string) string {
	themeName = strings.TrimSpace(themeName)
	if themeName == "" {
		return "templates/default-light"
	}
	if strings.HasPrefix(themeName, "templates/") || strings.HasPrefix(themeName, "templates\\") {
		return filepath.Clean(themeName)
	}
	return filepath.Join("templates", themeName)
}
