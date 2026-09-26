package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// 글로벌 플래그 변수
	sourceDir string
	pagesDir  string
	themeDir  string
	verbose   bool
)

// RootCmd 메인 루트 커맨드
var RootCmd = &cobra.Command{
	Use:   "site-cli",
	Short: "Pure Static Site Generator and GitHub Pages GitOps deployer for AI Info",
	Long: `site-cli is a high-performance Go-based Static Site Generator (SSG) 
designed to compile Markdown posts with decoupled external templates and deploy them to GitHub Pages.`,
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&sourceDir, "source", "s", "posts", "Markdown content source directory")
	RootCmd.PersistentFlags().StringVarP(&pagesDir, "pages", "p", "pages", "Standalone pages directory (e.g. pages/about.html)")
	RootCmd.PersistentFlags().StringVarP(&themeDir, "theme", "t", "templates/default", "Template theme directory (filesystem path)")
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
