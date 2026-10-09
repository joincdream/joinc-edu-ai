<script lang="ts">
  import type { MessageBundle } from '$types/messages';

  interface Props {
    messages: MessageBundle;
  }

  let { messages }: Props = $props();
</script>

<footer class="border-t border-slate-200 bg-white py-8 mt-auto">
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex flex-col sm:flex-row items-center justify-between gap-4 text-xs text-slate-500">
    <p>{messages.footer?.copyright || '© 2026 Joinc AI. All rights reserved.'}</p>
    <p class="text-slate-400">{messages.footer?.powered_by || 'Powered by Pure Static Joinc AI Generator'}</p>
  </div>
</footer>
