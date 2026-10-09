<script lang="ts">
  import { onMount } from 'svelte';
  import type { MessageBundle } from '$types/messages';

  interface Props {
    currentLang: string;
    messages: MessageBundle;
    switchUrl?: {
      ko?: string;
      en?: string;
    };
  }

  let { currentLang, messages, switchUrl }: Props = $props();

  let visible = $state(false);

  onMount(() => {
    if (typeof window === 'undefined') return;

    const userLang = (navigator.language || (navigator as any).userLanguage || '').toLowerCase();
    const isEnBrowser = userLang.startsWith('en');
    const isDismissed = localStorage.getItem('i18n_banner_dismissed') === 'true';

    // 한국어 페이지 접속 중이고, 브라우저가 영어이며, 이전에 닫지 않은 경우 배너 노출
    if ((currentLang === 'ko' || currentLang === '') && isEnBrowser && !isDismissed) {
      visible = true;
    }
  });

  function handleDismiss() {
    visible = false;
    if (typeof window !== 'undefined') {
      localStorage.setItem('i18n_banner_dismissed', 'true');
    }
  }

  function handleSwitch() {
    if (typeof window !== 'undefined') {
      localStorage.setItem('preferred_lang', 'en');
    }
  }
</script>

{#if visible}
  <div id="i18n-banner" class="bg-blue-50 border-b border-blue-200 text-blue-900 px-4 py-2 text-xs sm:text-sm">
    <div class="max-w-7xl mx-auto flex items-center justify-between gap-4">
      <div class="flex items-center gap-2">
        <span>🌐</span>
        <span>{messages.banner?.detected_text || 'Detected English browser. Switch to English?'}</span>
      </div>
      <div class="flex items-center gap-2">
        {#if switchUrl?.en}
          <a
            href={switchUrl.en}
            onclick={handleSwitch}
            class="px-2.5 py-1 bg-blue-600 hover:bg-blue-700 text-white rounded font-medium text-xs transition-colors"
          >
            {messages.banner?.switch_btn || 'English'}
          </a>
        {/if}
        <button
          onclick={handleDismiss}
          class="text-blue-600 hover:text-blue-800 font-bold px-1.5 py-0.5 text-sm cursor-pointer"
          aria-label="Dismiss"
        >
          &times;
        </button>
      </div>
    </div>
  </div>
{/if}
