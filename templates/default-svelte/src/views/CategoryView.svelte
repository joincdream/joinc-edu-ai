<script lang="ts">
  import type { PostExportData } from '$types/post';
  import type { CategoryExportData } from '$types/taxonomy';
  import type { MessageBundle } from '$types/messages';
  import type { AlternateLink } from '$types/site';
  import BaseLayout from '$layouts/BaseLayout.svelte';
  import CategorySidebar from '$components/layout/CategorySidebar.svelte';
  import PostCard from '$components/post/PostCard.svelte';

  interface Props {
    category: CategoryExportData;
    posts: PostExportData[];
    allPostsCount: number;
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
    category,
    posts,
    allPostsCount,
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
  title={category.name}
  {siteTitle}
  {currentLang}
  {baseUrl}
  currentPath="{baseUrl}category/{category.slug}/"
  {messages}
  {switchUrl}
  {alternateLangs}
>
  <div class="flex flex-col lg:flex-row gap-8">
    <!-- Sidebar Category Navigation -->
    <CategorySidebar
      {baseUrl}
      {categories}
      activeCategorySlug={category.slug}
      totalPostCount={allPostsCount}
      {messages}
    />

    <!-- Post Grid Section -->
    <main class="flex-1">
      <div class="mb-6 flex items-center justify-between h-8">
        <div class="flex items-center gap-3">
          <a
            href={baseUrl}
            class="text-xs text-blue-600 hover:text-blue-700 font-medium inline-flex items-center gap-1"
          >
            &larr; {messages.common?.back_to_list || 'Back to all posts'}
          </a>
          <span class="text-slate-300">|</span>
          <h1 class="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
            <span>{category.name}</span>
            <span class="text-xs px-2 py-0.5 rounded-full bg-blue-50 text-blue-700 border border-blue-200 font-mono font-semibold">
              {posts.length}{messages.nav?.post_count_unit || ' posts'}
            </span>
          </h1>
        </div>
      </div>

      {#if posts.length === 0}
        <div class="bg-white border border-slate-200 rounded-xl p-12 text-center text-slate-500 shadow-sm">
          {messages.empty?.no_category_posts || 'No posts found in this category.'}
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
