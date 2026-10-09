<script lang="ts">
  import { onMount } from 'svelte';
  import { modalState } from '$state/modal.svelte';

  onMount(async () => {
    if (typeof window === 'undefined') return;

    try {
      const mermaid = (await import('mermaid')).default;
      mermaid.initialize({
        startOnLoad: false,
        theme: 'default',
        themeVariables: {
          fontFamily: 'Pretendard, -apple-system, BlinkMacSystemFont, system-ui, sans-serif',
          fontSize: '13px',
        }
      });

      const rawBlocks = document.querySelectorAll<HTMLElement>('.mermaid-raw');
      for (let i = 0; i < rawBlocks.length; i++) {
        const rawEl = rawBlocks[i];
        const container = rawEl.closest<HTMLElement>('.mermaid-container');
        if (!container) continue;

        // 이미 렌더링된 SVG가 있다면 건너뜀
        if (container.querySelector('.mermaid-svg-wrapper')) continue;

        const code = rawEl.textContent?.trim() || '';
        const id = 'mermaid-svg-' + i;
        try {
          const { svg } = await mermaid.render(id, code);
          const svgWrapper = document.createElement('div');
          svgWrapper.className = 'mermaid-svg-wrapper cursor-zoom-in w-full flex justify-center py-2';
          svgWrapper.innerHTML = svg;

          svgWrapper.addEventListener('click', () => {
            modalState.open(svg);
          });

          container.appendChild(svgWrapper);
        } catch (err: any) {
          console.error('Mermaid render error:', err);
          container.innerHTML = `<pre class="text-xs text-rose-600 p-4">${err?.message || 'Diagram render error'}</pre>`;
        }
      }
    } catch (e) {
      console.error('Failed to load mermaid:', e);
    }
  });
</script>
