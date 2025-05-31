export type Theme = 'light' | 'dark' | 'auto'

export type Language = 'zh-CN' | 'zh-TW' | 'en-US' | 'ko-KR'

export interface AdminState {
  siderCollapsed: boolean
  theme: Theme
  language: Language
}

export function defaultSetting(): AdminState {
  return { siderCollapsed: false, theme: 'light', language: 'zh-CN' }
}
