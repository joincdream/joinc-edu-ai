<script lang="ts">
  import type { PostExportData } from '$types/post';
  import type { MessageBundle } from '$types/messages';
  import Badge from '$components/ui/Badge.svelte';
  import { slugify } from '../../utils/slugify';

  interface Props {
    post: PostExportData;
    baseUrl: string;
    messages: MessageBundle;
  }

  let { post, baseUrl, messages }: Props = $props();

  const categorySlug = $derived(
    post.frontmatter.category ? slugify(post.frontmatter.category) : 'general'
  );
</script>

<header class="mb-10 pb-8 border-b border-slate-200">
  <div class="flex items-center gap-2 mb-4">
    <Badge variant="primary" href="{baseUrl}category/{categorySlug}/">
      {post.frontmatter.category || 'General'}
    </Badge>
    <span class="text-xs text-slate-400 font-mono">
      {post.reading_time_minutes} {messages.common?.read_time_suffix || 'min read'}
    </span>
  </div>

  <h1 class="text-3xl sm:text-4xl font-extrabold text-slate-900 tracking-tight leading-tight mb-6">
    {post.frontmatter.title}
  </h1>

  <div class="flex flex-wrap items-center justify-between gap-4 text-xs text-slate-500">
    <div class="flex items-center gap-4">
      <span>
        {messages.detail?.author_label || 'Author:'}
        <strong class="text-slate-800">{messages.detail?.author_name || 'Sangbae Yoon'}</strong>
      </span>
      <span>•</span>
      <span>
        {messages.card?.published_prefix || 'Published:'} {post.frontmatter.created_date}
      </span>
    </div>

    {#if post.frontmatter.tags && post.frontmatter.tags.length > 0}
      <div class="flex items-center gap-1.5 flex-wrap">
        {#each post.frontmatter.tags as tag}
          <Badge variant="tag">
            #{tag}
          </Badge>
        {/each}
      </div>
    {/if}
  </div>
</header>
