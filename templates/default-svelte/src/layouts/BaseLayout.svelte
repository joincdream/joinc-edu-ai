<script lang="ts">
  import type { MessageBundle } from '$types/messages';
  import type { AlternateLink } from '$types/site';
  import Header from '$components/layout/Header.svelte';
  import Footer from '$components/layout/Footer.svelte';
  import LanguageBanner from '$components/layout/LanguageBanner.svelte';
  import Modal from '$components/ui/Modal.svelte';

  interface Props {
    title: string;
    siteTitle: string;
    currentLang: string;
    baseUrl: string;
    currentPath: string;
    messages: MessageBundle;
    switchUrl?: {
      ko?: string;
      en?: string;
    };
    alternateLangs?: AlternateLink[];
    children?: any;
  }

  let {
    title,
    siteTitle,
    currentLang,
    baseUrl,
    currentPath,
    messages,
    switchUrl,
    alternateLangs = [],
    children
  }: Props = $props();

  const fullTitle = $derived(
    title ? `${title} | ${siteTitle}` : siteTitle
  );
</script>

<svelte:head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>{fullTitle}</title>

  {#each alternateLangs as alt}
    <link rel="alternate" hreflang={alt.lang} href={alt.url} />
  {/each}

  <!-- Pretendard Font -->
  <link rel="stylesheet" as="style" crossorigin href="https://cdn.jsdelivr.net/gh/orioncactus/pretendard@v1.3.9/dist/web/static/pretendard.min.css" />

  <!-- PrismJS Light Syntax Highlighting -->
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism.min.css" />

  <!-- KaTeX CSS -->
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.css" />

  <!-- Tailwind CSS CDN -->
  <script src="https://cdn.tailwindcss.com"></script>
  <script>
    tailwind.config = {
      darkMode: 'class',
      theme: {
        extend: {
          fontFamily: {
            sans: ['Pretendard', '-apple-system', 'BlinkMacSystemFont', 'system-ui', 'Roboto', 'sans-serif'],
          },
          colors: {
            brand: {
              50: '#eff6ff',
              500: '#3b82f6',
              600: '#2563eb',
              700: '#1d4ed8',
            }
          }
        }
      }
    }
  </script>

  <!-- Mermaid.js CDN -->
  <script src="https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.min.js"></script>

  <!-- KaTeX Math Rendering CDN -->
  <script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.js"></script>
  <script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/contrib/auto-render.min.js"></script>

  <!-- PrismJS Scripts -->
  <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/prism.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/plugins/autoloader/prism-autoloader.min.js"></script>

  <!-- Custom Theme Styles -->
  <link rel="stylesheet" href="/assets/style.css" />
</svelte:head>

<div class="bg-slate-50 text-slate-900 min-h-screen flex flex-col font-sans antialiased selection:bg-blue-600 selection:text-white">
  <!-- Browser Language Smart Recommendation Banner -->
  <LanguageBanner
    {currentLang}
    {messages}
    {switchUrl}
  />

  <!-- Navigation Header -->
  <Header
    {baseUrl}
    {currentPath}
    {currentLang}
    {messages}
    {switchUrl}
  />

  <!-- Main Content Area -->
  <div class="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-8">
    {#if children}
      {@render children()}
    {/if}
  </div>

  <!-- Footer -->
  <Footer {messages} />

  <!-- Mermaid Click-to-Zoom Modal -->
  <Modal
    title={messages.detail?.modal_title || '다이어그램 확대 보기'}
    closeLabel={messages.detail?.modal_close || '닫기'}
  />

  <!-- Client-side Interactive Scripts -->
  {@html `<script>
    // 1. Mermaid 초기화 및 렌더링
    if (typeof mermaid !== 'undefined') {
      mermaid.initialize({
        startOnLoad: false,
        theme: 'default',
        themeVariables: {
          fontFamily: 'Pretendard, -apple-system, BlinkMacSystemFont, system-ui, sans-serif',
          fontSize: '13px',
        }
      });
    }

    async function initMermaid() {
      if (typeof mermaid === 'undefined') return;
      var rawBlocks = document.querySelectorAll(".mermaid-raw");
      for (var i = 0; i < rawBlocks.length; i++) {
        var rawEl = rawBlocks[i];
        var container = rawEl.closest(".mermaid-container");
        if (!container) continue;
        if (container.querySelector(".mermaid-svg-wrapper")) continue;

        var code = rawEl.textContent.trim();
        var id = "mermaid-svg-" + i;
        try {
          var res = await mermaid.render(id, code);
          var svg = res.svg;
          var svgWrapper = document.createElement("div");
          svgWrapper.className = "mermaid-svg-wrapper cursor-zoom-in w-full flex justify-center py-2";
          svgWrapper.innerHTML = svg;
          
          (function(currentSvg) {
            svgWrapper.addEventListener("click", function() {
              openMermaidModal(currentSvg);
            });
          })(svg);

          container.appendChild(svgWrapper);
        } catch (err) {
          console.error("Mermaid render error:", err);
          container.innerHTML = '<pre class="text-xs text-rose-600 p-4">' + err.message + '</pre>';
        }
      }
    }

    // 2. 줌 모달 로직
    var modal = document.getElementById("mermaid-modal");
    var modalBody = document.getElementById("modal-body");
    var closeBtn = document.getElementById("modal-close-btn");

    function openMermaidModal(svgHtml) {
      if (!modal || !modalBody) return;
      modalBody.innerHTML = svgHtml;
      var svg = modalBody.querySelector("svg");
      if (svg) {
        svg.removeAttribute("width");
        svg.removeAttribute("height");
        svg.style.width = "100%";
        svg.style.maxWidth = "100%";
        svg.style.maxHeight = "80vh";
        svg.style.height = "auto";
        svg.style.display = "block";
        svg.style.margin = "auto";
      }
      modal.classList.remove("hidden");
      document.body.classList.add("overflow-hidden");
    }

    function closeMermaidModal() {
      if (!modal || !modalBody) return;
      modal.classList.add("hidden");
      modalBody.innerHTML = "";
      document.body.classList.remove("overflow-hidden");
    }

    if (closeBtn) closeBtn.addEventListener("click", closeMermaidModal);
    if (modal) {
      modal.addEventListener("click", function(e) {
        if (e.target === modal) closeMermaidModal();
      });
    }
    document.addEventListener("keydown", function(e) {
      if (e.key === "Escape" && modal && !modal.classList.contains("hidden")) {
        closeMermaidModal();
      }
    });

    // 3. KaTeX 수식 자동 렌더링
    function initKaTeX() {
      if (typeof renderMathInElement === "function") {
        renderMathInElement(document.body, {
          delimiters: [
            { left: "$$", right: "$$", display: true }
          ],
          ignoredTags: ["script", "noscript", "style", "textarea", "pre", "code"],
          ignoredClasses: ["mermaid-raw", "mermaid-container"],
          throwOnError: false
        });
      }
    }

    if (document.readyState === "loading") {
      document.addEventListener("DOMContentLoaded", function() {
        initMermaid();
        initKaTeX();
      });
    } else {
      initMermaid();
      initKaTeX();
    }
  </script>`}
</div>
