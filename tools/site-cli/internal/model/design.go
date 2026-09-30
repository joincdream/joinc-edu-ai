package model

// DesignTokens Google Labs DESIGN.md 포맷의 머신 리더블 YAML 토큰 구조체
type DesignTokens struct {
	Version     string                      `yaml:"version,omitempty"`
	Name        string                      `yaml:"name"`
	Description string                      `yaml:"description,omitempty"`
	Colors      map[string]string           `yaml:"colors,omitempty"`
	Typography  map[string]TypographyToken  `yaml:"typography,omitempty"`
	Rounded     map[string]string           `yaml:"rounded,omitempty"`
	Spacing     map[string]string           `yaml:"spacing,omitempty"`
	Components  map[string]map[string]any   `yaml:"components,omitempty"`
}

// TypographyToken 개별 폰트 스타일 토큰 (스케일별 폰트 패밀리, 크기, 굵기, 줄간격)
type TypographyToken struct {
	FontFamily    string `yaml:"fontFamily,omitempty"`
	FontSize      string `yaml:"fontSize,omitempty"`
	FontWeight    any    `yaml:"fontWeight,omitempty"`
	LineHeight    any    `yaml:"lineHeight,omitempty"`
	LetterSpacing string `yaml:"letterSpacing,omitempty"`
}
