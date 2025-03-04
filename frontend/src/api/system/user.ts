import {get, post} from '@/utils/request'

// For current user, update himself info

export function getAuthInfo<T>() {
  return get<T>({
    url: '/user/getAuthInfo',
  })
}

export function updateInfo<T>(name: string) {
  return post<T>({
    url: '/user/updateInfo',
    data: { name },
  })
}

export function updatePassword<T>(oldPassword: string, newPassword: string) {
  return post<T>({
    url: '/user/updatePassword',
    data: { newPassword, oldPassword },
  })
}
