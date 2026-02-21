<script setup lang="ts">
import { computed, ref } from 'vue'

interface Props {
  uuid: string
  username?: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
}

const props = withDefaults(defineProps<Props>(), {
  size: 'lg',
  username: ''
})

const imageLoaded = ref(false)
const imageError = ref(false)

// Remove dashes from UUID for Crafatar
const cleanUuid = computed(() => props.uuid.replace(/-/g, ''))

// Size mapping for the bust image
const sizeMap = {
  sm: { width: 64, class: 'w-16 h-16' },
  md: { width: 128, class: 'w-32 h-32' },
  lg: { width: 200, class: 'w-48 h-48' },
  xl: { width: 300, class: 'w-72 h-72' }
}

const sizeConfig = computed(() => sizeMap[props.size])

// Crafatar renders - using body render for full character view
const skinUrl = computed(() => {
  const scale = sizeConfig.value.width > 128 ? 10 : 6
  return `https://crafatar.com/renders/body/${cleanUuid.value}?scale=${scale}&overlay`
})

// Fallback to bust if body fails
const bustUrl = computed(() => {
  return `https://crafatar.com/renders/head/${cleanUuid.value}?scale=6&overlay`
})

// Alternative: NameMC 3D render iframe (commented out as fallback option)
// const nameMcUrl = computed(() => {
//   return `https://namemc.com/skin/${cleanUuid.value}`
// })

function handleImageError() {
  imageError.value = true
}

function handleImageLoad() {
  imageLoaded.value = true
}
</script>

<template>
  <div class="relative flex items-center justify-center">
    <!-- Loading skeleton -->
    <div
      v-if="!imageLoaded && !imageError"
      class="absolute inset-0 flex items-center justify-center"
    >
      <div
        class="animate-pulse bg-white/5 rounded-lg"
        :class="sizeConfig.class"
      ></div>
    </div>

    <!-- Skin render -->
    <img
      v-if="!imageError"
      :src="skinUrl"
      :alt="`${username || 'Player'}'s skin`"
      class="transition-opacity duration-300 drop-shadow-2xl"
      :class="[
        sizeConfig.class,
        imageLoaded ? 'opacity-100' : 'opacity-0'
      ]"
      @load="handleImageLoad"
      @error="handleImageError"
      loading="lazy"
    />

    <!-- Fallback to bust render on error -->
    <img
      v-else
      :src="bustUrl"
      :alt="`${username || 'Player'}'s skin`"
      class="drop-shadow-2xl"
      :class="sizeConfig.class"
      loading="lazy"
    />
  </div>
</template>
