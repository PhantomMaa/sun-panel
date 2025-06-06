<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import { NButton, NDataTable, NDropdown, NTag, useDialog, useMessage } from 'naive-ui'
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import { deleteUser as usersDeleteUser, getList as usersGetList } from '../../../api/panel/users'
import { SvgIcon } from '../../common'
import EditUser from './EditUser/index.vue'
import { useAuthStore } from '@/store'
import { t } from '@/locales'
import { AdminAuthRole } from '@/enums/admin'

const message = useMessage()
const authStore = useAuthStore()
const tableIsLoading = ref<boolean>(false)
const editUserDialogShow = ref<boolean>(false)
const editUserUserInfo = ref<User.Info>()
const dialog = useDialog()

const createColumns = ({
  update,
}: {
  update: (row: User.Info) => void
}): DataTableColumns<User.Info> => {
  return [
    {
      title: t('common.username'),
      key: 'username',
      render(row: User.Info) {
        if (row.username === authStore.userInfo?.username)
          return `${row.username} (${t('adminSettingUsers.currentUseUsername')})`
        return row.username
      },
    },
    {
      title: t('common.nickName'),
      key: 'name',
    },
    {
      title: t('adminSettingUsers.role'),
      key: 'role',
      render(row) {
        switch (row.role) {
          case AdminAuthRole.admin:
            return h(NTag, { type: 'info' }, t('common.role.admin'))
          case AdminAuthRole.regularUser:
            return h(NTag, t('common.role.regularUser'))
          default:
            return '-'
        }
      },
    },
    {
      title: t('adminSettingUsers.provider'),
      key: 'oauthProvider',
      render(row) {
        return h('span', row.oauthProvider)
      },
    },
    {
      title: t('common.action'),
      key: '',
      render(row) {
        const btn = h(
          NButton,
          {
            strong: true,
            tertiary: true,
            size: 'small',
          },
          {
            default() {
              return h(
                SvgIcon, {
                  icon: 'mingcute:more-1-fill',
                },
              )
            },
          },
        )

        return h(NDropdown, {
          trigger: 'click',
          onSelect(key: string | number) {
            switch (key) {
              case 'update':
                update(row)
                break
              case 'delete':
                dialog.warning({
                  title: t('common.warning'),
                  content: t('adminSettingUsers.deletePromptContent', { name: row.name, username: row.username }),
                  positiveText: t('common.confirm'),
                  negativeText: t('common.cancel'),
                  onPositiveClick: () => {
                    deleteUser(row.id as number)
                  },
                })
                break

              default:
                break
            }
          },
          options: [
            {
              label: t('common.edit'),
              key: 'update',
            },
            {
              label: t('common.delete'),
              key: 'delete',
            },
          ],
        }, { default: () => btn })
      },
    },
  ]
}

const userList = ref<User.Info[]>()

const columns = createColumns({
  update(row: User.Info) {
    editUserUserInfo.value = row
    editUserDialogShow.value = true
  },
})
const pagination = reactive({
  page: 1,
  showSizePicker: true,
  pageSizes: [10, 30, 50, 100, 200],
  pageSize: 10,
  itemCount: 0,
  onChange: (page: number) => {
    pagination.page = page
    getList(null)
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize
    pagination.page = 1
    getList(null)
  },
  prefix(item: PaginationProps) {
    return t('adminSettingUsers.userCountText', { count: item.itemCount })
  },
})

// 添加
function handleAdd() {
  editUserDialogShow.value = true
  editUserUserInfo.value = {}
}

function handelDone() {
  editUserDialogShow.value = false
  message.success(t('common.success'))
  getList(null)
}

async function getList(page: number | null) {
  tableIsLoading.value = true
  const req: AdminUserManage.GetListRequest = {
    page: page || pagination.page,
    limit: pagination.pageSize,
  }
  const { data } = await usersGetList<Common.ListResponse<User.Info[]>>(req)
  pagination.itemCount = data.count
  if (data.list)
    userList.value = data.list
  tableIsLoading.value = false
}

async function deleteUser(id: number) {
  const { code } = await usersDeleteUser(id)
  if (code === 0) {
    message.success(t('common.deleteSuccess'))
    await getList(null)
  }
}

onMounted(() => {
  getList(null)
})
</script>

<template>
  <div class="overflow-auto pt-2">
    <div class="my-[10px]">
      <NButton type="primary" size="small" ghost @click="handleAdd">
        {{ $t('common.add') }}
      </NButton>
    </div>

    <NDataTable
      :columns="columns"
      :data="userList"
      :pagination="pagination"
      :bordered="false"
      :loading="tableIsLoading"
      :remote="true"
    />
    <EditUser v-model:visible="editUserDialogShow" :user-info="editUserUserInfo" @done="handelDone" />
  </div>
</template>
