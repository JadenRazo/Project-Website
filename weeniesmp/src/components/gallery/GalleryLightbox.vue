<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { X, ChevronLeft, ChevronRight, Heart, User, Calendar, Tag, Download, ExternalLink, Loader2 } from 'lucide-vue-next'
import type { GalleryImage } from '@/stores/galleryStore'
import { useGalleryStore } from '@/stores/galleryStore'

interface Props {
  image: GalleryImage | null
  open: boolean
  images?: GalleryImage[]
}

const props = withDefaults(defineProps<Props>(), {
  images: () => []
})

const emit = defineEmits<{
  close: []
  'update:image': [image: GalleryImage | null]
}>()

const galleryStore = useGalleryStore()
const imageLoading = ref(true)
const lightboxRef = ref<HTMLElement | null>(null)

const isLiked = computed(() => props.image ? galleryStore.isLiked(props.image.id) : false)

const currentIndex = computed(() => {
  if (!props.image || props.images.length === 0) return -1
  return props.images.findIndex(img => img.id === props.image!.id)
})

const hasPrevious = computed(() => currentIndex.value > 0)
const hasNext = computed(() => currentIndex.value < props.images.length - 1)

function close() {
  emit('close')
}

function goToPrevious() {
  if (hasPrevious.value) {
    imageLoading.value = true
    emit('update:image', props.images[currentIndex.value - 1])
  }
}

function goToNext() {
  if (hasNext.value) {
    imageLoading.value = true
    emit('update:image', props.images[currentIndex.value + 1])
  }
}

function handleLike() {
  if (props.image) {
    galleryStore.likeImage(props.image.id)
  }
}

function onImageLoad() {
  imageLoading.value = false
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

function downloadImage() {
  if (!props.image) return
  const link = document.createElement('a')
  link.href = props.image.url
  link.download = `${props.image.title.replace(/\s+/g, '-').toLowerCase()}.jpg`
  link.target = '_blank'
  link.rel = 'noopener noreferrer'
  link.click()
}

function openFullSize() {
  if (props.image) {
    window.open(props.image.url, '_blank', 'noopener,noreferrer')
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (!props.open) return

  switch (e.key) {
    case 'Escape':
      close()
      break
    case 'ArrowLeft':
      e.preventDefault()
      goToPrevious()
      break
    case 'ArrowRight':
      e.preventDefault()
      goToNext()
      break
  }
}

function handleBackdropClick(e: MouseEvent) {
  if (e.target === lightboxRef.value) {
    close()
  }
}

// Reset loading state when image changes
watch(() => props.image, () => {
  imageLoading.value = true
})

// Lock body scroll when open
watch(() => props.open, (isOpen) => {
  if (isOpen) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
})

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
  document.body.style.overflow = ''
})
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-opacity duration-300"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-200"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="open && image"
        ref="lightboxRef"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/95 backdrop-blur-sm"
        @click="handleBackdropClick"
      >
        <!-- Close Button -->
        <button
          @click="close"
          class="absolute top-4 right-4 p-2 text-gray-400 hover:text-white bg-white/5 hover:bg-white/10 rounded-full transition-colors z-10"
          aria-label="Close lightbox"
        >
          <X class="w-6 h-6" />
        </button>

        <!-- Navigation Buttons -->
        <button
          v-if="hasPrevious"
          @click="goToPrevious"
          class="absolute left-4 top-1/2 -translate-y-1/2 p-3 text-gray-400 hover:text-white bg-white/5 hover:bg-white/10 rounded-full transition-colors z-10"
          aria-label="Previous image"
        >
          <ChevronLeft class="w-6 h-6" />
        </button>

        <button
          v-if="hasNext"
          @click="goToNext"
          class="absolute right-4 top-1/2 -translate-y-1/2 p-3 text-gray-400 hover:text-white bg-white/5 hover:bg-white/10 rounded-full transition-colors z-10"
          aria-label="Next image"
        >
          <ChevronRight class="w-6 h-6" />
        </button>

        <!-- Content Container -->
        <div class="w-full max-w-6xl mx-4 flex flex-col lg:flex-row gap-6 max-h-[90vh]">
          <!-- Image Container -->
          <div class="flex-1 relative flex items-center justify-center min-h-0">
            <!-- Loading Spinner -->
            <div
              v-if="imageLoading"
              class="absolute inset-0 flex items-center justify-center"
            >
              <Loader2 class="w-10 h-10 text-weenie-gold animate-spin" />
            </div>

            <!-- Image -->
            <img
              :src="image.url"
              :alt="image.title"
              class="max-w-full max-h-[70vh] lg:max-h-[85vh] object-contain rounded-lg shadow-2xl transition-opacity duration-300"
              :class="{ 'opacity-0': imageLoading }"
              @load="onImageLoad"
            />
          </div>

          <!-- Info Panel -->
          <div class="lg:w-80 flex-shrink-0 bg-weenie-dark/80 backdrop-blur-sm rounded-xl p-6 overflow-y-auto">
            <!-- Title and Author -->
            <div class="mb-6">
              <h2 class="text-xl font-bold text-white mb-2">{{ image.title }}</h2>
              <div class="flex items-center gap-2 text-gray-400">
                <User class="w-4 h-4" />
                <span>{{ image.author }}</span>
              </div>
            </div>

            <!-- Description -->
            <p class="text-gray-400 text-sm mb-6 leading-relaxed">
              {{ image.description }}
            </p>

            <!-- Meta Info -->
            <div class="space-y-3 mb-6">
              <div class="flex items-center gap-2 text-sm text-gray-400">
                <Calendar class="w-4 h-4" />
                <span>{{ formatDate(image.createdAt) }}</span>
              </div>
              <div class="flex items-center gap-2 text-sm">
                <Tag class="w-4 h-4 text-gray-400" />
                <div class="flex flex-wrap gap-1.5">
                  <span
                    v-for="tag in image.tags"
                    :key="tag"
                    class="px-2 py-0.5 text-xs bg-white/10 text-gray-300 rounded-full"
                  >
                    {{ tag }}
                  </span>
                </div>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex flex-wrap gap-2">
              <button
                @click="handleLike"
                class="flex-1 inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-lg font-medium transition-colors"
                :class="isLiked
                  ? 'bg-weenie-red text-white'
                  : 'bg-white/5 text-gray-300 hover:bg-white/10 hover:text-white'"
              >
                <Heart class="w-4 h-4" :class="{ 'fill-current': isLiked }" />
                <span>{{ image.likes }}</span>
              </button>

              <button
                @click="downloadImage"
                class="p-2.5 bg-white/5 text-gray-300 hover:bg-white/10 hover:text-white rounded-lg transition-colors"
                title="Download image"
              >
                <Download class="w-4 h-4" />
              </button>

              <button
                @click="openFullSize"
                class="p-2.5 bg-white/5 text-gray-300 hover:bg-white/10 hover:text-white rounded-lg transition-colors"
                title="Open full size"
              >
                <ExternalLink class="w-4 h-4" />
              </button>
            </div>

            <!-- Image Counter -->
            <div v-if="images.length > 1" class="mt-6 pt-4 border-t border-white/10">
              <p class="text-center text-sm text-gray-500">
                {{ currentIndex + 1 }} of {{ images.length }}
              </p>
            </div>
          </div>
        </div>

        <!-- Keyboard Hint -->
        <div class="absolute bottom-4 left-1/2 -translate-x-1/2 hidden lg:flex items-center gap-4 text-xs text-gray-600">
          <span class="flex items-center gap-1">
            <kbd class="px-2 py-1 bg-white/5 rounded">ESC</kbd> Close
          </span>
          <span class="flex items-center gap-1">
            <kbd class="px-2 py-1 bg-white/5 rounded">&larr;</kbd>
            <kbd class="px-2 py-1 bg-white/5 rounded">&rarr;</kbd> Navigate
          </span>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
