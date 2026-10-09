<script lang="ts">
  import type { MessageBundle } from '$types/messages';

  interface Props {
    baseUrl: string;
    currentPath: string;
    currentLang: string;
    messages: MessageBundle;
    switchUrl?: {
      ko?: string;
      en?: string;
    };
  }

  let { baseUrl, currentPath, currentLang, messages, switchUrl }: Props = $props();

  const isAboutActive = $derived(
    currentPath === `${baseUrl}about/` ||
    currentPath === '/about/' ||
    currentPath === '/en/about/'
  );
</script>

<header class="sticky top-0 z-40 w-full border-b border-slate-200 bg-white/80 backdrop-blur-md">
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
    <div class="flex items-center gap-4">
      <a href={baseUrl} class="flex items-center gap-2 group">
        <div class="w-9 h-9 rounded-lg bg-blue-600 flex items-center justify-center text-white font-bold text-lg shadow-md shadow-blue-500/20 group-hover:scale-105 transition-transform">
          AI
        </div>
        <div class="flex flex-col">
          <span class="text-lg font-bold tracking-tight text-slate-900 group-hover:text-blue-600 transition-colors">
            {messages.common?.site_title || 'AI Info'}
          </span>
          <span class="text-xs text-slate-500 hidden sm:inline">
            {messages.common?.site_subtitle || ''}
          </span>
        </div>
      </a>
    </div>

    <nav class="flex items-center gap-3">
      <a
        href="{baseUrl}about/"
        class="px-3.5 py-1.5 text-sm font-medium transition-colors {isAboutActive ? 'bg-blue-50 text-blue-700 border border-blue-200 rounded-lg font-semibold' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100 rounded-lg'}"
      >
        {messages.nav?.about || 'About'}
      </a>

      {#if switchUrl}
        <!-- Language Switcher [ KO | EN ] -->
        <div class="flex items-center bg-slate-100 p-0.5 rounded-lg text-xs font-semibold border border-slate-200">
          <a
            href={switchUrl.ko || '/'}
            class="px-2.5 py-1 rounded-md transition-all {currentLang === 'ko' ? 'bg-white text-blue-600 shadow-xs font-bold' : 'text-slate-500 hover:text-slate-900'}"
          >
            KO
          </a>
          <a
            href={switchUrl.en || '/en/'}
            class="px-2.5 py-1 rounded-md transition-all {currentLang === 'en' ? 'bg-white text-blue-600 shadow-xs font-bold' : 'text-slate-500 hover:text-slate-900'}"
          >
            EN
          </a>
        </div>
      {/if}
    </nav>
  </div>
</header>
