package cli

import (
	"fmt"
	"log"
	"time"

	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/builder"
	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/server"
	"github.com/spf13/cobra"
)

var (
	servePort      int
	serveBind      string
	serveDrafts    bool
	noWatch        bool
	serveRedirects string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start local development HTTP server with live reloading (Hugo-like)",
	RunE: func(cmd *cobra.Command, args []string) error {
		distDir := "dist"

		hub := server.NewSSEHub()

		// 빌드 헬퍼 함수
		doBuild := func() error {
			opts := builder.Options{
				SourceDir:     sourceDir,
				PagesDir:      pagesDir,
				ThemeDir:      themeDir,
				OutputDir:     distDir,
				IncludeDrafts: serveDrafts,
				BaseURL:       "/",
				Clean:         false, // 빠른 증분 서빙을 위해 clean 배제
				ExtraHead:     server.LiveReloadScript,
				RedirectsFile: serveRedirects,
			}
			b := builder.NewBuilder(opts)
			res, err := b.Build()
			if err != nil {
				return err
			}
			log.Printf("[BUILD] Recompiled in %v (Posts: %d, Categories: %d)", res.Duration, res.TotalPosts, res.TotalCategories)
			return nil
		}

		// 1. 초기 1회 빌드 실행
		log.Println("[SERVE] Running initial compilation...")
		if err := doBuild(); err != nil {
			return fmt.Errorf("initial build failed: %w", err)
		}

		// 2. 파일 감시자 시작 (트리플 감시: posts/, templates/, pages/)
		if !noWatch {
			watchDirs := []string{sourceDir, themeDir, pagesDir}
			watcher, err := server.NewWatcher(watchDirs, 200*time.Millisecond, func() {
				if err := doBuild(); err != nil {
					log.Printf("[BUILD ERROR] %v", err)
					return
				}
				// 브라우저에 새로고침 신호 브로드캐스트
				hub.Broadcast("reload")
			})
			if err != nil {
				return fmt.Errorf("failed to start watcher: %w", err)
			}
			if err := watcher.Start(); err != nil {
				return err
			}
			defer watcher.Stop()
			log.Printf("[WATCHER] Watching for changes in %v", watchDirs)
		}

		// 3. 로컬 HTTP 서버 구동
		addr := fmt.Sprintf("%s:%d", serveBind, servePort)
		fmt.Printf("\n🚀 Local test server running at: http://%s\n", addr)
		fmt.Println("   Press Ctrl+C to stop.")

		return server.StartLocalServer(addr, distDir, hub)
	},
}

func init() {
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8080, "Port to bind local server to")
	serveCmd.Flags().StringVarP(&serveBind, "bind", "b", "127.0.0.1", "Host address to bind to")
	serveCmd.Flags().BoolVarP(&serveDrafts, "drafts", "D", true, "Include draft posts in local preview")
	serveCmd.Flags().BoolVar(&noWatch, "no-watch", false, "Disable file watching and Live Reload")
	serveCmd.Flags().StringVar(&serveRedirects, "redirects", "redirect.yaml", "Path to redirect configuration YAML file")
}
