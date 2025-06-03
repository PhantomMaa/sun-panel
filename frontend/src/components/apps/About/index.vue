<script setup lang="ts">
import { NDivider, NTag } from 'naive-ui'
import { onMounted, ref } from 'vue'
// 使用 public 目录下的资源，直接使用绝对路径
// 在 Vite 中，public 目录下的文件会被原样复制到构建输出目录，不会经过任何处理
const favicon = '/assets/favicon.svg'
import { getAbout } from '@/api/system/about'

interface Version {
  versionName: string
}

const versionName = ref('')
const frontVersion = import.meta.env.VITE_APP_VERSION || 'unknown'

onMounted(() => {
  getAbout<Version>().then((res) => {
    if (res.code === 0)
      versionName.value = res.data.versionName
  })
})
</script>

<template>
  <div class="pt-5">
    <div class="flex flex-col items-center justify-center">
      <img :src="favicon" width="100" height="100" alt="">
      <div class="text-3xl font-semibold">
        {{ $t('common.appName') }}
      </div>
    </div>
    <NDivider style="margin:10px 0">•</NDivider>
    <div class="flex mt-[10px] flex-wrap justify-center">
      <div class="flex items-center mx-[10px]">
        <img class="w-[20px] h-[20px] mr-[5px]" :src="favicon" alt="">
        <a href="https://moon.phantomlab.top/site/" target="_blank" class="link">Moon-Box Site</a>
      </div>
    </div>
    <div class="flex flex-col items-center justify-center text-base">
      <div class="mt-5">
        <NTag :bordered="false" size="small">
          {{ versionName }} & FV-{{ frontVersion }}
        </NTag>
      </div>
    </div>
  </div>
</template>

<style>
.link{
    color:rgb(0, 89, 255)
}
</style>
