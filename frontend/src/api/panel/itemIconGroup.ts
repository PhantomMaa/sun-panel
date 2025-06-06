import { get, post } from '../../utils/request'

export function edit<T>(req: Panel.ItemIconGroup) {
  return post<T>({
    url: '/panel/itemIconGroup/edit',
    data: req,
  })
}

export function getList<T>() {
  return get<T>({
    url: '/panel/itemIconGroup/getGroups',
  })
}

export function deletes<T>(ids: number[]) {
  return post<T>({
    url: '/panel/itemIconGroup/deletes',
    data: { ids },
  })
}

export function saveSort<T>(sortItems: Common.SortItemRequest[]) {
  return post<T>({
    url: '/panel/itemIconGroup/saveSort',
    data: { sortItems },
  })
}
