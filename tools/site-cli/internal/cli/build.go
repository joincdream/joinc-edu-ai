package cli

import (
	"fmt"

	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/builder"
	"github.com/spf13/cobra"
)

var (
	outputDir     string
	includeDrafts bool
	cleanBuild    bool
	baseURL       string
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Compile Markdown posts and templates into static HTML/CSS/JS",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := builder.Options{
			SourceDir:     sourceDir,
			PagesDir:      pagesDir,
			ThemeDir:      themeDir,
			OutputDir:     outputDir,
			IncludeDrafts: includeDrafts,
			BaseURL:       baseURL,
			Clean:         cleanBuild,
		}

		b := builder.NewBuilder(opts)
		res, err := b.Build()
		if err != nil {
			return err
		}

		fmt.Printf("\n[SUCCESS] Static site compiled successfully in %v\n", res.Duration)
		fmt.Printf("├── Source:     %s\n", sourceDir)
		fmt.Printf("├── Theme:      %s\n", themeDir)
		fmt.Printf("├── Total Posts:%d (Drafts included: %v)\n", res.TotalPosts, includeDrafts)
		fmt.Printf("├── Categories: %d\n", res.TotalCategories)
		fmt.Printf("└── Output Dir: %s\n\n", res.OutputDir)

		return nil
	},
}

func init() {
	buildCmd.Flags().StringVarP(&outputDir, "output", "o", "dist", "Output directory for compiled static files")
	buildCmd.Flags().BoolVarP(&includeDrafts, "drafts", "D", false, "Include posts marked as status: draft")
	buildCmd.Flags().BoolVarP(&cleanBuild, "clean", "c", true, "Clean output directory before build")
	buildCmd.Flags().StringVar(&baseURL, "base-url", "/", "Base URL path prefix")
}
