<script setup lang="ts">
import { ref, watch } from 'vue'
import { Search, SlidersHorizontal, X, ChevronDown } from 'lucide-vue-next'
import { useGalleryStore } from '@/stores/galleryStore'

const galleryStore = useGalleryStore()

const searchInput = ref(galleryStore.filters.search)
const showFilters = ref(false)
const sortDropdownOpen = ref(false)

const sortOptions = [
  { value: 'newest', label: 'Newest First' },
  { value: 'popular', label: 'Most Popular' },
  { value: 'random', label: 'Random' }
] as const

// Debounced search
let searchTimeout: ReturnType<typeof setTimeout>
watch(searchInput, (value) => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    galleryStore.setFilter('search', value)
  }, 300)
})

function toggleTag(tag: string) {
  galleryStore.toggleTag(tag)
}

function setSortBy(value: 'newest' | 'popular' | 'random') {
  galleryStore.setFilter('sortBy', value)
  sortDropdownOpen.value = false
}

function clearFilters() {
  galleryStore.clearFilters()
  searchInput.value = ''
}

const hasActiveFilters = () => {
  return galleryStore.filters.tags.length > 0 ||
    galleryStore.filters.search !== '' ||
    galleryStore.filters.sortBy !== 'newest'
}
</script>

<template>
  <div class="space-y-4">
    <!-- Search and Controls Bar -->
    <div class="flex flex-col sm:flex-row gap-3">
      <!-- Search Input -->
      <div class="relative flex-1">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" />
        <input
          v-model="searchInput"
          type="text"
          placeholder="Search images..."
          class="w-full pl-10 pr-4 py-2.5 bg-weenie-dark border border-white/10 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-weenie-gold/50 focus:ring-1 focus:ring-weenie-gold/50 transition-colors"
        />
        <button
          v-if="searchInput"
          @click="searchInput = ''; galleryStore.setFilter('search', '')"
          class="absolute right-3 top-1/2 -translate-y-1/2 p-1 text-gray-500 hover:text-white transition-colors"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Sort Dropdown -->
      <div class="relative">
        <button
          @click="sortDropdownOpen = !sortDropdownOpen"
          class="w-full sm:w-auto inline-flex items-center justify-between gap-2 px-4 py-2.5 bg-weenie-dark border border-white/10 rounded-lg text-gray-300 hover:text-white hover:border-white/20 transition-colors"
        >
          <span class="text-sm">
            {{ sortOptions.find(o => o.value === galleryStore.filters.sortBy)?.label }}
          </span>
          <ChevronDown class="w-4 h-4" :class="{ 'rotate-180': sortDropdownOpen }" />
        </button>

        <Transition
          enter-active-class="transition-all duration-200 ease-out"
          enter-from-class="opacity-0 -translate-y-2"
          enter-to-class="opacity-100 translate-y-0"
          leave-active-class="transition-all duration-150 ease-in"
          leave-from-class="opacity-100 translate-y-0"
          leave-to-class="opacity-0 -translate-y-2"
        >
          <div
            v-if="sortDropdownOpen"
            class="absolute right-0 mt-2 w-48 bg-weenie-dark border border-white/10 rounded-lg shadow-xl z-20 overflow-hidden"
          >
            <button
              v-for="option in sortOptions"
              :key="option.value"
              @click="setSortBy(option.value)"
              class="w-full px-4 py-2.5 text-left text-sm transition-colors"
              :class="galleryStore.filters.sortBy === option.value
                ? 'bg-weenie-gold/20 text-weenie-gold'
                : 'text-gray-300 hover:bg-white/5 hover:text-white'"
            >
              {{ option.label }}
            </button>
          </div>
        </Transition>
      </div>

      <!-- Filter Toggle (Mobile) -->
      <button
        @click="showFilters = !showFilters"
        class="sm:hidden inline-flex items-center justify-center gap-2 px-4 py-2.5 bg-weenie-dark border border-white/10 rounded-lg text-gray-300 hover:text-white transition-colors"
        :class="{ 'border-weenie-gold text-weenie-gold': galleryStore.filters.tags.length > 0 }"
      >
        <SlidersHorizontal class="w-4 h-4" />
        <span class="text-sm">Filters</span>
        <span
          v-if="galleryStore.filters.tags.length > 0"
          class="w-5 h-5 text-xs bg-weenie-gold text-black rounded-full flex items-center justify-center font-medium"
        >
          {{ galleryStore.filters.tags.length }}
        </span>
      </button>
    </div>

    <!-- Tag Filters -->
    <Transition
      enter-active-class="transition-all duration-300 ease-out"
      enter-from-class="max-h-0 opacity-0"
      enter-to-class="max-h-40 opacity-100"
      leave-active-class="transition-all duration-200 ease-in"
      leave-from-class="max-h-40 opacity-100"
      leave-to-class="max-h-0 opacity-0"
    >
      <div
        v-if="showFilters || true"
        class="overflow-hidden"
        :class="{ 'hidden sm:block': !showFilters }"
      >
        <div class="flex flex-wrap gap-2">
          <button
            v-for="tag in galleryStore.availableTags"
            :key="tag"
            @click="toggleTag(tag)"
            class="px-3 py-1.5 text-sm rounded-full border transition-all duration-200"
            :class="galleryStore.filters.tags.includes(tag)
              ? 'bg-weenie-gold text-black border-weenie-gold font-medium'
              : 'bg-transparent text-gray-400 border-white/10 hover:border-white/30 hover:text-white'"
          >
            {{ tag }}
          </button>

          <!-- Clear Filters Button -->
          <button
            v-if="hasActiveFilters()"
            @click="clearFilters"
            class="px-3 py-1.5 text-sm rounded-full border border-weenie-red/50 text-weenie-red hover:bg-weenie-red/10 transition-colors"
          >
            Clear All
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>
