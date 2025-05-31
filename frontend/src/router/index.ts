import type { App } from 'vue'
import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router'
import { setupPageGuard } from './permission'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/home/index.vue'),
  },

  {
    path: '/login',
    name: 'login',
    component: () => import('../views/login/index.vue'),
  },

  {
    path: '/404',
    name: '404',
    component: () => import('../views/exception/404/index.vue'),
  },

  {
    path: '/500',
    name: '500',
    component: () => import('../views/exception/500/index.vue'),
  },

  {
    path: '/test',
    name: 'test',
    component: () => import('../views/exception/test/index.vue'),
  },

  // 专门处理公开访问代码的路由
  // 匹配 10 位字母数字的路径，可以带或不带斜杠
  {
    path: '/:code([a-zA-Z0-9]{10})',
    name: 'PublicAccess',
    component: () => import('../views/home/index.vue'),
  },

  // 必须放在最后，作为兜底
  {
    path: '/:pathMatch(.*)*',
    name: 'notFound',
    redirect: '/404',
  },

  // adminRouter,
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ left: 0, top: 0 }),
})

setupPageGuard(router)

export async function setupRouter(app: App) {
  app.use(router)
  await router.isReady()
}
