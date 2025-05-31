<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NCard, NDivider, NForm, NFormItem, NInput, NSelect, useDialog, useMessage, NSwitch } from 'naive-ui'
import { ref, onMounted } from 'vue'
import { RoundCardModal, SvgIcon } from '../../common'
import { useAppStore, useAuthStore, usePanelState, useUserStore } from '@/store'
import { languageOptions } from '@/utils/defaultData'
import type { Language, Theme } from '@/store/modules/app/helper'
import { logout } from '@/api'
import { updateInfo, updatePassword } from '@/api/system/user'
import { enablePublicVisit, disablePublicVisit } from '@/api/panel/publicVisit'
import { updateLocalUserInfo } from '@/utils/cmn'
import { t } from '@/locales'

// 使用导入的 ApiResponse 类型

const userStore = useUserStore()
const authStore = useAuthStore()
const appStore = useAppStore()
const panelState = usePanelState()
const ms = useMessage()
const dialog = useDialog()

const languageValue = ref(appStore.language)
const themeValue = ref(appStore.theme)
const nickName = ref(authStore.userInfo?.name || '')
const isEditNickNameStatus = ref(false)
const formRef = ref<FormInst | null>(null)
const publicVisitEnabled = ref<boolean>(false)
const publicVisitUrl = ref<string>('')
const publicVisitLoading = ref<boolean>(false)
const themeOptions: { label: string; key: string; value: Theme }[] = [
  { label: t('apps.userInfo.themeStyle.dark'), key: 'dark', value: 'dark' },
  { label: t('apps.userInfo.themeStyle.light'), key: 'light', value: 'light' },
  { label: t('apps.userInfo.themeStyle.auto'), key: 'Auto', value: 'auto' },
]
const updatePasswordModalState = ref({
  show: false,
  loading: false,
  form: {
    password: '',
    oldPassword: '',
    confirmPassword: '',
  },
})

const updatePasswordModalFormRules: FormRules = {
  oldPassword: {
    required: true,
    trigger: 'blur',
    min: 6,
    max: 20,
    message: t('adminSettingUsers.formRules.passwordLimit'),
  },
  password: {
    required: true,
    trigger: 'blur',
    min: 6,
    max: 20,
    message: t('adminSettingUsers.formRules.passwordLimit'),
  },
  confirmPassword: {
    required: true,
    trigger: 'blur',
    min: 6,
    max: 20,
    message: t('adminSettingUsers.formRules.passwordLimit'),
  },
}

async function logoutApi() {
  await logout()
  userStore.resetUserInfo()
  authStore.removeToken()
  panelState.removeState()
  appStore.removeToken()
  ms.success(t('settingUserInfo.logoutSuccess'))
  // router.push({ path: '/login' })
  location.reload()// 强制刷新一下页面
}

function handleSaveInfo() {
  updateInfo(nickName.value).then(({ code, msg }) => {
    if (code === 0) {
      updateLocalUserInfo()
      isEditNickNameStatus.value = false
    }
    else {
      ms.error(`${t('common.editFail')}:${msg}`)
    }
  })
}

function handleUpdatePassword(e: MouseEvent) {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (errors) {
      console.log(errors)
      return
    }

    if (updatePasswordModalState.value.form.password !== updatePasswordModalState.value.form.confirmPassword) {
      ms.error(t('settingUserInfo.confirmPasswordInconsistentMsg'))
      return
    }
    updatePasswordModalState.value.loading = true
    updatePassword(updatePasswordModalState.value.form.oldPassword, updatePasswordModalState.value.form.password).then(({ code, msg }) => {
      if (code === 0) {
        // 成功
        updatePasswordModalState.value.show = false
        ms.success(t('common.success'))
      }
    }).finally(() => {
      updatePasswordModalState.value.loading = false
    }).catch(() => {
      ms.error(t('common.serverError'))
    })
  })
}

function handleLogout() {
  dialog.warning({
    title: t('common.warning'),
    content: t('settingUserInfo.confirmLogoutText'),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => {
      logoutApi()
    },
  })
}

function handleChangeLanuage(value: Language) {
  languageValue.value = value
  appStore.setLanguage(value)
  location.reload()
}

function handleChangeTheme(value: Theme) {
  themeValue.value = value
  appStore.setTheme(value)
  // location.reload()
}

// 组件挂载时获取公开访问代码状态
onMounted(() => {
  fetchPublicVisitCode()
})

// 从用户认证信息中获取公开访问代码状态
const fetchPublicVisitCode = () => {
  try {
    // 直接从 authStore 中获取公开访问代码
    const publiccode = authStore.userInfo?.publiccode
    if (publiccode) {
      // 使用类型断言确保类型安全
      const codeStr = String(publiccode)
      publicVisitEnabled.value = true
      // 构建公开访问URL
      const baseUrl = window.location.origin
      publicVisitUrl.value = `${baseUrl}/${codeStr}`
    } else {
      publicVisitEnabled.value = false
      publicVisitUrl.value = ''
    }
  } catch (error) {
    console.error('获取公开访问代码失败:', error)
    ms.error('获取公开访问代码状态失败')
  }
}

// 处理公开访问开关切换
const handleTogglePublicVisit = async (value: boolean) => {
  publicVisitLoading.value = true
  try {
    if (value) {
      // 开启公开访问
      const res = await enablePublicVisit()
      if (res.code === 0 && res.data && res.data.code) {
        // 使用类型断言确保类型安全
        const codeStr = String(res.data.code)
        publicVisitEnabled.value = true
        // 构建公开访问URL
        const baseUrl = window.location.origin
        publicVisitUrl.value = `${baseUrl}/${codeStr}`
        ms.success('公开访问已开启')
        // 更新用户信息
        updateLocalUserInfo()
      }
    } else {
      // 关闭公开访问
      await disablePublicVisit()
      publicVisitEnabled.value = false
      publicVisitUrl.value = ''
      ms.success('公开访问已关闭')
      // 更新用户信息
      updateLocalUserInfo()
    }
  } catch (error) {
    console.error('切换公开访问状态失败:', error)
    ms.error('切换公开访问状态失败')
    publicVisitEnabled.value = !value // 恢复开关状态
  } finally {
    publicVisitLoading.value = false
  }
}
</script>

<template>
  <div class="bg-slate-200 dark:bg-zinc-900 p-2 h-full">
    <NCard style="border-radius:10px" size="small">
      <div>
        <div class="text-slate-500 font-bold">
          {{ $t('common.username') }}
        </div>
        {{ authStore.userInfo?.username }}
      </div>

      <div class="mt-[10px]">
        <div class="text-slate-500 font-bold">
          {{ $t('common.nickName') }}
        </div>

        <div v-if="!isEditNickNameStatus">
          {{ authStore.userInfo?.name }}

          <NButton size="small" text type="info" @click="isEditNickNameStatus = !isEditNickNameStatus">
            {{ $t('common.edit') }}
          </NButton>
        </div>

        <div v-else class="flex items-center">
          <div class="max-w-[150px]">
            <NInput v-model:value="nickName" type="text" :placeholder="$t('common.inputPlaceholder')" />
          </div>
          <NButton size="small" quaternary type="info" @click="handleSaveInfo">
            {{ $t('common.save') }}
          </NButton>
        </div>
      </div>

      <div class="mt-[10px]">
        <div class="text-slate-500 font-bold">
          {{ $t('common.password') }}
        </div>

        <NButton size="small" text type="info" @click="updatePasswordModalState.show = !updatePasswordModalState.show">
          {{ $t('settingUserInfo.updatePassword') }}
        </NButton>
      </div>

      <NDivider style="margin: 10px 0;" dashed />

      <div class="mt-[10px]">
        <div class="text-slate-500 font-bold">
          {{ $t('common.language') }}
        </div>
        <div class="max-w-[200px]">
          <NSelect v-model:value="languageValue" :options="languageOptions" @update-value="handleChangeLanuage" />
        </div>
      </div>

      <div class="mt-[10px]">
        <div class="text-slate-500 font-bold">
          {{ $t('apps.userInfo.theme') }}
        </div>
        <div class="max-w-[200px]">
          <NSelect v-model:value="themeValue" :options="themeOptions" @update-value="handleChangeTheme" />
        </div>
      </div>

      <NDivider style="margin: 10px 0;" dashed />

      <div class="mt-[10px]">
        <div class="text-slate-500 font-bold">
          {{ $t('apps.userInfo.publicVisit') }}
        </div>
        <div class="max-w-[400px]">
          <div class="flex items-center">
            <span class="mr-[10px]">{{ $t('apps.userInfo.enablePublicVisit') }}</span>
            <NSwitch :value="publicVisitEnabled" :loading="publicVisitLoading" @update:value="handleTogglePublicVisit" />
          </div>
          <div v-if="publicVisitEnabled" class="mt-2">
            <div class="text-sm text-slate-500 mb-1">
              <a :href="publicVisitUrl" target="_blank" class="text-blue-500">{{ publicVisitUrl }}</a>
            </div>
          </div>
        </div>
      </div>
    </NCard>

    <NCard style="border-radius:10px" class="mt-[10px]" size="small">
      <NButton size="small" text type="error" @click="handleLogout">
        <template #icon>
          <SvgIcon icon="tabler:logout" />
        </template>
        {{ $t('settingUserInfo.logout') }}
      </NButton>
    </NCard>

    <RoundCardModal v-model:show="updatePasswordModalState.show" size="small" preset="card" style="width: 400px" :title="$t('settingUserInfo.updatePassword')">
      <NForm ref="formRef" :model="updatePasswordModalState.form" :rules="updatePasswordModalFormRules">
        <NFormItem path="oldPassword" :label="$t('settingUserInfo.oldPassword')">
          <NInput v-model:value="updatePasswordModalState.form.oldPassword" :maxlength="20" type="password" :placeholder="$t('settingUserInfo.oldPassword')" />
        </NFormItem>

        <NFormItem path="password" :label="$t('settingUserInfo.newPassword')">
          <NInput v-model:value="updatePasswordModalState.form.password" :maxlength="20" type="password" :placeholder="$t('settingUserInfo.newPassword')" />
        </NFormItem>

        <NFormItem path="confirmPassword" :label="$t('settingUserInfo.confirmPassword')">
          <NInput v-model:value="updatePasswordModalState.form.confirmPassword" :maxlength="20" type="password" :placeholder="$t('settingUserInfo.confirmPassword')" />
        </NFormItem>
      </NForm>

      <template #footer>
        <div class="float-right">
          <NButton type="success" size="small" :loading="updatePasswordModalState.loading" @click="handleUpdatePassword">
            {{ $t('common.save') }}
          </NButton>
        </div>
      </template>
    </RoundCardModal>
  </div>
</template>
