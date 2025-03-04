import { get } from '@/utils/request'

export function getAbout<T>() {
  return get<T>({
    url: '/about',
  })
}
