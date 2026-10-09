import { createServer } from 'vite';
import { render } from 'svelte/server';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import type { SiteDataBundle, LanguageDataBundle, AlternateLink } from '../src/types/site';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const rootDir = path.resolve(__dirname, '..');

// CLI 파라미터 파싱
function parseArgs() {
  const args = process.argv.slice(2);
  let dataPath = path.resolve(rootDir, '../../dist/.cache/site-data.json');
  let outDir = path.resolve(rootDir, '../../dist');

  for (let i = 0; i < args.length; i++) {
    if (args[i] === '--data' && args[i + 1]) {
      dataPath = path.resolve(process.cwd(), args[i + 1]);
      i++;
    } else if (args[i] === '--out' && args[i + 1]) {
      outDir = path.resolve(process.cwd(), args[i + 1]);
      i++;
    }
  }

  // fallback to .cache/site-data.json if dist/.cache/site-data.json doesn't exist
  if (!fs.existsSync(dataPath)) {
    const fallback = path.resolve(rootDir, '../../.cache/site-data.json');
    if (fs.existsSync(fallback)) {
      dataPath = fallback;
    }
  }

  return { dataPath, outDir };
}

function copyDirRecursive(src: string, dest: string) {
  if (!fs.existsSync(src)) return;
  fs.mkdirSync(dest, { recursive: true });
  const entries = fs.readdirSync(src, { withFileTypes: true });
  for (const entry of entries) {
    const srcPath = path.join(src, entry.name);
    const destPath = path.join(dest, entry.name);
    if (entry.isDirectory()) {
      copyDirRecursive(srcPath, destPath);
    } else {
      fs.copyFileSync(srcPath, destPath);
    }
  }
}

async function buildSSG() {
  const startTime = Date.now();
  const { dataPath, outDir } = parseArgs();

  console.log(`[SSG] Reading site data from: ${dataPath}`);
  if (!fs.existsSync(dataPath)) {
    console.error(`[SSG] Error: Site data bundle not found at ${dataPath}`);
    process.exit(1);
  }

  const rawData = fs.readFileSync(dataPath, 'utf-8');
  const siteData: SiteDataBundle = JSON.parse(rawData);

  console.log(`[SSG] Target output directory: ${outDir}`);

  // Vite 프로덕션 SSR 런타임 생성
  process.env.NODE_ENV = 'production';
  const vite = await createServer({
    root: rootDir,
    mode: 'production',
    server: { middlewareMode: true },
    appType: 'custom',
  });

  try {
    // 뷰 컴포넌트 로드
    const { default: HomeView } = await vite.ssrLoadModule('./src/views/HomeView.svelte');
    const { default: DetailView } = await vite.ssrLoadModule('./src/views/DetailView.svelte');
    const { default: CategoryView } = await vite.ssrLoadModule('./src/views/CategoryView.svelte');
    const { default: StaticPageView } = await vite.ssrLoadModule('./src/views/StaticPageView.svelte');

    let totalRendered = 0;

    const renderPage = (Component: any, props: any, targetHtmlPath: string, lang: string) => {
      const result = render(Component, { props });
      const html = `<!DOCTYPE html>
<html lang="${lang}">
<head>
${result.head}
</head>
<body>
${result.body}
</body>
</html>`;

      const dir = path.dirname(targetHtmlPath);
      fs.mkdirSync(dir, { recursive: true });
      fs.writeFileSync(targetHtmlPath, html, 'utf-8');
      totalRendered++;
    };

    const languages = siteData.languages;
    const langKeys = Object.keys(languages);

    for (const lang of langKeys) {
      const bundle: LanguageDataBundle = languages[lang];
      const isKorean = lang === 'ko';
      const baseOutDir = isKorean ? outDir : path.join(outDir, lang);
      const baseUrl = bundle.base_url;

      const peerLang = isKorean ? 'en' : 'ko';
      const peerBundle = languages[peerLang];

      // 포스트 목록 최신 날짜순(내림차순, Descending) 정렬 확정
      bundle.posts.sort((a, b) => {
        const dateA = a.frontmatter.created_date || '';
        const dateB = b.frontmatter.created_date || '';
        return dateB.localeCompare(dateA);
      });

      console.log(`[SSG] Pre-rendering [${lang.toUpperCase()}] site (Base URL: ${baseUrl})...`);

      // 1. HomeView 렌더링
      const homeSwitchUrl = {
        ko: '/',
        en: '/en/'
      };
      const homeAlternateLangs: AlternateLink[] = [
        { lang: 'ko', url: '/' },
        { lang: 'en', url: '/en/' },
        { lang: 'x-default', url: '/' }
      ];

      renderPage(HomeView, {
        posts: bundle.posts,
        categories: bundle.categories,
        currentLang: lang,
        baseUrl,
        messages: bundle.messages,
        switchUrl: homeSwitchUrl,
        alternateLangs: homeAlternateLangs
      }, path.join(baseOutDir, 'index.html'), lang);

      // 2. DetailView (포스트 상세) 렌더링
      const peerPostMap = new Map((peerBundle?.posts || []).map(p => [p.slug, p]));

      for (const post of bundle.posts) {
        const hasPeer = peerPostMap.has(post.slug);
        const switchUrl = {
          ko: `/posts/${post.slug}/`,
          en: hasPeer ? `/en/posts/${post.slug}/` : '/en/'
        };
        if (!isKorean) {
          switchUrl.ko = hasPeer ? `/posts/${post.slug}/` : '/';
          switchUrl.en = `/en/posts/${post.slug}/`;
        }

        const altLangs: AlternateLink[] = [];
        if (hasPeer) {
          altLangs.push(
            { lang: 'ko', url: `/posts/${post.slug}/` },
            { lang: 'en', url: `/en/posts/${post.slug}/` },
            { lang: 'x-default', url: `/posts/${post.slug}/` }
          );
        }

        const destFile = path.join(baseOutDir, 'posts', post.slug, 'index.html');
        renderPage(DetailView, {
          post,
          currentLang: lang,
          baseUrl,
          messages: bundle.messages,
          switchUrl,
          alternateLangs: altLangs
        }, destFile, lang);
      }

      // 3. CategoryView (카테고리 아카이브) 렌더링
      const allPostsCount = bundle.posts.length;
      for (const cat of bundle.categories) {
        const catPosts = bundle.posts.filter(p => cat.post_slugs.includes(p.slug));
        const switchUrl = {
          ko: `/category/${cat.slug}/`,
          en: `/en/category/${cat.slug}/`
        };

        const destFile = path.join(baseOutDir, 'category', cat.slug, 'index.html');
        renderPage(CategoryView, {
          category: cat,
          posts: catPosts,
          allPostsCount,
          categories: bundle.categories,
          currentLang: lang,
          baseUrl,
          messages: bundle.messages,
          switchUrl
        }, destFile, lang);
      }

      // 4. StaticPageView (About 등 독립 페이지) 렌더링
      const peerPageMap = new Map((peerBundle?.pages || []).map(p => [p.slug, p]));
      for (const page of bundle.pages || []) {
        const hasPeer = peerPageMap.has(page.slug);
        const switchUrl = {
          ko: `/${page.slug}/`,
          en: hasPeer ? `/en/${page.slug}/` : '/en/'
        };
        if (!isKorean) {
          switchUrl.ko = hasPeer ? `/${page.slug}/` : '/';
          switchUrl.en = `/en/${page.slug}/`;
        }

        const altLangs: AlternateLink[] = [];
        if (hasPeer) {
          altLangs.push(
            { lang: 'ko', url: `/${page.slug}/` },
            { lang: 'en', url: `/en/${page.slug}/` },
            { lang: 'x-default', url: `/${page.slug}/` }
          );
        }

        const destFile = path.join(baseOutDir, page.slug, 'index.html');
        renderPage(StaticPageView, {
          page,
          currentLang: lang,
          baseUrl,
          messages: bundle.messages,
          switchUrl,
          alternateLangs: altLangs
        }, destFile, lang);
      }
    }

    // 5. 테마 정적 에셋 복사 (templates/default-svelte/assets/ -> dist/assets/)
    const srcAssets = path.join(rootDir, 'assets');
    const destAssets = path.join(outDir, 'assets');
    copyDirRecursive(srcAssets, destAssets);

    const elapsed = Date.now() - startTime;
    console.log(`[SSG] [SUCCESS] Completed pre-rendering ${totalRendered} static HTML pages in ${elapsed}ms!`);
  } finally {
    await vite.close();
  }
}

buildSSG().catch(err => {
  console.error('[SSG] Pre-rendering failed:', err);
  process.exit(1);
});
