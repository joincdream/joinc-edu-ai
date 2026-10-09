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

  const displayExcerpt = $derived(
    post.frontmatter.summary || post.frontmatter.description || post.excerpt
  );
</script>

<article class="bg-white border border-slate-200 hover:border-blue-400 rounded-xl overflow-hidden shadow-sm hover:shadow-md transition-all hover:-translate-y-0.5 group">
  <div class="p-6 flex flex-col justify-between">
    <div>
      <div class="flex items-center justify-between gap-2 mb-3">
        <div class="flex items-center gap-2">
          <Badge variant="primary" href="{baseUrl}category/{categorySlug}/">
            {post.frontmatter.category || 'General'}
          </Badge>
          {#if post.frontmatter.is_draft || post.frontmatter.status === 'draft'}
            <Badge variant="draft">
              DRAFT
            </Badge>
          {/if}
        </div>
        <span class="text-xs text-slate-400 font-mono">
          {post.reading_time_minutes} {messages.common?.read_time_suffix || 'min read'}
        </span>
      </div>

      <h2 class="text-xl font-bold text-slate-900 group-hover:text-blue-600 transition-colors mb-2 leading-snug">
        <a href="{baseUrl}posts/{post.slug}/">
          {post.frontmatter.title}
        </a>
      </h2>

      <p class="text-sm text-slate-600 line-clamp-2 mb-4 leading-relaxed">
        {displayExcerpt}
      </p>
    </div>

    <div class="pt-4 border-t border-slate-100 flex items-center justify-between">
      <div class="flex items-center gap-1.5 flex-wrap">
        {#each post.frontmatter.tags || [] as tag}
          <Badge variant="tag">
            #{tag}
          </Badge>
        {/each}
      </div>

      <div class="flex items-center gap-3 text-xs text-slate-400">
        <span>{post.frontmatter.created_date}</span>
        <span>&bull;</span>
        <a href="{baseUrl}posts/{post.slug}/" class="text-blue-600 hover:text-blue-700 font-medium inline-flex items-center gap-1 group-hover:translate-x-0.5 transition-transform">
          {messages.card?.read_more || 'Read More'} &rarr;
        </a>
      </div>
    </div>
  </div>
</article>
