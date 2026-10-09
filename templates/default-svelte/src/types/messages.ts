/**
 * 다국어 UI 메시지 번들 인터페이스 (messages.yaml 매핑)
 */
export interface MessageBundle {
  common: Record<string, string>;
  nav: Record<string, string>;
  card: Record<string, string>;
  detail: Record<string, string>;
  empty: Record<string, string>;
  footer: Record<string, string>;
  banner: Record<string, string>;
}
