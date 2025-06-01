import { defineStore } from 'pinia'
import { getStorage, removeStorage as hRemoveStorage, setStorage } from './helper'

export interface AuthState {
  userInfo: User.Info | null
  token: string | null
}

const defaultState: AuthState = {
  userInfo: null,
  token: '',
}

export const useAuthStore = defineStore('auth-store', {
  state: (): AuthState => getStorage() || defaultState,

  actions: {
    setToken(token: string) {
      this.token = token
    },

    setUserInfo(userInfo: User.Info) {
      this.userInfo = userInfo
    },

    saveStorage() {
      setStorage(this.$state)
    },

    removeStorage() {
      this.$state = defaultState
      hRemoveStorage()
    },
  },

})
