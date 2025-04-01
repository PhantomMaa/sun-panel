<script setup lang="ts">
import { computed, ref } from 'vue'
import { NEllipsis } from 'naive-ui'
import { ItemIcon } from '../../../../components/common'
import { PanelPanelConfigStyleEnum } from '@/enums'

interface Prop {
  itemInfo?: Panel.ItemInfo
  size?: number // 默认70
  forceBackground?: string // 强制背景色
  iconTextColor?: string
  iconTextInfoHideDescription: boolean
  iconTextIconHideTitle: boolean
  style: PanelPanelConfigStyleEnum
}

const props = withDefaults(defineProps<Prop>(), {
  size: 70,
})

const defaultBackground = '#2a2a2a6b'

const calculateLuminance = (color: string) => {
  const hex = color.replace(/^#/, '')
  const r = parseInt(hex.substring(0, 2), 16)
  const g = parseInt(hex.substring(2, 4), 16)
  const b = parseInt(hex.substring(4, 6), 16)
  return (0.299 * r + 0.587 * g + 0.114 * b) / 255
}

const textColor = computed(() => {
  const luminance = calculateLuminance(props.itemInfo?.icon?.backgroundColor || defaultBackground)
  return luminance > 0.5 ? 'black' : 'white'
})

// Card tilt effect variables
const isHovering = ref(false)
const cardTransform = ref({ x: 0, y: 0, scale: 1 })

// Handle mouse events
const handleMouseEnter = () => {
  isHovering.value = true
}

const handleMouseLeave = () => {
  isHovering.value = false
  // Reset transform on mouse leave
  cardTransform.value = { x: 0, y: 0, scale: 1 }
}

const handleMouseMove = (e: MouseEvent, element: EventTarget | null) => {
  if (!(element instanceof HTMLElement) || !isHovering.value) return
  
  const rect = element.getBoundingClientRect()
  const centerX = rect.left + rect.width / 2
  const centerY = rect.top + rect.height / 2
  
  // Calculate distance from center (normalized to -1 to 1)
  const x = (e.clientX - centerX) / (rect.width / 2)
  const y = (e.clientY - centerY) / (rect.height / 2)
  
  // Update transform values (limit tilt to 10 degrees)
  cardTransform.value = {
    x: y * -10, // Invert Y axis for natural tilt
    y: x * 10,  // X axis tilt
    scale: 1.05 // Slight scale up on hover
  }
  
  // Update glow position
  if (element) {
    // Calculate relative position for the glow effect (0-100%)
    const relativeX = ((e.clientX - rect.left) / rect.width) * 100
    const relativeY = ((e.clientY - rect.top) / rect.height) * 100
    element.style.setProperty('--x', `${relativeX}%`)
    element.style.setProperty('--y', `${relativeY}%`)
  }
}

// Computed styles for card transform
const cardStyle = computed(() => {
  if (!isHovering.value) {
    return {
      transform: 'perspective(1000px) rotateX(0deg) rotateY(0deg) scale(1)',
      transition: 'all 0.5s ease-out'
    }
  }
  
  return {
    transform: `perspective(1000px) rotateX(${cardTransform.value.x}deg) rotateY(${cardTransform.value.y}deg) scale(${cardTransform.value.scale})`,
    transition: 'transform 0.1s ease-out'
  }
})
</script>

<template>
  <div class="app-icon w-full">
    <!-- 详情图标 -->
    <div
      v-if="style === PanelPanelConfigStyleEnum.info"
      class="app-icon-info w-full rounded-2xl transition-all duration-200 hover:shadow-[0_0_25px_rgba(0,0,0,0.3)] flex card-container"
      :style="[
        { background: itemInfo?.icon?.backgroundColor || defaultBackground },
        cardStyle
      ]"
      @mouseenter="handleMouseEnter"
      @mouseleave="handleMouseLeave"
      @mousemove="(e) => handleMouseMove(e, e.currentTarget)"
    >
      <!-- 图标 -->
      <div class="app-icon-info-icon w-[70px] h-[70px]">
        <div class="w-[70px] h-full flex items-center justify-center">
          <ItemIcon :item-icon="itemInfo?.icon" force-background="transparent" :size="50" class="overflow-hidden rounded-xl" />
        </div>
      </div>

      <!-- 文字 -->
      <!-- 如果为纯白色，将自动根据背景的明暗计算字体的黑白色 -->
      <div class="text-white flex items-center" :style="{ color: (iconTextColor === '#ffffff') ? textColor : iconTextColor, maxWidth: 'calc(100% - 80px)' }">
        <div class="app-icon-info-text-box w-full">
          <div class="app-icon-info-text-box-title font-semibold w-full">
            <NEllipsis>
              {{ itemInfo?.title }}
            </NEllipsis>
          </div>
          <div v-if="!iconTextInfoHideDescription" class="app-icon-info-text-box-description">
            <NEllipsis :line-clamp="2" class="text-xs">
              {{ itemInfo?.description }}
            </NEllipsis>
          </div>
        </div>
      </div>
      
      <!-- Hover glow effect -->
      <div class="card-glow"></div>
    </div>

    <!-- 极简(小)图标（APP） -->
    <div v-if="style === PanelPanelConfigStyleEnum.icon" class="app-icon-small">
      <div
        class="app-icon-small-icon overflow-hidden rounded-2xl sunpanel w-[70px] h-[70px] mx-auto transition-all duration-200 hover:shadow-[0_0_25px_rgba(0,0,0,0.3)] card-container"
        :title="itemInfo?.description"
        :style="cardStyle"
        @mouseenter="handleMouseEnter"
        @mouseleave="handleMouseLeave"
        @mousemove="(e) => handleMouseMove(e, e.currentTarget)"
      >
        <ItemIcon :item-icon="itemInfo?.icon" />
        <!-- Hover glow effect -->
        <div class="card-glow"></div>
      </div>
      <div
        v-if="!iconTextIconHideTitle"
        class="app-icon-small-title text-center app-icon-text-shadow cursor-pointer mt-[2px]"
        :style="{ color: iconTextColor }"
      >
        <span>{{ itemInfo?.title }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card-container {
  position: relative;
  transform-style: preserve-3d;
  will-change: transform;
  overflow: hidden;
  border: 1px solid transparent;
  backface-visibility: hidden;
}

.card-glow {
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  background: radial-gradient(circle at var(--x, 50%) var(--y, 50%), rgba(255, 255, 255, 0.2) 0%, rgba(255, 255, 255, 0) 60%);
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.3s ease;
}

.card-container:hover .card-glow {
  opacity: 1;
}

.app-icon-info:hover, .app-icon-small-icon:hover {
  box-shadow: 0 20px 30px -10px rgba(0, 0, 0, 0.3);
  transform: translateY(-5px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}
</style>
