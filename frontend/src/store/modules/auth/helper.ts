import type { AuthState } from './index'
import { ss } from '@/utils/storage'

const LOCAL_NAME = 'userStorage'

export function setStorage(state: AuthState) {
  return ss.set(LOCAL_NAME, state)
}

export function getStorage() {
  return ss.get(LOCAL_NAME)
}

export function removeToken() {
  ss.remove(LOCAL_NAME)
}
