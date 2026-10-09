import type { MessageBundle } from '$types/messages';

let currentLang = $state<string>('ko');
let messages = $state<MessageBundle | null>(null);

export function setI18n(lang: string, bundle: MessageBundle) {
  currentLang = lang;
  messages = bundle;
}

export function getLang(): string {
  return currentLang;
}

export function getMessages(): MessageBundle | null {
  return messages;
}

/**
 * i18n 메시지 키 조회 헬퍼
 * @param section 'common' | 'nav' | 'card' | 'detail' | 'empty' | 'footer' | 'banner'
 * @param key 메시지 키
 * @param fallback 기본값
 */
export function t(section: keyof MessageBundle, key: string, fallback: string = ''): string {
  if (!messages) return fallback;
  const targetSection = messages[section];
  if (targetSection && targetSection[key] !== undefined && targetSection[key] !== '') {
    return targetSection[key];
  }
  return fallback;
}
