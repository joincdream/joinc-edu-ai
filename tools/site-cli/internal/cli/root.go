package cli

import (
	"fmt"
	"os"

	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/model"
	"github.com/spf13/cobra"
)

var (
	// 글로벌 플래그 변수
	configPath string
	sourceDir  string
	pagesDir   string
	themeDir   string
	verbose    bool

	// 로드된 전역 설정
	loadedConfig model.SiteConfig
)

// RootCmd 메인 루트 커맨드
var RootCmd = &cobra.Command{
	Use:   "site-cli",
	Short: "Pure Static Site Generator and GitHub Pages GitOps deployer for AI Info",
	Long: `site-cli is a high-performance Go-based Static Site Generator (SSG) 
designed to compile Markdown posts with decoupled external templates and deploy them to GitHub Pages.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := model.LoadSiteConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config file %q: %w", configPath, err)
		}
		loadedConfig = cfg

		// CLI 플래그가 직접 입력되지 않은 경우 config.yaml 설정값을 우선 적용
		if !cmd.Flags().Changed("theme") && cfg.Theme != "" {
			themeDir = model.ResolveThemeDir(cfg.Theme)
		}
		if !cmd.Flags().Changed("source") && cfg.Source != "" {
			sourceDir = cfg.Source
		}
		if !cmd.Flags().Changed("pages") && cfg.Pages != "" {
			pagesDir = cfg.Pages
		}

		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().StringVar(&configPath, "config", "config.yaml", "Path to site configuration file")
	RootCmd.PersistentFlags().StringVarP(&sourceDir, "source", "s", "posts", "Markdown content source directory")
	RootCmd.PersistentFlags().StringVar(&pagesDir, "pages", "pages", "Standalone pages directory (e.g. pages/about.html)")
	RootCmd.PersistentFlags().StringVarP(&themeDir, "theme", "t", "templates/default-light", "Template theme directory (filesystem path)")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose debug output")

	// 서브커맨드 등록
	RootCmd.AddCommand(buildCmd)
	RootCmd.AddCommand(serveCmd)
	RootCmd.AddCommand(deployCmd)
	RootCmd.AddCommand(categoriesCmd)
}

// Execute CLI 엔트리포인트 실행 함수
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
