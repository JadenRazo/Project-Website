<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { Newspaper, ChevronLeft, ChevronRight, Loader2 } from 'lucide-vue-next'
import NewsCard from '@/components/news/NewsCard.vue'
import NewsFilters from '@/components/news/NewsFilters.vue'
import { useNewsStore } from '@/stores/newsStore'

const newsStore = useNewsStore()

const activeCategory = ref('all')
const searchQuery = ref('')
const currentPage = ref(1)
const articlesPerPage = 6

async function loadArticles() {
  await newsStore.fetchArticles({
    category: activeCategory.value,
    search: searchQuery.value,
    page: currentPage.value,
    limit: articlesPerPage
  })
}

// Watch for filter changes
watch([activeCategory, searchQuery], () => {
  currentPage.value = 1
  loadArticles()
})

// Watch for page changes
watch(currentPage, () => {
  loadArticles()
  // Scroll to top of articles section
  window.scrollTo({ top: 200, behavior: 'smooth' })
})

function goToPage(page: number) {
  if (page >= 1 && page <= newsStore.totalPages) {
    currentPage.value = page
  }
}

function nextPage() {
  if (currentPage.value < newsStore.totalPages) {
    currentPage.value++
  }
}

function prevPage() {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

// Generate page numbers for pagination
function getPageNumbers(): (number | string)[] {
  const total = newsStore.totalPages
  const current = currentPage.value
  const pages: (number | string)[] = []

  if (total <= 7) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    pages.push(1)

    if (current > 3) {
      pages.push('...')
    }

    for (let i = Math.max(2, current - 1); i <= Math.min(total - 1, current + 1); i++) {
      if (!pages.includes(i)) {
        pages.push(i)
      }
    }

    if (current < total - 2) {
      pages.push('...')
    }

    if (!pages.includes(total)) {
      pages.push(total)
    }
  }

  return pages
}

onMounted(() => {
  loadArticles()
})
</script>

<template>
  <div class="min-h-screen pt-24 pb-16 bg-weenie-darker">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="text-center mb-12">
        <Newspaper class="w-12 h-12 text-weenie-gold mx-auto mb-4" />
        <h1 class="text-4xl md:text-5xl font-bold mb-4">
          <span class="gradient-text">News & Updates</span>
        </h1>
        <p class="text-gray-400 max-w-xl mx-auto">
          Stay up to date with the latest announcements, events, and patch notes.
        </p>
      </div>

      <!-- Filters -->
      <div class="mb-10">
        <NewsFilters
          v-model:active-category="activeCategory"
          v-model:search-query="searchQuery"
        />
      </div>

      <!-- Loading State -->
      <div v-if="newsStore.loading" class="flex flex-col items-center justify-center py-20">
        <Loader2 class="w-10 h-10 text-weenie-gold animate-spin mb-4" />
        <p class="text-gray-400">Loading articles...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="newsStore.error" class="text-center py-20">
        <p class="text-red-400 mb-4">{{ newsStore.error }}</p>
        <button
          @click="loadArticles"
          class="px-6 py-2 bg-weenie-red text-white rounded-lg hover:bg-weenie-red/80 transition-colors"
        >
          Retry
        </button>
      </div>

      <!-- Empty State -->
      <div v-else-if="newsStore.articles.length === 0" class="text-center py-20">
        <Newspaper class="w-16 h-16 text-gray-600 mx-auto mb-4" />
        <h3 class="text-xl font-semibold text-white mb-2">No articles found</h3>
        <p class="text-gray-400 mb-6">
          {{ searchQuery ? `No results for "${searchQuery}"` : 'No articles in this category yet.' }}
        </p>
        <button
          v-if="searchQuery || activeCategory !== 'all'"
          @click="searchQuery = ''; activeCategory = 'all'"
          class="px-6 py-2 text-weenie-gold border border-weenie-gold rounded-lg hover:bg-weenie-gold/10 transition-colors"
        >
          Clear filters
        </button>
      </div>

      <!-- Articles Grid -->
      <template v-else>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
          <NewsCard
            v-for="article in newsStore.articles"
            :key="article.id"
            :article="article"
          />
        </div>

        <!-- Pagination -->
        <div v-if="newsStore.totalPages > 1" class="flex items-center justify-center gap-2">
          <!-- Previous Button -->
          <button
            @click="prevPage"
            :disabled="currentPage === 1"
            class="p-2 rounded-lg transition-all"
            :class="currentPage === 1
              ? 'text-gray-600 cursor-not-allowed'
              : 'text-gray-400 hover:text-white hover:bg-weenie-dark'"
          >
            <ChevronLeft class="w-5 h-5" />
          </button>

          <!-- Page Numbers -->
          <div class="flex items-center gap-1">
            <template v-for="(page, index) in getPageNumbers()" :key="index">
              <span v-if="page === '...'" class="px-3 py-2 text-gray-500">...</span>
              <button
                v-else
                @click="goToPage(page as number)"
                class="min-w-[40px] h-10 rounded-lg font-medium transition-all"
                :class="currentPage === page
                  ? 'bg-weenie-gradient text-white'
                  : 'text-gray-400 hover:text-white hover:bg-weenie-dark'"
              >
                {{ page }}
              </button>
            </template>
          </div>

          <!-- Next Button -->
          <button
            @click="nextPage"
            :disabled="currentPage === newsStore.totalPages"
            class="p-2 rounded-lg transition-all"
            :class="currentPage === newsStore.totalPages
              ? 'text-gray-600 cursor-not-allowed'
              : 'text-gray-400 hover:text-white hover:bg-weenie-dark'"
          >
            <ChevronRight class="w-5 h-5" />
          </button>
        </div>

        <!-- Results Info -->
        <div class="text-center mt-6 text-sm text-gray-500">
          Showing {{ newsStore.articles.length }} of {{ newsStore.totalArticles }} articles
        </div>
      </template>
    </div>
  </div>
</template>
