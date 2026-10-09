<script lang="ts">
  interface Props {
    text: string;
    successMessage: string;
    class?: string;
  }

  let { text, successMessage, class: className = '' }: Props = $props();

  let copied = $state(false);

  function handleCopy() {
    if (typeof window === 'undefined') return;
    navigator.clipboard.writeText(window.location.href).then(() => {
      copied = true;
      alert(successMessage);
      setTimeout(() => {
        copied = false;
      }, 2000);
    }).catch(err => {
      console.error('Failed to copy link:', err);
    });
  }
</script>

<button
  type="button"
  onclick={handleCopy}
  class={className || "px-4 py-2 text-sm font-medium bg-blue-50 hover:bg-blue-100 border border-blue-200 text-blue-700 rounded-lg transition-colors cursor-pointer"}
>
  {text}
</button>
