package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/model"
	"github.com/joincdream/joinc-ai.io/tools/site-cli/internal/parser"
	"github.com/spf13/cobra"
)

var (
	listDrafts bool
	outputJSON bool
)

var categoriesCmd = &cobra.Command{
	Use:     "list-categories",
	Aliases: []string{"categories"},
	Short:   "Inspect categories and post distribution across markdown files",
	RunE: func(cmd *cobra.Command, args []string) error {
		posts, err := parser.ScanPosts(sourceDir, listDrafts)
		if err != nil {
			return err
		}

		taxonomy := model.NewTaxonomyIndex(posts)

		if outputJSON {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(taxonomy.Categories)
		}

		fmt.Printf("\n📂 Total Posts Scanned: %d (Drafts included: %v)\n", len(taxonomy.AllPosts), listDrafts)
		fmt.Printf("%-25s %-25s %s\n", "CATEGORY NAME", "SLUG", "POSTS")
		fmt.Println("----------------------------------------------------------------------")
		for _, cat := range taxonomy.Categories {
			fmt.Printf("%-25s %-25s %d\n", cat.Name, cat.Slug, cat.PostCount)
		}
		fmt.Println()

		return nil
	},
}

func init() {
	categoriesCmd.Flags().BoolVarP(&listDrafts, "drafts", "D", false, "Include draft posts in counts")
	categoriesCmd.Flags().BoolVar(&outputJSON, "json", false, "Output results in JSON format")
}
