package cli

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/builder"
	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/server"
	"github.com/spf13/cobra"
)

var (
	servePort      int
	serveBind      string
	serveDir       string
	serveWatch     bool
	serveDrafts    bool
	serveRedirects string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve compiled static site (dist/) or start live reload dev server with --watch",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. 디렉터리 존재 여부 확인
		info, err := os.Stat(serveDir)
		if err != nil || !info.IsDir() {
			if !serveWatch {
				return fmt.Errorf("target directory %q not found. Please run 'make build' first", serveDir)
			}
		}

		var hub *server.SSEHub

		// 2. --watch 플래그가 활성화된 경우에만 Live Reload 및 증분 빌드 활성화
		if serveWatch {
			hub = server.NewSSEHub()

			doBuild := func() error {
				opts := builder.Options{
					SourceDir:     sourceDir,
					PagesDir:      pagesDir,
					ThemeDir:      themeDir,
					OutputDir:     serveDir,
					IncludeDrafts: serveDrafts,
					BaseURL:       "/",
					Clean:         false,
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

			log.Println("[WATCH] Running initial compilation for watch mode...")
			if err := doBuild(); err != nil {
				return fmt.Errorf("initial build failed: %w", err)
			}

			watchDirs := []string{sourceDir, themeDir, pagesDir}
			watcher, err := server.NewWatcher(watchDirs, 200*time.Millisecond, func() {
				if err := doBuild(); err != nil {
					log.Printf("[BUILD ERROR] %v", err)
					return
				}
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

		// 3. 로컬 HTTP 서버 구동 (순수 정적 서빙)
		addr := fmt.Sprintf("%s:%d", serveBind, servePort)
		fmt.Printf("\n🚀 Serving %q at: http://%s\n", serveDir, addr)
		fmt.Println("   Press Ctrl+C to stop.")

		return server.StartLocalServer(addr, serveDir, hub)
	},
}

func init() {
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8080, "Port to bind local server to")
	serveCmd.Flags().StringVarP(&serveBind, "bind", "b", "127.0.0.1", "Host address to bind to")
	serveCmd.Flags().StringVarP(&serveDir, "dir", "d", "dist", "Directory to serve static files from")
	serveCmd.Flags().BoolVarP(&serveWatch, "watch", "w", false, "Enable file watcher and live reload compiler")
	serveCmd.Flags().BoolVarP(&serveDrafts, "drafts", "D", true, "Include draft posts in watch mode")
	serveCmd.Flags().StringVar(&serveRedirects, "redirects", "redirect.yaml", "Path to redirect configuration YAML file")
}
