<script lang="ts">
  interface Props {
    title?: string;
    closeLabel?: string;
  }

  let { title = '다이어그램 확대 보기', closeLabel = '닫기' }: Props = $props();
</script>

<!-- Mermaid Click-to-Zoom Modal -->
<div
  id="mermaid-modal"
  class="fixed inset-0 z-50 hidden bg-slate-900/50 backdrop-blur-sm p-2 sm:p-4 md:p-6 flex items-center justify-center"
  role="dialog"
  aria-modal="true"
  tabindex="-1"
>
  <div class="relative w-full h-full max-w-[98vw] max-h-[96vh] bg-white border border-slate-200 rounded-2xl p-4 sm:p-6 shadow-2xl flex flex-col">
    <div class="flex items-center justify-between pb-3 sm:pb-4 border-b border-slate-200">
      <div class="flex items-center gap-2">
        <span class="inline-block w-2.5 h-2.5 rounded-full bg-blue-600 animate-pulse"></span>
        <span class="text-sm font-semibold text-slate-800">{title}</span>
      </div>
      <button
        id="modal-close-btn"
        class="px-3.5 py-1.5 text-xs sm:text-sm bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-lg transition-colors border border-slate-200 font-medium flex items-center gap-1.5 cursor-pointer"
      >
        <span>{closeLabel}</span>
        <kbd class="px-1.5 py-0.5 text-[10px] bg-white border border-slate-300 rounded text-slate-500">ESC</kbd>
      </button>
    </div>
    <div id="modal-body" class="flex-1 overflow-auto flex items-center justify-center p-2 sm:p-6 bg-slate-50 rounded-xl">
      <!-- Expanded SVG will be injected here -->
    </div>
  </div>
</div>
