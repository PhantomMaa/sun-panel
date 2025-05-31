import { post } from '@/utils/request'

// 定义接口类型
export interface PublicVisitCodeData {
  code: string
}

export interface ApiResponse<T> {
  data: T
  msg: string
  code: number
}

// 启用公开访问代码
export function enablePublicVisit() {
  return post<ApiResponse<PublicVisitCodeData>>({
    url: '/panel/publicVisit/enable',
  })
}

// 禁用公开访问代码
export function disablePublicVisit() {
  return post<ApiResponse<{}>>({
    url: '/panel/publicVisit/disable',
  })
}

