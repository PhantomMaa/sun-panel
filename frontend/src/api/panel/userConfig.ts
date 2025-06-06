import { get, post } from '@/utils/request'

export function set<T>(req: Panel.userConfig) {
  return post<T>({
    url: '/panel/userConfig/setConfig',
    data: req,
  })
}

export function getUserConfig<T>() {
  return get<T>({
    url: '/panel/userConfig/getConfig',
  })
}
