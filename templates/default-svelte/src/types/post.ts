/**
 * Frontmatter 메타데이터 인터페이스
 */
export interface Frontmatter {
  title: string;
  category?: string;
  tags?: string[];
  summary?: string;
  description?: string;
  created_date?: string;
  updated_date?: string;
  status?: string;
  is_draft?: boolean;
  slug?: string;
}

/**
 * 목차(TOC) 트리 항목 인터페이스
 */
export interface TOCItem {
  id: string;
  title: string;
  level: number;
  children?: TOCItem[];
}

/**
 * 포스트 내보내기 데이터 DTO 인터페이스
 */
export interface PostExportData {
  id: string;
  slug: string;
  file_path: string;
  frontmatter: Frontmatter;
  html_content: string;
  toc: TOCItem[];
  reading_time_minutes: number;
  excerpt: string;
  lang: string;
  alternate_url?: string;
}
