/**
 * 문자열을 URL-Safe 슬러그로 변환 (Go 백엔드 model.Slugify와 1:1 일치)
 */
export function slugify(str: string): string {
  if (!str) return 'default';
  const slug = str
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9가-힣]+/g, '-')
    .replace(/^-+|-+$/g, '');
  return slug || 'default';
}
