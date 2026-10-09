package parser

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/model"
)

// ScanPosts 지정된 루트 디렉터리를 재귀 탐색하여 기본(한국어) 마크다운 포스트를 수집합니다.
// includeDrafts가 false인 경우 status: draft 포스트는 결과에서 제외됩니다.
func ScanPosts(rootDir string, includeDrafts bool) ([]*model.Post, error) {
	return ScanPostsByLang(rootDir, includeDrafts, "ko")
}

// ScanPostsByLang 언어 코드("ko" 또는 "en")에 맞춰 해당 언어의 마크다운 포스트만 수집합니다.
// - "en": *.en.md 접미사 파일만 수집
// - "ko" 또는 기타: *.en.md를 제외한 순수 *.md 파일 수집
func ScanPostsByLang(rootDir string, includeDrafts bool, lang string) ([]*model.Post, error) {
	info, err := os.Stat(rootDir)
	if err != nil {
		return nil, fmt.Errorf("scanner: failed to access source directory %q: %w", rootDir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("scanner: source path %q is not a directory", rootDir)
	}

	var posts []*model.Post
	isEnglish := strings.ToLower(lang) == "en"

	err = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// 무시할 디렉터리 필터링
		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "dist" || name == "templates" {
				return filepath.SkipDir
			}
			return nil
		}

		lowerPath := strings.ToLower(path)

		// .md 확장자 파일만 파싱 대상
		if filepath.Ext(lowerPath) != ".md" {
			return nil
		}

		// 언어별 파일 필터링: 영문은 *.en.md 만, 한국어는 *.en.md 제외
		hasEnSuffix := strings.HasSuffix(lowerPath, ".en.md")
		if isEnglish && !hasEnSuffix {
			return nil
		}
		if !isEnglish && hasEnSuffix {
			return nil
		}

		post, parseErr := ParseFile(path)
		if parseErr != nil {
			// [SE-04 & P-05] Fail-Fast: 메타데이터 깨짐 발생 시 즉시 중단 및 원인 통보
			return fmt.Errorf("scanner error: %w", parseErr)
		}
		if post == nil {
			// Frontmatter가 없는 관리용 문서는 수집 제외
			return nil
		}

		// 드래프트 필터링
		if !includeDrafts && post.Frontmatter.IsDraft() {
			return nil
		}

		posts = append(posts, post)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return posts, nil
}
