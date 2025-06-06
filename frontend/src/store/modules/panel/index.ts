import { defineStore } from 'pinia'
import { defaultState, defaultStatePanelConfig, getLocalState, removeLocalState, setLocalState } from './helper'
import { getUserConfig } from '@/api/panel/userConfig'
import { router } from '@/router'
import type { PanelStateNetworkModeEnum } from '@/enums'
export const usePanelState = defineStore('panel', {
  state: (): Panel.State => getLocalState() || defaultState(),

  getters: {

  },

  actions: {
    setLeftSiderCollapsed(Collapsed: boolean) {
      this.leftSiderCollapsed = Collapsed
      // this.recordState()
    },

    setRightSiderCollapsed(Collapsed: boolean) {
      this.rightSiderCollapsed = Collapsed
      // this.recordState()
    },

    setNetworkMode(mode: PanelStateNetworkModeEnum) {
      this.networkMode = mode
      this.recordState()
    },

    // 获取云端（搭建的服务器）的面板配置
    updatePanelConfigByCloud() {
      getUserConfig<Panel.userConfig>().then((res) => {
        if (res.code === 0)
          this.panelConfig = { ...defaultStatePanelConfig(), ...res.data.panel }
        else
          this.resetPanelConfig()
        this.recordState()
      }).catch(() => {
        this.resetPanelConfig()
        this.recordState()
      })
    },

    resetPanelConfig() {
      this.panelConfig = defaultStatePanelConfig()
    },

    async reloadRoute(id?: number) {
      await router.push({ name: 'AppletDialog', params: { aiAppletId: id } })
    },

    recordState() {
      setLocalState(this.$state)
    },

    removeState() {
      removeLocalState()
    },
  },
})
