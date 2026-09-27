package builder

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/model"
	"gopkg.in/yaml.v3"
)

// loadRedirectConfig 지정된 YAML 파일에서 리다이렉트 설정을 로드합니다.
func loadRedirectConfig(path string) (*model.RedirectConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg model.RedirectConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse redirect yaml %s: %w", path, err)
	}

	return &cfg, nil
}

// generateRedirects redirect.yaml 설정을 읽어 정적 리다이렉트 페이지들을 dist에 생성합니다.
func (b *Builder) generateRedirects() error {
	if b.opts.RedirectsFile == "" {
		return nil
	}

	// 파일이 존재하지 않는 경우 건너뜀 (Optional)
	if _, err := os.Stat(b.opts.RedirectsFile); os.IsNotExist(err) {
		return nil
	}

	cfg, err := loadRedirectConfig(b.opts.RedirectsFile)
	if err != nil {
		return fmt.Errorf("builder: %w", err)
	}

	// 1. query_redirects 처리
	for _, qr := range cfg.QueryRedirects {
		if qr.SourcePath == "" {
			continue
		}
		if qr.QueryKey == "" {
			qr.QueryKey = "id"
		}
		if qr.DefaultTarget == "" {
			qr.DefaultTarget = "/"
		}
		if qr.Mappings == nil {
			qr.Mappings = make(map[string]string)
		}

		htmlContent, err := buildQueryRedirectHTML(qr)
		if err != nil {
			return fmt.Errorf("builder: failed to generate query redirect HTML for %s: %w", qr.SourcePath, err)
		}

		relPath := strings.TrimPrefix(strings.TrimSuffix(qr.SourcePath, "/"), "/")
		targetDir := filepath.Join(b.opts.OutputDir, relPath)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("builder: failed to create redirect dir %s: %w", targetDir, err)
		}

		destPath := filepath.Join(targetDir, "index.html")
		if err := os.WriteFile(destPath, []byte(htmlContent), 0644); err != nil {
			return fmt.Errorf("builder: failed to write redirect file %s: %w", destPath, err)
		}
		log.Printf("[REDIRECT] Generated query redirect: %s (key=%s, targets=%d)", qr.SourcePath, qr.QueryKey, len(qr.Mappings))
	}

	// 2. path_redirects 처리
	for _, pr := range cfg.PathRedirects {
		if pr.SourcePath == "" || pr.TargetPath == "" {
			continue
		}

		htmlContent, err := buildPathRedirectHTML(pr)
		if err != nil {
			return fmt.Errorf("builder: failed to generate path redirect HTML for %s: %w", pr.SourcePath, err)
		}

		relPath := strings.TrimPrefix(strings.TrimSuffix(pr.SourcePath, "/"), "/")
		targetDir := filepath.Join(b.opts.OutputDir, relPath)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("builder: failed to create path redirect dir %s: %w", targetDir, err)
		}

		destPath := filepath.Join(targetDir, "index.html")
		if err := os.WriteFile(destPath, []byte(htmlContent), 0644); err != nil {
			return fmt.Errorf("builder: failed to write path redirect file %s: %w", destPath, err)
		}
		log.Printf("[REDIRECT] Generated path redirect: %s -> %s", pr.SourcePath, pr.TargetPath)
	}

	return nil
}

// buildQueryRedirectHTML 쿼리 파라미터 기반 리다이렉트 HTML 템플릿 생성
func buildQueryRedirectHTML(qr model.QueryRedirect) (string, error) {
	mappingsJSON, err := json.Marshal(qr.Mappings)
	if err != nil {
		return "", err
	}
	queryKeyJSON, err := json.Marshal(qr.QueryKey)
	if err != nil {
		return "", err
	}
	defaultTargetJSON, err := json.Marshal(qr.DefaultTarget)
	if err != nil {
		return "", err
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="ko" class="dark">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>새로운 페이지로 이동 중 - joinc-ai.io</title>
  <noscript>
    <meta http-equiv="refresh" content="0; url=%s">
  </noscript>
  <style>
    body {
      margin: 0;
      padding: 0;
      background-color: #020617;
      color: #f8fafc;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
      display: flex;
      align-items: center;
      justify-content: center;
      min-height: 100vh;
    }
    .container {
      text-align: center;
      padding: 2rem;
      max-width: 480px;
    }
    .spinner {
      width: 44px;
      height: 44px;
      border: 3px solid rgba(56, 189, 248, 0.2);
      border-top-color: #38bdf8;
      border-radius: 50%%;
      animation: spin 0.8s linear infinite;
      margin: 0 auto 1.5rem auto;
    }
    @keyframes spin {
      to { transform: rotate(360deg); }
    }
    h1 {
      font-size: 1.25rem;
      font-weight: 600;
      margin-bottom: 0.5rem;
      color: #f1f5f9;
    }
    p {
      font-size: 0.95rem;
      color: #94a3b8;
      margin-bottom: 1.5rem;
      line-height: 1.5;
    }
    a {
      color: #38bdf8;
      text-decoration: underline;
      text-underline-offset: 4px;
      font-size: 0.875rem;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="spinner"></div>
    <h1>새로운 페이지로 안전하게 이동 중입니다...</h1>
    <p>기존 URL 체계 변경에 따라 해당 글의 새 주소로 자동 이동합니다.</p>
    <noscript>
      <p>자바스크립트가 비활성화되어 있습니다. <a href="%s">여기를 클릭하여 이동하세요</a>.</p>
    </noscript>
    <div id="fallback" style="display:none;">
      <a id="targetLink" href="%s">자동으로 이동하지 않으면 여기를 클릭하세요.</a>
    </div>
  </div>
  <script>
    (function() {
      const params = new URLSearchParams(window.location.search);
      const queryKey = %s;
      const targetId = params.get(queryKey);
      const mappings = %s;
      const defaultTarget = %s;
      const destination = (targetId && mappings[targetId]) ? mappings[targetId] : defaultTarget;

      const link = document.getElementById("targetLink");
      if (link) link.href = destination;
      const fallback = document.getElementById("fallback");
      if (fallback) fallback.style.display = "block";

      window.location.replace(destination);
    })();
  </script>
</body>
</html>`,
		qr.DefaultTarget,
		qr.DefaultTarget,
		qr.DefaultTarget,
		string(queryKeyJSON),
		string(mappingsJSON),
		string(defaultTargetJSON),
	)

	return html, nil
}

// buildPathRedirectHTML 단순 경로 1:1 리다이렉트 HTML 생성
func buildPathRedirectHTML(pr model.PathRedirect) (string, error) {
	targetJSON, err := json.Marshal(pr.TargetPath)
	if err != nil {
		return "", err
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="ko" class="dark">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="refresh" content="0; url=%s">
  <title>새로운 페이지로 이동 중 - joinc-ai.io</title>
  <style>
    body {
      margin: 0;
      padding: 0;
      background-color: #020617;
      color: #f8fafc;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
      display: flex;
      align-items: center;
      justify-content: center;
      min-height: 100vh;
    }
    .container {
      text-align: center;
      padding: 2rem;
      max-width: 480px;
    }
    .spinner {
      width: 44px;
      height: 44px;
      border: 3px solid rgba(56, 189, 248, 0.2);
      border-top-color: #38bdf8;
      border-radius: 50%%;
      animation: spin 0.8s linear infinite;
      margin: 0 auto 1.5rem auto;
    }
    @keyframes spin {
      to { transform: rotate(360deg); }
    }
    h1 {
      font-size: 1.25rem;
      font-weight: 600;
      margin-bottom: 0.5rem;
      color: #f1f5f9;
    }
    p {
      font-size: 0.95rem;
      color: #94a3b8;
      margin-bottom: 1.5rem;
      line-height: 1.5;
    }
    a {
      color: #38bdf8;
      text-decoration: underline;
      text-underline-offset: 4px;
      font-size: 0.875rem;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="spinner"></div>
    <h1>새로운 페이지로 안전하게 이동 중입니다...</h1>
    <p>해당 주소로 자동 리다이렉트합니다.</p>
    <a href="%s">자동으로 이동하지 않으면 여기를 클릭하세요.</a>
  </div>
  <script>
    window.location.replace(%s);
  </script>
</body>
</html>`,
		pr.TargetPath,
		pr.TargetPath,
		string(targetJSON),
	)

	return html, nil
}
