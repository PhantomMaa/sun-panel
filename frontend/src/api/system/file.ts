import { get, post } from '@/utils/request'

export function getList<T>() {
  return get<T>({
    url: '/file/getList',
  })
}

export function deleteFile<T>(id: number) {
  return post<T>({
    url: '/file/delete',
    data: { id },
  })
}
