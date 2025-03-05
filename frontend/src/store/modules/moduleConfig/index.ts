import { defineStore } from 'pinia'
import type { ModuleConfigState } from './helper'
import { getLocalState, setLocalState } from './helper'
import { getValueByName, save } from '@/api/system/moduleConfig'

export const useModuleConfig = defineStore('module-config-store', {
  state: (): ModuleConfigState => getLocalState(),
  actions: {
    // 获取值
    async getValueByNameFromCloud<T>(name: string) {
      const moduleName = `module-${name}`
      return await getValueByName<T>(moduleName)
    },

    // 保存到网络
    async saveToCloud(name: string, value: any) {
      const moduleName = `module-${name}`
      // 保存至网络
      return save(moduleName, value)
    },

    recordState() {
      setLocalState(this.$state)
    },
  },
})
