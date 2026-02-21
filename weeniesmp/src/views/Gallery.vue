<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Image, Loader2, Filter, Upload } from 'lucide-vue-next'
import GalleryGrid from '@/components/gallery/GalleryGrid.vue'
import GalleryFilters from '@/components/gallery/GalleryFilters.vue'
import GalleryLightbox from '@/components/gallery/GalleryLightbox.vue'
import { useGalleryStore, type GalleryImage } from '@/stores/galleryStore'
import { useUserStore } from '@/stores/userStore'

const galleryStore = useGalleryStore()
const userStore = useUserStore()

const selectedImage = ref<GalleryImage | null>(null)
const lightboxOpen = ref(false)
const showFilters = ref(false)

const isLoggedIn = computed(() => !!userStore.user)

function openLightbox(image: GalleryImage) {
  selectedImage.value = image
  lightboxOpen.value = true
}

function closeLightbox() {
  lightboxOpen.value = false
  selectedImage.value = null
}

function navigateImage(direction: 'prev' | 'next') {
  if (!selectedImage.value) return

  const currentIndex = galleryStore.filteredImages.findIndex(img => img.id === selectedImage.value?.id)
  let newIndex: number

  if (direction === 'prev') {
    newIndex = currentIndex > 0 ? currentIndex - 1 : galleryStore.filteredImages.length - 1
  } else {
    newIndex = currentIndex < galleryStore.filteredImages.length - 1 ? currentIndex + 1 : 0
  }

  selectedImage.value = galleryStore.filteredImages[newIndex]
}

async function loadImages() {
  await galleryStore.fetchImages(1)
}

function handleFilterChange(tags: string[], sort: 'newest' | 'popular' | 'random') {
  galleryStore.filters.tags = tags
  galleryStore.filters.sortBy = sort
}

onMounted(() => {
  loadImages()
})
</script>

<template>
  <div class="min-h-screen pt-24 pb-16 bg-weenie-darker">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="text-center mb-12">
        <Image class="w-12 h-12 text-weenie-gold mx-auto mb-4" />
        <h1 class="text-4xl md:text-5xl font-bold mb-4">
          <span class="gradient-text">Community Gallery</span>
        </h1>
        <p class="text-gray-400 max-w-xl mx-auto">
          Explore amazing builds, screenshots, and moments from our community.
        </p>
      </div>

      <!-- Actions Bar -->
      <div class="flex items-center justify-between mb-8">
        <button
          @click="showFilters = !showFilters"
          class="flex items-center gap-2 px-4 py-2 bg-weenie-dark text-gray-400 hover:text-white rounded-lg transition-colors"
        >
          <Filter class="w-5 h-5" />
          Filters
        </button>

        <button
          v-if="isLoggedIn"
          class="flex items-center gap-2 px-4 py-2 bg-weenie-gradient text-white rounded-lg hover:opacity-90 transition-opacity"
        >
          <Upload class="w-5 h-5" />
          Upload
        </button>
        <RouterLink
          v-else
          to="/dashboard"
          class="text-sm text-gray-400 hover:text-weenie-gold transition-colors"
        >
          Log in to upload
        </RouterLink>
      </div>

      <!-- Filters -->
      <Transition
        enter-active-class="transition duration-200 ease-out"
        enter-from-class="opacity-0 -translate-y-2"
        enter-to-class="opacity-100 translate-y-0"
        leave-active-class="transition duration-150 ease-in"
        leave-from-class="opacity-100 translate-y-0"
        leave-to-class="opacity-0 -translate-y-2"
      >
        <div v-if="showFilters" class="mb-8">
          <GalleryFilters
            :active-tags="galleryStore.filters.tags"
            :sort-by="galleryStore.filters.sortBy"
            :available-tags="galleryStore.availableTags"
            @change="handleFilterChange"
          />
        </div>
      </Transition>

      <!-- Loading State -->
      <div v-if="galleryStore.loading" class="flex flex-col items-center justify-center py-20">
        <Loader2 class="w-10 h-10 text-weenie-gold animate-spin mb-4" />
        <p class="text-gray-400">Loading gallery...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="galleryStore.error" class="text-center py-20">
        <p class="text-red-400 mb-4">{{ galleryStore.error }}</p>
        <button
          @click="loadImages"
          class="px-6 py-2 bg-weenie-red text-white rounded-lg hover:bg-weenie-red/80 transition-colors"
        >
          Retry
        </button>
      </div>

      <!-- Empty State -->
      <div v-else-if="galleryStore.filteredImages.length === 0" class="text-center py-20">
        <Image class="w-16 h-16 text-gray-600 mx-auto mb-4" />
        <h3 class="text-xl font-semibold text-white mb-2">No images found</h3>
        <p class="text-gray-400 mb-6">
          {{ galleryStore.filters.tags.length > 0 ? 'No images with these tags yet.' : 'Be the first to share your screenshots!' }}
        </p>
        <button
          v-if="galleryStore.filters.tags.length > 0"
          @click="galleryStore.clearFilters()"
          class="px-6 py-2 text-weenie-gold border border-weenie-gold rounded-lg hover:bg-weenie-gold/10 transition-colors"
        >
          Clear filters
        </button>
      </div>

      <!-- Gallery Grid -->
      <template v-else>
        <GalleryGrid
          :images="galleryStore.filteredImages"
          @select="openLightbox"
        />

        <!-- Load More -->
        <div v-if="galleryStore.hasMore" class="text-center mt-12">
          <button
            @click="galleryStore.loadMore"
            :disabled="galleryStore.loadingMore"
            class="px-8 py-3 bg-weenie-dark text-white rounded-lg hover:bg-weenie-dark/80 transition-colors disabled:opacity-50"
          >
            <Loader2 v-if="galleryStore.loadingMore" class="w-5 h-5 animate-spin inline mr-2" />
            {{ galleryStore.loadingMore ? 'Loading...' : 'Load More' }}
          </button>
        </div>
      </template>
    </div>

    <!-- Lightbox -->
    <GalleryLightbox
      :image="selectedImage"
      :open="lightboxOpen"
      @close="closeLightbox"
      @prev="navigateImage('prev')"
      @next="navigateImage('next')"
    />
  </div>
</template>
