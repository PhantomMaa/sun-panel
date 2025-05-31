import {get, post} from '@/utils/request'

// For current user, update himself info

export function getUser() {
  return get<User.Info>({
    url: '/user/getUser',
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
