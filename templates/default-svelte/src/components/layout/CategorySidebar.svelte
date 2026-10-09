<script lang="ts">
  import type { CategoryExportData } from '$types/taxonomy';
  import type { MessageBundle } from '$types/messages';

  interface Props {
    baseUrl: string;
    categories: CategoryExportData[];
    activeCategorySlug?: string;
    totalPostCount: number;
    messages: MessageBundle;
  }

  let {
    baseUrl,
    categories,
    activeCategorySlug = '',
    totalPostCount,
    messages
  }: Props = $props();

  const isAllPostsActive = $derived(!activeCategorySlug);
</script>

<aside class="w-full lg:w-64 flex-shrink-0">
  <div class="sticky top-24">
    <div class="mb-6 flex items-center h-8">
      <h2 class="text-xs font-bold text-slate-500 uppercase tracking-wider">
        {messages.nav?.categories_title || 'Categories'}
      </h2>
    </div>
    <div class="bg-white border border-slate-200 rounded-xl p-4 shadow-sm">
      <nav class="space-y-1.5">
        <a
          href={baseUrl}
          class="flex items-center justify-between px-3 py-2 rounded-lg text-sm font-medium transition-colors {isAllPostsActive ? 'bg-blue-50 text-blue-700 border border-blue-200 font-semibold' : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'}"
        >
          <span>{messages.nav?.all_posts || 'All Posts'}</span>
          <span class="px-2 py-0.5 text-xs rounded-full bg-slate-100 text-slate-600 font-mono">
            {totalPostCount}
          </span>
        </a>

        {#each categories as cat (cat.slug)}
          {@const isActive = activeCategorySlug === cat.slug}
          <a
            href="{baseUrl}category/{cat.slug}/"
            class="flex items-center justify-between px-3 py-2 rounded-lg text-sm font-medium transition-colors {isActive ? 'bg-blue-50 text-blue-700 border border-blue-200 font-semibold' : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'}"
          >
            <span class="truncate">{cat.name}</span>
            <span class="px-2 py-0.5 text-xs rounded-full bg-slate-100 text-slate-600 font-mono">
              {cat.post_count}
            </span>
          </a>
        {/each}
      </nav>
    </div>
  </div>
</aside>
