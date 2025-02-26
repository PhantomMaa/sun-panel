import { get, post } from '../../utils/request'

export function getList<T>() {
  return get<T>({
    url: '/file/getList',
  })
}

export function deletes<T>(ids: number[]) {
  return post<T>({
    url: '/file/deletes',
    data: { ids },
  })
}
