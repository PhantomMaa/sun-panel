<script setup lang="ts">
import { NConfigProvider } from 'naive-ui'
import { NaiveProvider } from './components/common'
import { useTheme } from './hooks/useTheme'
import { useLanguage } from './hooks/useLanguage'
import { usePanelState } from './store'
import { watch } from 'vue'

const { theme, themeOverrides } = useTheme()
const { language } = useLanguage()
const panelState = usePanelState()

function updateFavicon(url: string | undefined) {
  let link = document.querySelector("link[rel~='icon']") as HTMLLinkElement
  if (!link) {
    link = document.createElement('link')
    link.rel = 'icon'
    document.head.appendChild(link)
  }
  link.href = url || '/favicon.svg'
}

watch(() => panelState.panelConfig.logoImageSrc, (newUrl) => {
  updateFavicon(newUrl)
})

// 首次加载时更新
updateFavicon(panelState.panelConfig.logoImageSrc)
</script>

<template>
  <NConfigProvider
    class="h-full"
    :theme="theme"
    :theme-overrides="themeOverrides"
    :locale="language"
  >
    <NaiveProvider>
      <RouterView />
    </NaiveProvider>
  </NConfigProvider>
</template>
