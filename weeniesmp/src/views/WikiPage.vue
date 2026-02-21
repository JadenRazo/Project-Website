<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { useWikiStore } from '@/stores/wikiStore'
import WikiSidebar from '@/components/wiki/WikiSidebar.vue'
import TableOfContents from '@/components/wiki/TableOfContents.vue'
import { ChevronLeft, ChevronRight, Calendar, Menu, X, Home } from 'lucide-vue-next'

const route = useRoute()
const wikiStore = useWikiStore()

const sidebarOpen = ref(false)
const slug = computed(() => route.params.slug as string)

// Configure marked
marked.setOptions({
  gfm: true,
  breaks: true
})

// Fetch page on mount and route change
onMounted(() => {
  wikiStore.fetchCategories()
  if (slug.value) {
    wikiStore.fetchPage(slug.value)
  }
})

watch(slug, (newSlug) => {
  if (newSlug) {
    wikiStore.fetchPage(newSlug)
    sidebarOpen.value = false
  }
})

// Computed
const renderedContent = computed(() => {
  if (!wikiStore.currentPage?.content) return ''
  const html = marked(wikiStore.currentPage.content) as string
  return DOMPurify.sanitize(html)
})

const adjacentPages = computed(() => {
  return wikiStore.getAdjacentPages(slug.value)
})

const formattedDate = computed(() => {
  if (!wikiStore.currentPage?.lastUpdated) return ''
  const date = new Date(wikiStore.currentPage.lastUpdated)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
})

const categoryName = computed(() => {
  const category = wikiStore.categories.find(c => c.id === wikiStore.currentPage?.category)
  return category?.name ?? ''
})

// Close sidebar on escape
function handleEscape(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    sidebarOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleEscape)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleEscape)
})
</script>

<template>
  <div class="min-h-screen pt-16 bg-weenie-darker">
    <!-- Mobile Sidebar Toggle -->
    <button
      @click="sidebarOpen = !sidebarOpen"
      class="lg:hidden fixed bottom-4 right-4 z-50 p-4 bg-weenie-dark border border-white/10 rounded-full shadow-lg hover:bg-white/5 transition-colors"
    >
      <Menu v-if="!sidebarOpen" class="w-6 h-6 text-white" />
      <X v-else class="w-6 h-6 text-white" />
    </button>

    <div class="flex">
      <!-- Sidebar Overlay (Mobile) -->
      <Transition
        enter-active-class="transition-opacity duration-300"
        enter-from-class="opacity-0"
        enter-to-class="opacity-100"
        leave-active-class="transition-opacity duration-200"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <div
          v-if="sidebarOpen"
          class="lg:hidden fixed inset-0 bg-black/50 z-40"
          @click="sidebarOpen = false"
        />
      </Transition>

      <!-- Sidebar -->
      <aside
        :class="[
          'fixed lg:sticky top-16 left-0 z-40 h-[calc(100vh-4rem)] w-72 bg-weenie-dark border-r border-white/10 overflow-y-auto transition-transform duration-300 lg:translate-x-0',
          sidebarOpen ? 'translate-x-0' : '-translate-x-full'
        ]"
      >
        <WikiSidebar
          :categories="wikiStore.categories"
          :current-slug="slug"
          @close="sidebarOpen = false"
        />
      </aside>

      <!-- Main Content -->
      <main class="flex-1 min-w-0 lg:ml-0">
        <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <!-- Breadcrumb -->
          <nav class="flex items-center gap-2 text-sm text-gray-400 mb-6">
            <router-link to="/wiki" class="hover:text-white transition-colors flex items-center gap-1">
              <Home class="w-4 h-4" />
              Wiki
            </router-link>
            <ChevronRight class="w-4 h-4" />
            <span v-if="categoryName" class="text-gray-500">{{ categoryName }}</span>
            <ChevronRight v-if="categoryName" class="w-4 h-4" />
            <span class="text-white">{{ wikiStore.currentPage?.title }}</span>
          </nav>

          <!-- Loading State -->
          <div v-if="wikiStore.loading" class="animate-pulse">
            <div class="h-10 bg-white/5 rounded w-1/2 mb-4"></div>
            <div class="h-4 bg-white/5 rounded w-full mb-2"></div>
            <div class="h-4 bg-white/5 rounded w-3/4 mb-2"></div>
            <div class="h-4 bg-white/5 rounded w-5/6"></div>
          </div>

          <!-- Error State -->
          <div v-else-if="wikiStore.error" class="card p-8 text-center">
            <p class="text-red-400 mb-4">{{ wikiStore.error }}</p>
            <router-link
              to="/wiki"
              class="inline-flex items-center gap-2 px-4 py-2 bg-weenie-gold text-black font-medium rounded-lg hover:bg-weenie-gold/90 transition-colors"
            >
              <ChevronLeft class="w-4 h-4" />
              Back to Wiki
            </router-link>
          </div>

          <!-- Content -->
          <article v-else-if="wikiStore.currentPage" class="wiki-content">
            <!-- Last Updated -->
            <div class="flex items-center gap-2 text-sm text-gray-500 mb-6">
              <Calendar class="w-4 h-4" />
              Last updated: {{ formattedDate }}
            </div>

            <!-- Main Content with Table of Contents -->
            <div class="lg:flex lg:gap-8">
              <!-- Article Content -->
              <div
                class="prose prose-invert prose-weenie max-w-none flex-1"
                v-html="renderedContent"
              />

              <!-- Desktop Table of Contents -->
              <div class="hidden xl:block w-64 flex-shrink-0">
                <TableOfContents :content="wikiStore.currentPage.content" />
              </div>
            </div>

            <!-- Previous/Next Navigation -->
            <nav class="mt-12 pt-8 border-t border-white/10">
              <div class="flex justify-between gap-4">
                <router-link
                  v-if="adjacentPages.prev"
                  :to="`/wiki/${adjacentPages.prev.slug}`"
                  class="flex-1 card p-4 hover:bg-white/5 transition-colors group"
                >
                  <div class="flex items-center gap-2 text-sm text-gray-400 mb-1">
                    <ChevronLeft class="w-4 h-4" />
                    Previous
                  </div>
                  <div class="text-white group-hover:text-weenie-gold transition-colors font-medium">
                    {{ adjacentPages.prev.title }}
                  </div>
                </router-link>
                <div v-else class="flex-1" />

                <router-link
                  v-if="adjacentPages.next"
                  :to="`/wiki/${adjacentPages.next.slug}`"
                  class="flex-1 card p-4 hover:bg-white/5 transition-colors group text-right"
                >
                  <div class="flex items-center justify-end gap-2 text-sm text-gray-400 mb-1">
                    Next
                    <ChevronRight class="w-4 h-4" />
                  </div>
                  <div class="text-white group-hover:text-weenie-gold transition-colors font-medium">
                    {{ adjacentPages.next.title }}
                  </div>
                </router-link>
              </div>
            </nav>
          </article>
        </div>
      </main>
    </div>
  </div>
</template>

<style>
/* Wiki content styles */
.wiki-content .prose {
  color: #d1d5db;
}

.wiki-content .prose h1 {
  @apply text-3xl font-bold text-white mb-6 mt-0;
}

.wiki-content .prose h2 {
  @apply text-2xl font-semibold text-white mt-10 mb-4 pb-2 border-b border-white/10;
}

.wiki-content .prose h3 {
  @apply text-xl font-semibold text-white mt-8 mb-3;
}

.wiki-content .prose h4 {
  @apply text-lg font-medium text-white mt-6 mb-2;
}

.wiki-content .prose p {
  @apply mb-4 leading-relaxed;
}

.wiki-content .prose ul,
.wiki-content .prose ol {
  @apply mb-4 pl-6;
}

.wiki-content .prose li {
  @apply mb-2;
}

.wiki-content .prose ul li {
  @apply list-disc;
}

.wiki-content .prose ol li {
  @apply list-decimal;
}

.wiki-content .prose a {
  @apply text-weenie-gold hover:underline;
}

.wiki-content .prose code {
  @apply bg-white/10 text-weenie-gold px-1.5 py-0.5 rounded text-sm font-mono;
}

.wiki-content .prose pre {
  @apply bg-black/50 border border-white/10 rounded-lg p-4 overflow-x-auto mb-4;
}

.wiki-content .prose pre code {
  @apply bg-transparent p-0 text-gray-300;
}

.wiki-content .prose blockquote {
  @apply border-l-4 border-weenie-gold pl-4 italic text-gray-400 my-4;
}

.wiki-content .prose table {
  @apply w-full border-collapse mb-4;
}

.wiki-content .prose th {
  @apply text-left p-3 bg-white/5 text-white font-semibold border-b border-white/10;
}

.wiki-content .prose td {
  @apply p-3 border-b border-white/5;
}

.wiki-content .prose tr:hover td {
  @apply bg-white/5;
}

.wiki-content .prose hr {
  @apply border-white/10 my-8;
}

.wiki-content .prose strong {
  @apply text-white font-semibold;
}

.wiki-content .prose img {
  @apply rounded-lg max-w-full h-auto;
}
</style>
