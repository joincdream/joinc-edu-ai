import type { PostExportData } from './post';
import type { CategoryExportData } from './taxonomy';
import type { MessageBundle } from './messages';

/**
 * 독립 단일 페이지 DTO 인터페이스 (About 등)
 */
export interface PageExportData {
  slug: string;
  title: string;
  html_content: string;
  lang: string;
  alternate_url?: string;
}

/**
 * 언어별 데이터 세트 번들
 */
export interface LanguageDataBundle {
  lang: string;
  base_url: string;
  messages: MessageBundle;
  categories: CategoryExportData[];
  posts: PostExportData[];
  pages: PageExportData[];
}

/**
 * 전체 사이트 데이터 전송 계약 객체 (SSOT)
 */
export interface SiteDataBundle {
  version: string;
  generated_at: string;
  languages: {
    ko: LanguageDataBundle;
    en: LanguageDataBundle;
    [key: string]: LanguageDataBundle;
  };
}

/**
 * SEO 대체 언어 링크
 */
export interface AlternateLink {
  lang: string;
  url: string;
}
