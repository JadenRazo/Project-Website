<script setup lang="ts">
import { ref, computed } from 'vue'
import { Search, X } from 'lucide-vue-next'

const searchQuery = ref('')

const emit = defineEmits<{
  (e: 'search', query: string): void
  (e: 'clear'): void
}>()

const hasQuery = computed(() => searchQuery.value.trim().length > 0)

function handleSearch() {
  const query = searchQuery.value.trim()
  if (query) {
    emit('search', query)
  }
}

function handleClear() {
  searchQuery.value = ''
  emit('clear')
}

function handleInput() {
  if (!searchQuery.value.trim()) {
    emit('clear')
  }
}
</script>

<template>
  <div class="mb-8">
    <div class="max-w-md mx-auto relative">
      <div class="relative">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400 pointer-events-none" />

        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search for a player..."
          class="w-full pl-10 pr-10 py-3 bg-weenie-dark border border-white/10 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-weenie-red/50 focus:ring-2 focus:ring-weenie-red/20 transition-all"
          @input="handleInput"
          @keyup.enter="handleSearch"
        />

        <button
          v-if="hasQuery"
          @click="handleClear"
          class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-white transition-colors"
          type="button"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <p class="text-gray-500 text-sm mt-2 text-center">
        Press Enter or type to filter players
      </p>
    </div>
  </div>
</template>
