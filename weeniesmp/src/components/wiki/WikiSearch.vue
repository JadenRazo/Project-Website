<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useWikiStore } from '@/stores/wikiStore'
import { Search, FileText, X } from 'lucide-vue-next'

const emit = defineEmits<{
  select: [slug: string]
}>()

const wikiStore = useWikiStore()

const inputRef = ref<HTMLInputElement | null>(null)
const isOpen = ref(false)
const selectedIndex = ref(0)
const searchInput = ref('')

// Debounce search
let debounceTimer: ReturnType<typeof setTimeout> | null = null

watch(searchInput, (value) => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
  debounceTimer = setTimeout(() => {
    wikiStore.setSearchQuery(value)
    selectedIndex.value = 0
    isOpen.value = value.trim().length > 0
  }, 200)
})

const results = computed(() => wikiStore.searchResults)

function handleKeydown(e: KeyboardEvent) {
  if (!isOpen.value || results.value.length === 0) return

  switch (e.key) {
    case 'ArrowDown':
      e.preventDefault()
      selectedIndex.value = Math.min(selectedIndex.value + 1, results.value.length - 1)
      break
    case 'ArrowUp':
      e.preventDefault()
      selectedIndex.value = Math.max(selectedIndex.value - 1, 0)
      break
    case 'Enter':
      e.preventDefault()
      if (results.value[selectedIndex.value]) {
        selectPage(results.value[selectedIndex.value].slug)
      }
      break
    case 'Escape':
      e.preventDefault()
      closeSearch()
      break
  }
}

function selectPage(slug: string) {
  emit('select', slug)
  closeSearch()
}

function closeSearch() {
  isOpen.value = false
  searchInput.value = ''
  wikiStore.setSearchQuery('')
}

function handleFocus() {
  if (searchInput.value.trim().length > 0) {
    isOpen.value = true
  }
}

function handleBlur() {
  // Delay close to allow click on results
  setTimeout(() => {
    if (!inputRef.value?.contains(document.activeElement)) {
      isOpen.value = false
    }
  }, 200)
}

// Global keyboard shortcut (Ctrl/Cmd + K)
function handleGlobalKeydown(e: KeyboardEvent) {
  if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
    e.preventDefault()
    inputRef.value?.focus()
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleGlobalKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleGlobalKeydown)
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
})

// Highlight matching text
function highlightMatch(text: string, query: string): string {
  if (!query) return text
  const regex = new RegExp(`(${query})`, 'gi')
  return text.replace(regex, '<mark class="bg-weenie-gold/30 text-white">$1</mark>')
}
</script>

<template>
  <div class="relative">
    <!-- Search Input -->
    <div class="relative">
      <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
      <input
        ref="inputRef"
        v-model="searchInput"
        type="text"
        placeholder="Search wiki... (Ctrl+K)"
        class="w-full pl-10 pr-10 py-3 bg-white/5 border border-white/10 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-weenie-gold/50 focus:ring-1 focus:ring-weenie-gold/50 transition-colors"
        @keydown="handleKeydown"
        @focus="handleFocus"
        @blur="handleBlur"
      />
      <button
        v-if="searchInput"
        @click="closeSearch"
        class="absolute right-3 top-1/2 -translate-y-1/2 p-1 hover:bg-white/10 rounded transition-colors"
      >
        <X class="w-4 h-4 text-gray-400" />
      </button>
    </div>

    <!-- Search Results Dropdown -->
    <Transition
      enter-active-class="transition-all duration-200 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition-all duration-150 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div
        v-if="isOpen && (results.length > 0 || searchInput.trim())"
        class="absolute z-50 w-full mt-2 bg-weenie-dark border border-white/10 rounded-lg shadow-xl overflow-hidden"
      >
        <!-- Results -->
        <ul v-if="results.length > 0" class="max-h-80 overflow-y-auto">
          <li v-for="(page, index) in results" :key="page.slug">
            <button
              @click="selectPage(page.slug)"
              @mouseenter="selectedIndex = index"
              class="w-full flex items-start gap-3 px-4 py-3 text-left transition-colors"
              :class="index === selectedIndex ? 'bg-white/10' : 'hover:bg-white/5'"
            >
              <FileText class="w-5 h-5 text-weenie-gold flex-shrink-0 mt-0.5" />
              <div class="min-w-0 flex-1">
                <div
                  class="font-medium text-white truncate"
                  v-html="highlightMatch(page.title, searchInput)"
                />
                <div class="text-sm text-gray-500 truncate">
                  {{ page.category }}
                </div>
              </div>
            </button>
          </li>
        </ul>

        <!-- No Results -->
        <div v-else-if="searchInput.trim()" class="px-4 py-6 text-center text-gray-400">
          No pages found for "{{ searchInput }}"
        </div>

        <!-- Keyboard Hints -->
        <div class="px-4 py-2 bg-black/30 border-t border-white/5 flex items-center gap-4 text-xs text-gray-500">
          <span class="flex items-center gap-1">
            <kbd class="px-1.5 py-0.5 bg-white/10 rounded">↑↓</kbd>
            Navigate
          </span>
          <span class="flex items-center gap-1">
            <kbd class="px-1.5 py-0.5 bg-white/10 rounded">Enter</kbd>
            Select
          </span>
          <span class="flex items-center gap-1">
            <kbd class="px-1.5 py-0.5 bg-white/10 rounded">Esc</kbd>
            Close
          </span>
        </div>
      </div>
    </Transition>
  </div>
</template>
