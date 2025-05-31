import type { Router } from 'vue-router'
import { useAuthStore } from '@/store'

export function setupPageGuard(router: Router) {
  router.beforeEach(async (to, from, next) => {
    // 获取当前token和用户信息
    const authStore = useAuthStore()
    const token = authStore.token

    // 如果是登录页且有token，直接跳转到首页
    if (to.path === '/login' && token) {
      next('/')
      return
    }

    // 非管理员路由拦截
    if (authStore.userInfo && authStore.userInfo.role !== 1 && to.path.includes('admin')) {
      console.log('Non-admin trying to access admin page')
      next({ name: '404' })
      return
    }

    // 其他情况正常放行
    next()
  })
}
