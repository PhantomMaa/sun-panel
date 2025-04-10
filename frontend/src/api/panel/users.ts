import { get, post } from '@/utils/request'

// For admin to manage users

export function edit<T>(param: User.Info) {
  let url = '/panel/users/create'
  if (param.id)
    url = '/panel/users/update'

  return post<T>({
    url,
    data: param,
  })
}

export function getList<T>(param: AdminUserManage.GetListRequest) {
  return get<T>({
    url: '/panel/users/getList',
    data: param,
  })
}

export function deleteUser<T>(userId: number) {
  return post<T>({
    url: '/panel/users/delete',
    data: { userId },
  })
}
