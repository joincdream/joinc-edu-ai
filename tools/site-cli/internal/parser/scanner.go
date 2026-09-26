package parser

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/model"
)

// ScanPosts 지정된 루트 디렉터리를 재귀 탐색하여 모든 마크다운 포스트를 수집합니다.
// includeDrafts가 false인 경우 status: draft 포스트는 결과에서 제외됩니다.
func ScanPosts(rootDir string, includeDrafts bool) ([]*model.Post, error) {
	info, err := os.Stat(rootDir)
	if err != nil {
		return nil, fmt.Errorf("scanner: failed to access source directory %q: %w", rootDir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("scanner: source path %q is not a directory", rootDir)
	}

	var posts []*model.Post

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

		// .md 확장자 파일만 파싱 대상
		if strings.ToLower(filepath.Ext(path)) != ".md" {
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
