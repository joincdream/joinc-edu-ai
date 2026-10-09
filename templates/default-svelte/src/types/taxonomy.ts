/**
 * 카테고리 내보내기 데이터 DTO 인터페이스
 */
export interface CategoryExportData {
  name: string;
  slug: string;
  post_count: number;
  post_slugs: string[];
}
