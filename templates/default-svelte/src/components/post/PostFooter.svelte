<script lang="ts">
  import type { MessageBundle } from '$types/messages';
  import CopyButton from '$components/ui/CopyButton.svelte';

  interface Props {
    baseUrl: string;
    messages: MessageBundle;
  }

  let { baseUrl, messages }: Props = $props();
</script>

<div class="mt-16 pt-8 border-t border-slate-200 flex items-center justify-between">
  <a
    href={baseUrl}
    class="px-4 py-2 text-sm font-medium bg-white hover:bg-slate-50 border border-slate-200 text-slate-700 rounded-lg shadow-sm transition-colors"
  >
    &larr; {messages.common?.back_to_list || 'Back to all posts'}
  </a>

  <CopyButton
    text={messages.common?.copy_link || 'Copy Link'}
    successMessage={messages.common?.copy_success || 'Link copied to clipboard!'}
  />
</div>
