import moment from 'moment'
import { useAuthStore, useUserStore } from '@/store'
import { getUser } from '@/api/system/user'

const userStore = useUserStore()
const authStore = useAuthStore()

/**
 * 生成指定时间格式
 * @param format 时间格式 默认：'YYYY-MM-DD HH:mm:ss'
 * @returns string
 */
export function buildTimeString(format?: string): string {
  if (!format)
    format = 'YYYY-MM-DD HH:mm:ss'

  return moment().format(format)
}

export function timeFormat(timeString?: string) {
  return moment(timeString).format('YYYY-MM-DD HH:mm:ss')
}

export function setTitle(titile: string) {
  document.title = titile
}

export function getTitle(titile: string) {
  document.title = titile
}

export async function updateLocalUserInfo() {
  try {
    const { data } = await getUser()
    if (data) {
      userStore.updateUserInfo({ name: data.name, logined: data.logined })
      authStore.setUserInfo(data)
    }
  }
  catch (error) {
    console.error('Failed to update local user info:', error)
  }
}

// 复制文字到剪切板
export async function copyToClipboard(text: string): Promise<boolean> {
  if (navigator.clipboard) {
    // 使用 Clipboard API
    try {
      await navigator.clipboard.writeText(text)
      return true
    }
    catch (err) {
      console.error('copy fail', err)
      return false
    }
  }
  else {
    // 兼容旧版浏览器
    const textArea = document.createElement('textarea')
    textArea.value = text
    document.body.appendChild(textArea)
    textArea.select()

    try {
      document.execCommand('copy')
      return true
    }
    catch (err) {
      console.error('copy fail', err)
      return false
    }
    finally {
      document.body.removeChild(textArea)
    }
  }
}

export function bytesToSize(bytes: number) {
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB']
  if (bytes === 0)
    return '0B'
  const i = parseInt(String(Math.floor(Math.log(bytes) / Math.log(1024))))
  return `${(bytes / 1024 ** i).toFixed(1)} ${sizes[i]}`
}
