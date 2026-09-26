package markdown

import (
	"strings"
	"testing"
)

func TestConverter_Convert(t *testing.T) {
	c := NewConverter()

	t.Run("Convert markdown headings to TOC and add IDs", func(t *testing.T) {
		input := []byte(`# Main Title

## Section 1: Introduction
Some introductory text.

### Subsection 1.1: Background
More text.

## Section 2: Architecture
Architecture text.
`)

		html, toc, err := c.Convert(input)
		if err != nil {
			t.Fatalf("Convert failed: %v", err)
		}

		// 1. TOC 구조체 검증
		if len(toc) != 3 {
			t.Fatalf("expected 3 TOC items, got %d", len(toc))
		}
		if toc[0].Title != "Section 1: Introduction" || toc[0].Level != 2 {
			t.Errorf("TOC[0] mismatch: %+v", toc[0])
		}
		if toc[1].Title != "Subsection 1.1: Background" || toc[1].Level != 3 {
			t.Errorf("TOC[1] mismatch: %+v", toc[1])
		}
		if toc[2].Title != "Section 2: Architecture" || toc[2].Level != 2 {
			t.Errorf("TOC[2] mismatch: %+v", toc[2])
		}

		// 2. HTML 내 id 속성 주입 검증
		htmlStr := string(html)
		if !strings.Contains(htmlStr, `id="section-1-introduction"`) {
			t.Errorf("HTML missing id='section-1-introduction': %s", htmlStr)
		}
		if !strings.Contains(htmlStr, `id="section-2-architecture"`) {
			t.Errorf("HTML missing id='section-2-architecture': %s", htmlStr)
		}
	})

	t.Run("Convert mermaid code block to container", func(t *testing.T) {
		input := []byte(`Here is a diagram:

` + "```mermaid" + `
flowchart LR
    A --> B
` + "```" + `

End of document.
`)

		html, _, err := c.Convert(input)
		if err != nil {
			t.Fatalf("Convert failed: %v", err)
		}

		htmlStr := string(html)
		if !strings.Contains(htmlStr, `class="mermaid-container`) {
			t.Errorf("expected mermaid-container in HTML, got:\n%s", htmlStr)
		}
		if !strings.Contains(htmlStr, `class="mermaid-raw"`) {
			t.Errorf("expected mermaid-raw in HTML, got:\n%s", htmlStr)
		}
		if !strings.Contains(htmlStr, `flowchart LR`) {
			t.Errorf("expected mermaid content in HTML, got:\n%s", htmlStr)
		}
	})

	t.Run("GFM Table and TaskList rendering", func(t *testing.T) {
		input := []byte(`
| Column 1 | Column 2 |
| :--- | :--- |
| Val 1 | Val 2 |

- [x] Done item
- [ ] Todo item
`)

		html, _, err := c.Convert(input)
		if err != nil {
			t.Fatalf("Convert failed: %v", err)
		}

		htmlStr := string(html)
		if !strings.Contains(htmlStr, "<table>") {
			t.Errorf("expected <table> in HTML, got:\n%s", htmlStr)
		}
		if !strings.Contains(htmlStr, `type="checkbox"`) {
			t.Errorf("expected checkbox in HTML, got:\n%s", htmlStr)
		}
	})
}
