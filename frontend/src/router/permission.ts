import type { Router } from 'vue-router'
import { useAuthStore } from '@/store'

export function setupPageGuard(router: Router) {
  router.beforeEach(async (to, from, next) => {
    // 获取当前token和用户信息
    const authStore = useAuthStore()

    if (to.path === '/login' && authStore.token) {
      next('/')
      return
    }

    // 其他情况正常放行
    next()
  })
}
