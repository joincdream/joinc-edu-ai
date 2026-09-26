package cli

import (
	"fmt"
	"log"

	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/builder"
	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/deployer"
	"github.com/spf13/cobra"
)

var (
	deployBranch  string
	deployRemote  string
	deployCNAME   string
	deployMessage string
	noBuild       bool
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy compiled static site to GitHub Pages with orphan branch push",
	RunE: func(cmd *cobra.Command, args []string) error {
		distDir := "dist"

		// 1. 배포 전 선행 빌드 (프로덕션 모드: 드래프트 배제, Clean 활성화)
		if !noBuild {
			log.Println("[DEPLOY] Compiling production static site...")
			opts := builder.Options{
				SourceDir:     sourceDir,
				ThemeDir:      themeDir,
				OutputDir:     distDir,
				IncludeDrafts: false, // 프로덕션 배포에는 드래프트 절대 배제
				BaseURL:       "/",
				Clean:         true,
			}
			b := builder.NewBuilder(opts)
			res, err := b.Build()
			if err != nil {
				return fmt.Errorf("pre-deploy build failed: %w", err)
			}
			log.Printf("[DEPLOY] Production build finished in %v (%d posts compiled)", res.Duration, res.TotalPosts)
		}

		// 2. GitHub Pages 푸시
		deployOpts := deployer.DeployOptions{
			DistDir: distDir,
			Branch:  deployBranch,
			Remote:  deployRemote,
			CNAME:   deployCNAME,
			Message: deployMessage,
		}

		return deployer.Deploy(deployOpts)
	},
}

func init() {
	deployCmd.Flags().StringVarP(&deployBranch, "branch", "B", "gh-pages", "Target GitHub Pages branch")
	deployCmd.Flags().StringVarP(&deployRemote, "remote", "r", "origin", "Target Git remote")
	deployCmd.Flags().StringVar(&deployCNAME, "cname", "www.joinc.co.kr", "CNAME domain to inject into dist/")
	deployCmd.Flags().StringVarP(&deployMessage, "message", "m", "deploy: publish static site via site-cli", "Git commit message")
	deployCmd.Flags().BoolVar(&noBuild, "no-build", false, "Skip build step and deploy existing dist/ directory directly")
}
