<script lang="ts">
  import { onMount } from 'svelte';

  onMount(async () => {
    if (typeof window === 'undefined') return;

    try {
      // @ts-ignore
      const renderMathInElement = (await import('katex/contrib/auto-render')).default;
      renderMathInElement(document.body, {
        delimiters: [
          { left: "$$", right: "$$", display: true },
          { left: "$", right: "$", display: false },
          { left: "\\(", right: "\\)", display: false },
          { left: "\\[", right: "\\]", display: true }
        ],
        ignoredTags: ["script", "noscript", "style", "textarea", "pre", "code"],
        ignoredClasses: ["mermaid-raw", "mermaid-container"],
        throwOnError: false
      });
    } catch (e) {
      console.error('Failed to auto-render KaTeX:', e);
    }
  });
</script>
