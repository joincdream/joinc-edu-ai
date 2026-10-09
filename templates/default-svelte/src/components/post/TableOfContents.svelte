<script lang="ts">
  import type { TOCItem } from '$types/post';
  import type { MessageBundle } from '$types/messages';

  interface Props {
    toc: TOCItem[];
    messages: MessageBundle;
  }

  let { toc, messages }: Props = $props();
</script>

{#if toc && toc.length > 0}
  <aside class="w-full lg:w-72 flex-shrink-0">
    <div class="sticky top-24 space-y-4">
      <div class="bg-white border border-slate-200 rounded-xl p-5 shadow-sm">
        <h2 class="text-sm font-semibold text-slate-500 uppercase tracking-wider mb-4">
          {messages.common?.table_of_contents || 'Table of Contents'}
        </h2>
        <nav class="space-y-2 text-sm max-h-[calc(100vh-14rem)] overflow-y-auto pr-1">
          {#each toc as item (item.id)}
            <a
              href="#{item.id}"
              class="block text-slate-600 hover:text-blue-600 transition-colors {item.level === 3 ? 'pl-4 text-xs' : item.level >= 4 ? 'pl-8 text-xs' : ''}"
            >
              {item.title}
            </a>
          {/each}
        </nav>
      </div>

      <div class="p-4 bg-slate-50 border border-slate-200 rounded-xl text-xs text-slate-500 flex items-center gap-2">
        <svg class="w-4 h-4 text-blue-600 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        <span>{messages.detail?.mermaid_zoom_tip || 'Click on any Mermaid diagram to expand and zoom.'}</span>
      </div>
    </div>
  </aside>
{/if}
