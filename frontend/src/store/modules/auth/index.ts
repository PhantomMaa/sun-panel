import { defineStore } from 'pinia'
import { getStorage, removeToken as hRemoveToken, setStorage } from './helper'

export interface AuthState {
  token: string | null
  userInfo: User.Info | null
}

const defaultState: AuthState = {
  token: null,
  userInfo: null,
}

export const useAuthStore = defineStore('auth-store', {
  state: (): AuthState => getStorage() || defaultState,

  actions: {
    setToken(token: string) {
      this.token = token
      this.saveStorage()
    },

    setUserInfo(userInfo: User.Info) {
      this.userInfo = userInfo
      this.saveStorage()
    },

    saveStorage() {
      setStorage(this.$state)
    },

    removeToken() {
      this.$state = defaultState
      hRemoveToken()
    },
  },

})
