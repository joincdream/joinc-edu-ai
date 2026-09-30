package model

import (
	"errors"
	"fmt"
	"html/template"
	"regexp"
	"sort"
	"strings"
)

var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// TagList 문자열("a, b") 또는 배열(["a", "b"])을 모두 수용하는 유연한 태그 타입
type TagList []string

// UnmarshalYAML yaml.v2 및 yaml.v3 호환 언마샬링
func (t *TagList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var single string
	if err := unmarshal(&single); err == nil {
		parts := strings.Split(single, ",")
		var res []string
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				res = append(res, trimmed)
			}
		}
		*t = res
		return nil
	}

	var slice []string
	if err := unmarshal(&slice); err == nil {
		var res []string
		for _, s := range slice {
			trimmed := strings.TrimSpace(s)
			if trimmed != "" {
				res = append(res, trimmed)
			}
		}
		*t = res
		return nil
	}

	return nil
}

// Frontmatter 마크다운 상단 YAML 메타데이터
type Frontmatter struct {
	Title         string   `yaml:"title"`
	Category      string   `yaml:"category"`
	Tags          TagList  `yaml:"tags"`
	CreatedDate   string   `yaml:"created_date"`
	PublishedDate string   `yaml:"published_date"`
	Summary       string   `yaml:"summary"`
	Description   string   `yaml:"description"`
	Thumbnail     string   `yaml:"thumbnail"`
	Status        string   `yaml:"status"` // "published" 또는 "draft"
	Draft         bool     `yaml:"draft"`  // draft: true 플래그 지원
	PostID        int      `yaml:"post_id,omitempty"`
}

// Validate Frontmatter 필수 필드 및 날짜 포맷 검증
func (f *Frontmatter) Validate() error {
	if strings.TrimSpace(f.Title) == "" {
		return errors.New("frontmatter: 'title' is required")
	}
	if strings.TrimSpace(f.Category) == "" {
		f.Category = "General"
	}
	if strings.TrimSpace(f.CreatedDate) == "" {
		return errors.New("frontmatter: 'created_date' is required")
	}
	if !dateRegex.MatchString(f.CreatedDate) {
		return fmt.Errorf("frontmatter: invalid 'created_date' format: %q (must be YYYY-MM-DD)", f.CreatedDate)
	}
	if f.Status == "" {
		f.Status = "published"
	}
	if strings.TrimSpace(f.Summary) == "" && strings.TrimSpace(f.Description) != "" {
		f.Summary = f.Description
	} else if strings.TrimSpace(f.Description) == "" && strings.TrimSpace(f.Summary) != "" {
		f.Description = f.Summary
	}
	return nil
}

// IsDraft 드래프트 포스트 여부 확인
func (f *Frontmatter) IsDraft() bool {
	if f.Draft {
		return true
	}
	return strings.ToLower(strings.TrimSpace(f.Status)) == "draft"
}

// TOCItem 마크다운 본문에서 추출된 목차 노드 (H2, H3)
type TOCItem struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Level    int       `json:"level"` // 2 (H2) or 3 (H3)
	Children []TOCItem `json:"children,omitempty"`
}

// Post 단일 마크다운 포스트 모델
type Post struct {
	ID                 string        `json:"id"`
	Slug               string        `json:"slug"`
	FilePath           string        `json:"file_path"`
	Frontmatter        Frontmatter   `json:"frontmatter"`
	RawContent         string        `json:"raw_content"`
	HTMLContent        template.HTML `json:"html_content"`
	TOC                []TOCItem     `json:"toc"`
	ReadingTimeMinutes int           `json:"reading_time_minutes"`
}

// Category 카테고리 색인 모델
type Category struct {
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	PostCount int     `json:"post_count"`
	Posts     []*Post `json:"posts,omitempty"`
}

// AddPost 카테고리에 포스트 추가 및 카운트 증가
func (c *Category) AddPost(p *Post) {
	c.Posts = append(c.Posts, p)
	c.PostCount = len(c.Posts)
}

// TaxonomyIndex 전체 카테고리/태그 색인 결과
type TaxonomyIndex struct {
	AllPosts    []*Post              `json:"all_posts"`
	Categories  []*Category          `json:"categories"`
	CategoryMap map[string]*Category `json:"-"`
}

// NewTaxonomyIndex 포스트 목록으로부터 최신순 정렬 및 카테고리 역색인 생성
func NewTaxonomyIndex(posts []*Post) *TaxonomyIndex {
	sortedPosts := make([]*Post, len(posts))
	copy(sortedPosts, posts)
	sort.Slice(sortedPosts, func(i, j int) bool {
		return sortedPosts[i].Frontmatter.CreatedDate > sortedPosts[j].Frontmatter.CreatedDate
	})

	catMap := make(map[string]*Category)
	var catList []*Category

	for _, p := range sortedPosts {
		catName := strings.TrimSpace(p.Frontmatter.Category)
		if catName == "" {
			catName = "General"
		}
		catSlug := Slugify(catName)

		cat, exists := catMap[catSlug]
		if !exists {
			cat = &Category{
				Name: catName,
				Slug: catSlug,
			}
			catMap[catSlug] = cat
			catList = append(catList, cat)
		}
		cat.AddPost(p)
	}

	sort.Slice(catList, func(i, j int) bool {
		return strings.ToLower(catList[i].Name) < strings.ToLower(catList[j].Name)
	})

	return &TaxonomyIndex{
		AllPosts:    sortedPosts,
		Categories:  catList,
		CategoryMap: catMap,
	}
}

// Slugify 문자열을 URL-Safe 슬러그로 변환
func Slugify(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	reg := regexp.MustCompile(`[^a-z0-9가-힣]+`)
	s = reg.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "default"
	}
	return s
}

// MessageBundle messages.yaml 매핑 모델
type MessageBundle struct {
	Common map[string]string `yaml:"common"`
	Nav    map[string]string `yaml:"nav"`
	Card   map[string]string `yaml:"card"`
	Detail map[string]string `yaml:"detail"`
	Empty  map[string]string `yaml:"empty"`
	Footer map[string]string `yaml:"footer"`
}

// Get helper: 키 경로 또는 맵에서 안전하게 문자열 반환 (누락 시 fallback 반환)
func (m *MessageBundle) Get(section, key, fallback string) string {
	var targetMap map[string]string
	switch strings.ToLower(section) {
	case "common":
		targetMap = m.Common
	case "nav":
		targetMap = m.Nav
	case "card":
		targetMap = m.Card
	case "detail":
		targetMap = m.Detail
	case "empty":
		targetMap = m.Empty
	case "footer":
		targetMap = m.Footer
	}

	if targetMap != nil {
		if val, ok := targetMap[key]; ok && strings.TrimSpace(val) != "" {
			return val
		}
	}
	return fallback
}

// TemplateContext html/template 렌더링 시 주입되는 컨텍스트
type TemplateContext struct {
	SiteTitle      string
	SiteSubtitle   string
	BaseURL        string
	CurrentPath    string
	Messages       MessageBundle
	Design         *DesignTokens
	Categories     []*Category
	ActiveCategory *Category
	Posts          []*Post
	Post           *Post
	TOC            []TOCItem
	ExtraHead      template.HTML
}
