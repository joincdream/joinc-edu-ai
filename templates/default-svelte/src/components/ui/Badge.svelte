<script lang="ts">
  interface Props {
    variant?: 'primary' | 'draft' | 'tag' | 'count';
    href?: string;
    children?: any;
  }

  let { variant = 'primary', href = undefined, children }: Props = $props();

  const styles = {
    primary: 'inline-block px-2.5 py-0.5 text-xs font-semibold rounded-md bg-blue-50 text-blue-700 border border-blue-200 hover:bg-blue-100 transition-colors',
    draft: 'px-2 py-0.5 text-[11px] font-semibold rounded bg-amber-50 text-amber-700 border border-amber-200',
    tag: 'px-2 py-0.5 rounded bg-slate-100 text-slate-600 border border-slate-200 text-[11px]',
    count: 'px-2 py-0.5 text-xs rounded-full bg-slate-100 text-slate-600 font-mono'
  };
</script>

{#if href}
  <a {href} class={styles[variant]}>
    {#if children}
      {@render children()}
    {/if}
  </a>
{:else}
  <span class={styles[variant]}>
    {#if children}
      {@render children()}
    {/if}
  </span>
{/if}
