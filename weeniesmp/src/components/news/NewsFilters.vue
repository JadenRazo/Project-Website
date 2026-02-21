<script setup lang="ts">
import { ref, watch } from 'vue'
import { Search, X } from 'lucide-vue-next'
import { useNewsStore } from '@/stores/newsStore'

const props = defineProps<{
  activeCategory: string
  searchQuery: string
}>()

const emit = defineEmits<{
  (e: 'update:activeCategory', value: string): void
  (e: 'update:searchQuery', value: string): void
}>()

const newsStore = useNewsStore()
const localSearch = ref(props.searchQuery)
let debounceTimer: ReturnType<typeof setTimeout> | null = null

// Debounce search input
watch(localSearch, (newValue) => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
  debounceTimer = setTimeout(() => {
    emit('update:searchQuery', newValue)
  }, 300)
})

// Sync external changes to local search
watch(() => props.searchQuery, (newValue) => {
  if (newValue !== localSearch.value) {
    localSearch.value = newValue
  }
})

function clearSearch() {
  localSearch.value = ''
  emit('update:searchQuery', '')
}

function setCategory(categoryId: string) {
  emit('update:activeCategory', categoryId)
}

const categoryColors: Record<string, string> = {
  all: 'bg-weenie-gradient',
  update: 'bg-blue-500',
  event: 'bg-purple-500',
  patch: 'bg-green-500',
  announcement: 'bg-orange-500'
}
</script>

<template>
  <div class="space-y-6">
    <!-- Search Bar -->
    <div class="relative max-w-md mx-auto">
      <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-500" />
      <input
        v-model="localSearch"
        type="text"
        placeholder="Search news..."
        class="w-full pl-12 pr-12 py-3 bg-weenie-dark border border-gray-700 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:border-weenie-gold transition-colors"
      />
      <button
        v-if="localSearch"
        @click="clearSearch"
        class="absolute right-4 top-1/2 -translate-y-1/2 text-gray-500 hover:text-white transition-colors"
      >
        <X class="w-5 h-5" />
      </button>
    </div>

    <!-- Category Tabs -->
    <div class="flex flex-wrap justify-center gap-2">
      <button
        v-for="category in newsStore.categories"
        :key="category.id"
        @click="setCategory(category.id)"
        class="flex items-center gap-2 px-5 py-2.5 rounded-lg font-medium transition-all duration-300"
        :class="
          activeCategory === category.id
            ? `${categoryColors[category.id] || 'bg-weenie-gradient'} text-white shadow-lg`
            : 'bg-weenie-dark/50 text-gray-400 hover:text-white hover:bg-weenie-dark'
        "
      >
        {{ category.name }}
      </button>
    </div>
  </div>
</template>
