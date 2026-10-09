<script lang="ts">
  import type { PostExportData } from '$types/post';
  import type { MessageBundle } from '$types/messages';
  import type { AlternateLink } from '$types/site';
  import BaseLayout from '$layouts/BaseLayout.svelte';
  import PostHeader from '$components/post/PostHeader.svelte';
  import PostFooter from '$components/post/PostFooter.svelte';
  import TableOfContents from '$components/post/TableOfContents.svelte';

  interface Props {
    post: PostExportData;
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
    post,
    currentLang,
    baseUrl,
    messages,
    switchUrl,
    alternateLangs = []
  }: Props = $props();

  const siteTitle = $derived(messages.common?.site_title || 'AI Info');
</script>

<BaseLayout
  title={post.frontmatter.title}
  {siteTitle}
  {currentLang}
  {baseUrl}
  currentPath="{baseUrl}posts/{post.slug}/"
  {messages}
  {switchUrl}
  {alternateLangs}
>
  <div class="flex flex-col lg:flex-row gap-10">
    <!-- Main Post Content -->
    <article class="flex-1 min-w-0">
      <!-- Breadcrumb & Back -->
      <div class="mb-6">
        <a
          href={baseUrl}
          class="text-xs font-medium text-blue-600 hover:text-blue-700 flex items-center gap-1 transition-colors"
        >
          &larr; {messages.common?.back_to_list || 'Back to all posts'}
        </a>
      </div>

      <!-- Post Header -->
      <PostHeader
        {post}
        {baseUrl}
        {messages}
      />

      <!-- Post Body -->
      <div class="prose max-w-none">
        {@html post.html_content}
      </div>

      <!-- Post Footer -->
      <PostFooter
        {baseUrl}
        {messages}
      />
    </article>

    <!-- Sidebar Table of Contents (TOC) -->
    <TableOfContents
      toc={post.toc}
      {messages}
    />
  </div>
</BaseLayout>
