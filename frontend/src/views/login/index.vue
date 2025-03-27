<script setup lang="ts">
import { NButton, NCard, NForm, NFormItem, NGradientText, NInput, NSelect, useMessage, NDivider } from 'naive-ui'
import { ref, onMounted } from 'vue'
import { login } from '../../api'
import { useAppStore, useAuthStore } from '../../store'
import { SvgIcon, SvgIconOnline } from '../../components/common'
import { router } from '../../router'
import { t } from '../../locales'
import { languageOptions } from '../../utils/defaultData'
import type { Language } from '../../store/modules/app/helper'
import service from '../../utils/request/axios'
import { getAuthInfo } from '@/api/system/user'

const authStore = useAuthStore()
const appStore = useAppStore()
const ms = useMessage()
const loading = ref(false)
const languageValue = ref<Language>(appStore.language)
const oauthEnabled = ref(false)
const oauthProviders = ref<string[]>([])
const oauthLoading = ref(false)
const loadingProvider = ref<string>('')

const form = ref<Login.LoginReqest>({
  username: '',
  password: '',
})

onMounted(async () => {
  try {
    // 获取OAuth配置
    const res = await service.get('oauth/config')
    if (res.data.code === 0) {
      oauthEnabled.value = res.data.data.enabled
      oauthProviders.value = res.data.data.providers || []
    }
  } catch (error) {
    console.error('Failed to fetch OAuth config:', error)
  }

  // 检查URL中是否有token参数（OAuth回调）
  // 从hash部分获取token参数
  const hashParts = window.location.hash.split('?')
  if (hashParts.length > 1) {
    const hashParams = new URLSearchParams(hashParts[1])
    const token = hashParams.get('token')
    
    if (token) {
      try {
        // 直接保存token
        authStore.setToken(token)
        
        try {
          const { data } = await getAuthInfo()
          if (data) {
            authStore.setUserInfo(data)
            // 显示欢迎消息
            ms.success(`Hi ${data.name}, ${t('login.welcomeMessage')}`)
            router.push({ path: '/' })
          } else {
            console.error('Failed to get user info:', data)
          }
        }
        catch (error) {
          console.error('Failed to update local user info:', error)
        }
      } catch (error) {
        console.error('Error during OAuth login:', error)
      }
    }
  }
})

const loginPost = async () => {
  loading.value = true
  try {
    const res = await login<Login.LoginResponse>(form.value)
    if (res.code === 0) {
      authStore.setToken(res.data.token)
      authStore.setUserInfo(res.data)

      setTimeout(() => {
        ms.success(`Hi ${res.data.name},${t('login.welcomeMessage')}`)
        loading.value = false
        router.push({ path: '/' })
      }, 500)
    }
    else {
      loading.value = false
    }
  }
  catch (error) {
    loading.value = false
    // 请检查网络或者服务器错误
    console.log(error)
  }
}

function handleSubmit() {
  // 点击登录按钮触发
  loginPost()
}

function handleChangeLanuage(value: Language) {
  languageValue.value = value
  appStore.setLanguage(value)
}

function handleOAuthLogin(provider: string) {
  oauthLoading.value = true
  loadingProvider.value = provider
  
  // 构建OAuth URL
  const oauthUrl = `/api/oauth/${provider}`
  
  // 直接在当前窗口打开OAuth登录页面
  window.location.href = oauthUrl
}

</script>

<template>
  <div class="login-container">
    <NCard class="login-card" style="border-radius: 20px;">
      <div class="mb-5 flex items-center justify-end">
        <div class="mr-2">
          <SvgIcon icon="ion-language" style="width: 20;height: 20;" />
        </div>
        <div class="min-w-[100px]">
          <NSelect v-model:value="languageValue" size="small" :options="languageOptions" @update-value="handleChangeLanuage" />
        </div>
      </div>

      <div class="login-title  ">
        <NGradientText :size="30" type="success" class="!font-bold">
          {{ $t('common.appName') }}
        </NGradientText>
      </div>
      <NForm :model="form" label-width="100px" @keydown.enter="handleSubmit">
        <NFormItem>
          <NInput v-model:value="form.username" :placeholder="$t('login.usernamePlaceholder')">
            <template #prefix>
              <SvgIcon icon="ph:user-bold" />
            </template>
          </NInput>
        </NFormItem>

        <NFormItem>
          <NInput v-model:value="form.password" type="password" :placeholder="$t('login.passwordPlaceholder')">
            <template #prefix>
              <SvgIcon icon="mdi:password-outline" />
            </template>
          </NInput>
        </NFormItem>

        <NFormItem style="margin-top: 10px">
          <NButton type="primary" block :loading="loading" @click="handleSubmit">
            {{ $t('login.loginButton') }}
          </NButton>
        </NFormItem>

        <!-- OAuth登录按钮 -->
        <div v-if="oauthEnabled && oauthProviders.length > 0">
          <NDivider>{{ $t('login.thirdPartyLogin') }}</NDivider>
          <div class="oauth-buttons flex flex-col gap-2">
            <NButton 
              v-if="oauthProviders.includes('github')" 
              quaternary 
              class="oauth-button" 
              :loading="oauthLoading && loadingProvider === 'github'"
              :disabled="oauthLoading && loadingProvider !== 'github'"
              @click="handleOAuthLogin('github')"
            >
              <template #icon>
                <SvgIconOnline icon="mdi:github" />
              </template>
              GitHub
            </NButton>
            <NButton 
              v-if="oauthProviders.includes('google')" 
              quaternary 
              class="oauth-button" 
              :loading="oauthLoading && loadingProvider === 'google'"
              :disabled="oauthLoading && loadingProvider !== 'google'"
              @click="handleOAuthLogin('google')"
            >
              <template #icon>
                <SvgIconOnline icon="mdi:google" />
              </template>
              Google
            </NButton>
          </div>
        </div>
      </NForm>
    </NCard>
  </div>
</template>

  <style>
    .login-container {
        padding: 20px;
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100vh;
        background-color: #f2f6ff;
    }

    /* 夜间模式 */
    .dark .login-container{
      background-color: rgb(43, 43, 43);
    }

    @media (min-width: 600px) {
        .login-card {
            width: auto;
            margin: 0px 10px;
        }
        .login-button {
            width: 100%;
        }
    }

    .login-card {
        margin: 20px;
        min-width:400px;
    }

  .login-title{
    text-align: center;
    margin: 20px;
  }

  .oauth-buttons {
    display: flex;
    flex-direction: column;
    gap: 10px;
    margin-top: 10px;
  }

  .oauth-button {
    min-width: 100px;
    width: 100%;
  }
  </style>
