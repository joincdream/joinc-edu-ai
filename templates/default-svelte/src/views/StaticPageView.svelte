<script lang="ts">
  import type { PageExportData } from '$types/site';
  import type { MessageBundle } from '$types/messages';
  import type { AlternateLink } from '$types/site';
  import BaseLayout from '$layouts/BaseLayout.svelte';

  interface Props {
    page: PageExportData;
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
    page,
    currentLang,
    baseUrl,
    messages,
    switchUrl,
    alternateLangs = []
  }: Props = $props();

  const siteTitle = $derived(messages.common?.site_title || 'AI Info');
</script>

<BaseLayout
  title={page.title}
  {siteTitle}
  {currentLang}
  {baseUrl}
  currentPath="{baseUrl}{page.slug}/"
  {messages}
  {switchUrl}
  {alternateLangs}
>
  <article class="max-w-4xl mx-auto py-6">
    <!-- Page Header -->
    <header class="mb-10 pb-6 border-b border-slate-200">
      <h1 class="text-3xl sm:text-4xl font-extrabold text-slate-900 tracking-tight leading-tight mb-4">
        {page.title}
      </h1>
    </header>

    <!-- Page Body -->
    <div class="prose max-w-none">
      {@html page.html_content}
    </div>
  </article>
</BaseLayout>
