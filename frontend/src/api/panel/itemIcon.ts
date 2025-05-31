import { get, post } from '@/utils/request'

export function addMultiple<T>(req: Panel.ItemInfo[]) {
  return post<T>({
    url: '/panel/itemIcon/addMultiple',
    data: req,
  })
}

export function edit<T>(req: Panel.ItemInfo) {
  return post<T>({
    url: '/panel/itemIcon/edit',
    data: req,
  })
}

export function getListByGroupId<T>(itemIconGroupId: number | undefined) {
  return get<T>({
    url: '/panel/itemIcon/getIcons',
    data: { itemIconGroupId },
  })
}

export function deleteItem<T>(id: number) {
  return post<T>({
    url: '/panel/itemIcon/delete',
    data: { id },
  })
}

export function saveSort<T>(data: Panel.ItemIconSortRequest) {
  return post<T>({
    url: '/panel/itemIcon/saveSort',
    data,
  })
}

export function getSiteFavicon<T>(url: string) {
  return post<T>({
    url: '/panel/itemIcon/getSiteFavicon',
    data: { url },
  })
}
