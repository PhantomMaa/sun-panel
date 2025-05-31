import type { AuthState } from './index'
import { ss } from '@/utils/storage'

const LOCAL_NAME = 'authStorage'

export function setStorage(state: AuthState) {
  return ss.set(LOCAL_NAME, state)
}

export function getStorage() {
  return ss.get(LOCAL_NAME)
}

export function removeStorage() {
  ss.remove(LOCAL_NAME)
}
