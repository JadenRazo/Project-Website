<script setup lang="ts">
import { ref, computed } from 'vue'
import { Heart, User } from 'lucide-vue-next'
import type { GalleryImage } from '@/stores/galleryStore'
import { useGalleryStore } from '@/stores/galleryStore'

interface Props {
  image: GalleryImage
}

const props = defineProps<Props>()
const emit = defineEmits<{
  click: [image: GalleryImage]
}>()

const galleryStore = useGalleryStore()
const imageLoaded = ref(false)
const imageError = ref(false)

const isLiked = computed(() => galleryStore.isLiked(props.image.id))

function onImageLoad() {
  imageLoaded.value = true
}

function onImageError() {
  imageError.value = true
  imageLoaded.value = true
}

function handleClick() {
  emit('click', props.image)
}

function handleLike(e: Event) {
  e.stopPropagation()
  galleryStore.likeImage(props.image.id)
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}
</script>

<template>
  <div
    class="group relative rounded-xl overflow-hidden bg-weenie-dark border border-white/10 cursor-pointer transition-all duration-300 hover:border-weenie-gold/50 hover:shadow-lg hover:shadow-weenie-gold/10"
    @click="handleClick"
  >
    <!-- Image Container -->
    <div class="relative aspect-[4/3] overflow-hidden bg-weenie-darker">
      <!-- Loading Skeleton -->
      <div
        v-if="!imageLoaded"
        class="absolute inset-0 bg-weenie-dark animate-pulse"
      />

      <!-- Error State -->
      <div
        v-if="imageError"
        class="absolute inset-0 flex items-center justify-center bg-weenie-dark"
      >
        <p class="text-gray-500 text-sm">Failed to load</p>
      </div>

      <!-- Image -->
      <img
        :src="image.thumbnail"
        :alt="image.title"
        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
        :class="{ 'opacity-0': !imageLoaded || imageError }"
        loading="lazy"
        @load="onImageLoad"
        @error="onImageError"
      />

      <!-- Hover Overlay -->
      <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />

      <!-- Hover Content -->
      <div class="absolute inset-0 flex flex-col justify-end p-4 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
        <h3 class="text-white font-semibold text-lg line-clamp-1 mb-1">
          {{ image.title }}
        </h3>
        <div class="flex items-center gap-2 text-sm text-gray-300">
          <User class="w-3.5 h-3.5" />
          <span>{{ image.author }}</span>
          <span class="text-gray-500">{{ formatDate(image.createdAt) }}</span>
        </div>
      </div>

      <!-- Tags (always visible) -->
      <div class="absolute top-3 left-3 flex flex-wrap gap-1.5">
        <span
          v-for="tag in image.tags.slice(0, 2)"
          :key="tag"
          class="px-2 py-0.5 text-xs font-medium bg-black/60 backdrop-blur-sm text-white rounded-full"
        >
          {{ tag }}
        </span>
      </div>

      <!-- Like Button -->
      <button
        @click="handleLike"
        class="absolute top-3 right-3 p-2 rounded-full bg-black/60 backdrop-blur-sm text-white opacity-0 group-hover:opacity-100 transition-all duration-300 hover:bg-weenie-red hover:scale-110"
        :class="{ '!opacity-100 bg-weenie-red': isLiked }"
      >
        <Heart class="w-4 h-4" :class="{ 'fill-current': isLiked }" />
      </button>
    </div>

    <!-- Info Bar -->
    <div class="p-3 flex items-center justify-between border-t border-white/5">
      <p class="text-sm text-gray-400 truncate max-w-[70%]">
        {{ image.title }}
      </p>
      <div class="flex items-center gap-1.5 text-sm text-gray-500">
        <Heart class="w-3.5 h-3.5" :class="{ 'text-weenie-red fill-weenie-red': isLiked }" />
        <span>{{ image.likes }}</span>
      </div>
    </div>
  </div>
</template>
