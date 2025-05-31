import { defineStore } from 'pinia'
import { getStorage, removeStorage as hRemoveStorage, setStorage } from './helper'

export interface AuthState {
  userInfo: User.Info | null
}

const defaultState: AuthState = {
  userInfo: null,
}

export const useAuthStore = defineStore('auth-store', {
  state: (): AuthState => getStorage() || defaultState,

  actions: {
    setUserInfo(userInfo: User.Info) {
      this.userInfo = userInfo
      this.saveStorage()
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
