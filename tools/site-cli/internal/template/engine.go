package template

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/model"
	"gopkg.in/yaml.v3"
)

// Engine 로컬 파일시스템 기반 템플릿 컴파일러
type Engine struct {
	themeDir    string
	messages    model.MessageBundle
	messagesMap map[string]model.MessageBundle
	design      *model.DesignTokens
	funcMap     template.FuncMap
}

// NewEngine 주어진 테마 디렉터리로부터 템플릿 엔진을 초기화합니다.
// 바이너리 임베딩 없이 로컬 파일시스템의 파일들을 직접 읽습니다.
func NewEngine(themeDir string) (*Engine, error) {
	info, err := os.Stat(themeDir)
	if err != nil {
		return nil, fmt.Errorf("template: theme directory %q not found: %w", themeDir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("template: theme path %q is not a directory", themeDir)
	}

	messagesMap := make(map[string]model.MessageBundle)

	// 1. 기본 messages.yaml 로드 (한국어/기본)
	msgPath := filepath.Join(themeDir, "messages.yaml")
	msgBytes, err := os.ReadFile(msgPath)
	if err != nil {
		return nil, fmt.Errorf("template: failed to read messages bundle %q: %w", msgPath, err)
	}

	var messages model.MessageBundle
	if err := yaml.Unmarshal(msgBytes, &messages); err != nil {
		return nil, fmt.Errorf("template: failed to parse messages YAML %q: %w", msgPath, err)
	}
	messagesMap["ko"] = messages

	// 테마 디렉터리 내 추가 언어 번들 파일(messages_en.yaml 등) 자동 로드
	entries, _ := os.ReadDir(themeDir)
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasPrefix(name, "messages_") || !strings.HasSuffix(name, ".yaml") {
			continue
		}
		langCode := strings.TrimSuffix(strings.TrimPrefix(name, "messages_"), ".yaml")
		langCode = strings.ToLower(strings.TrimSpace(langCode))
		if langCode == "" {
			continue
		}

		subBytes, subErr := os.ReadFile(filepath.Join(themeDir, name))
		if subErr == nil {
			var subBundle model.MessageBundle
			if yamlErr := yaml.Unmarshal(subBytes, &subBundle); yamlErr == nil {
				messagesMap[langCode] = subBundle
			}
		}
	}

	// 2. DESIGN.md 로드 (선택적)
	design, err := loadDesign(themeDir)
	if err != nil {
		return nil, err
	}

	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"join": func(elems []string, sep string) string {
			return strings.Join(elems, sep)
		},
		"lower": strings.ToLower,
		"dict": func(values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict requires an even number of arguments")
			}
			dict := make(map[string]any, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	}

	return &Engine{
		themeDir:    themeDir,
		messages:    messages,
		messagesMap: messagesMap,
		design:      design,
		funcMap:     funcMap,
	}, nil
}

// GetMessages 로드된 기본 메시지 리소스 번들을 반환합니다.
func (e *Engine) GetMessages() model.MessageBundle {
	return e.messages
}

// GetMessagesFor 지정된 언어 코드의 메시지 번들을 반환합니다 (없으면 기본 메시지 번들 fallback).
func (e *Engine) GetMessagesFor(lang string) model.MessageBundle {
	lang = strings.ToLower(strings.TrimSpace(lang))
	if bundle, exists := e.messagesMap[lang]; exists {
		return bundle
	}
	return e.messages
}

// GetDesign 로드된 디자인 시스템 토큰을 반환합니다 (없으면 nil).
func (e *Engine) GetDesign() *model.DesignTokens {
	return e.design
}

// loadDesign 테마 디렉터리에서 DESIGN.md를 로드하여 파싱합니다. 없으면 nil을 반환합니다.
func loadDesign(themeDir string) (*model.DesignTokens, error) {
	designPath := filepath.Join(themeDir, "DESIGN.md")
	content, err := os.ReadFile(designPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("template: failed to read DESIGN.md %q: %w", designPath, err)
	}

	tokens, err := parseDesignTokens(content)
	if err != nil {
		return nil, fmt.Errorf("template: failed to parse %q: %w", designPath, err)
	}
	return tokens, nil
}

// parseDesignTokens DESIGN.md 파일의 Frontmatter를 파싱하여 DesignTokens 구조체로 반환합니다.
func parseDesignTokens(content []byte) (*model.DesignTokens, error) {
	str := strings.TrimSpace(string(content))
	if !strings.HasPrefix(str, "---") {
		return nil, nil
	}

	rest := str[3:]
	if idx := strings.IndexAny(rest, "\r\n"); idx != -1 {
		rest = rest[idx:]
	}

	endIdx := strings.Index(rest, "\n---")
	if endIdx == -1 {
		return nil, fmt.Errorf("invalid DESIGN.md format: closing delimiter '---' not found")
	}

	frontmatter := rest[:endIdx]
	var tokens model.DesignTokens
	if err := yaml.Unmarshal([]byte(frontmatter), &tokens); err != nil {
		return nil, fmt.Errorf("failed to unmarshal DESIGN.md YAML: %w", err)
	}
	return &tokens, nil
}

// HasTemplate 특정 템플릿 파일이 테마 디렉터리에 존재하는지 확인합니다.
func (e *Engine) HasTemplate(pageName string) bool {
	pageFilePath := filepath.Join(e.themeDir, pageName)
	_, err := os.Stat(pageFilePath)
	return err == nil
}

// loadTemplate 지정된 페이지 템플릿과 기본 레이아웃을 조합하여 컴파일합니다.
func (e *Engine) loadTemplate(pageName string) (*template.Template, error) {
	pageFilePath := filepath.Join(e.themeDir, pageName)
	if _, err := os.Stat(pageFilePath); err != nil {
		return nil, fmt.Errorf("template: page file not found: %s", pageFilePath)
	}

	// 1. base.html 및 components 파일들 수집
	var filesToParse []string
	basePath := filepath.Join(e.themeDir, "base.html")
	if _, err := os.Stat(basePath); err == nil {
		filesToParse = append(filesToParse, basePath)
	}

	// components 하위 파일 수집
	componentsDir := filepath.Join(e.themeDir, "components")
	if info, err := os.Stat(componentsDir); err == nil && info.IsDir() {
		_ = filepath.WalkDir(componentsDir, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr == nil && !d.IsDir() && strings.ToLower(filepath.Ext(path)) == ".html" {
				filesToParse = append(filesToParse, path)
			}
			return nil
		})
	}

	// 2. 대상 페이지 파일 추가
	if filepath.Clean(pageFilePath) != filepath.Clean(basePath) {
		filesToParse = append(filesToParse, pageFilePath)
	}

	t := template.New("base").Funcs(e.funcMap)
	t, err := t.ParseFiles(filesToParse...)
	if err != nil {
		return nil, fmt.Errorf("template: failed to parse HTML files for %s: %w", pageName, err)
	}

	return t, nil
}

// RenderPage 주어진 템플릿 이름과 컨텍스트로 HTML을 렌더링합니다.
func (e *Engine) RenderPage(w io.Writer, templateName string, ctx *model.TemplateContext) error {
	if ctx.CurrentLang != "" {
		ctx.Messages = e.GetMessagesFor(ctx.CurrentLang)
	} else if ctx.Messages.Common == nil {
		ctx.Messages = e.messages
	}
	if ctx.Design == nil {
		ctx.Design = e.design
	}

	t, err := e.loadTemplate(templateName)
	if err != nil {
		return err
	}

	// base 템플릿을 실행하여 자식 페이지의 block "content"를 렌더링
	if err := t.ExecuteTemplate(w, "base", ctx); err != nil {
		return fmt.Errorf("template: failed to execute %s: %w", templateName, err)
	}

	return nil
}

// RenderCustomTemplate 테마 외부의 커스텀 템플릿 파일(예: pages/about.html)과 base 레이아웃을 조합하여 렌더링합니다.
func (e *Engine) RenderCustomTemplate(w io.Writer, templatePath string, ctx *model.TemplateContext) error {
	if ctx.CurrentLang != "" {
		ctx.Messages = e.GetMessagesFor(ctx.CurrentLang)
	} else if ctx.Messages.Common == nil {
		ctx.Messages = e.messages
	}
	if ctx.Design == nil {
		ctx.Design = e.design
	}

	if _, err := os.Stat(templatePath); err != nil {
		return fmt.Errorf("template: custom template file not found: %s: %w", templatePath, err)
	}

	var filesToParse []string
	basePath := filepath.Join(e.themeDir, "base.html")
	if _, err := os.Stat(basePath); err == nil {
		filesToParse = append(filesToParse, basePath)
	}

	componentsDir := filepath.Join(e.themeDir, "components")
	if info, err := os.Stat(componentsDir); err == nil && info.IsDir() {
		_ = filepath.WalkDir(componentsDir, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr == nil && !d.IsDir() && strings.ToLower(filepath.Ext(path)) == ".html" {
				filesToParse = append(filesToParse, path)
			}
			return nil
		})
	}

	filesToParse = append(filesToParse, templatePath)

	t := template.New("base").Funcs(e.funcMap)
	t, err := t.ParseFiles(filesToParse...)
	if err != nil {
		return fmt.Errorf("template: failed to parse custom template %s: %w", templatePath, err)
	}

	if err := t.ExecuteTemplate(w, "base", ctx); err != nil {
		return fmt.Errorf("template: failed to execute custom template %s: %w", templatePath, err)
	}

	return nil
}

