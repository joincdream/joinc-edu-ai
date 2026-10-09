<script lang="ts">
  import type { PostExportData } from '$types/post';
  import type { CategoryExportData } from '$types/taxonomy';
  import type { MessageBundle } from '$types/messages';
  import type { AlternateLink } from '$types/site';
  import BaseLayout from '$layouts/BaseLayout.svelte';
  import CategorySidebar from '$components/layout/CategorySidebar.svelte';
  import PostCard from '$components/post/PostCard.svelte';

  interface Props {
    posts: PostExportData[];
    categories: CategoryExportData[];
    currentLang: string;
    baseUrl: string;
    messages: MessageBundle;
    switchUrl?: {
      ko?: string;
      en?: string;
    };
    alternateLangs?: AlternateLink[];
  }

  let {
    posts,
    categories,
    currentLang,
    baseUrl,
    messages,
    switchUrl,
    alternateLangs = []
  }: Props = $props();

  const siteTitle = $derived(messages.common?.site_title || 'AI Info');
</script>

<BaseLayout
  title=""
  {siteTitle}
  {currentLang}
  {baseUrl}
  currentPath={baseUrl}
  {messages}
  {switchUrl}
  {alternateLangs}
>
  <div class="flex flex-col lg:flex-row gap-8">
    <!-- Sidebar Category Navigation -->
    <CategorySidebar
      {baseUrl}
      {categories}
      totalPostCount={posts.length}
      {messages}
    />

    <!-- Post Grid Section -->
    <main class="flex-1">
      <div class="mb-6 flex items-center justify-between h-8">
        <h1 class="text-2xl font-bold text-slate-900 tracking-tight">
          {messages.nav?.all_posts || 'All Posts'}
        </h1>
        <span class="text-xs text-slate-500 font-medium">
          {messages.nav?.total_count_prefix || 'Total '}{posts.length}{messages.nav?.post_count_unit || ' posts'}
        </span>
      </div>

      {#if posts.length === 0}
        <div class="bg-white border border-slate-200 rounded-xl p-12 text-center text-slate-500 shadow-sm">
          {messages.empty?.no_posts_found || 'No posts available.'}
        </div>
      {:else}
        <div class="flex flex-col gap-5">
          {#each posts as post (post.slug)}
            <PostCard
              {post}
              {baseUrl}
              {messages}
            />
          {/each}
        </div>
      {/if}
    </main>
  </div>
</BaseLayout>
