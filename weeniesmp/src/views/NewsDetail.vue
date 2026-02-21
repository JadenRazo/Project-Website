<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { ArrowLeft, Calendar, User, Tag, Clock, Loader2 } from 'lucide-vue-next'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { useNewsStore } from '@/stores/newsStore'
import NewsCard from '@/components/news/NewsCard.vue'

const route = useRoute()
const newsStore = useNewsStore()

const slug = computed(() => route.params.slug as string)

const relatedArticles = computed(() => {
  if (!newsStore.currentArticle) return []
  return newsStore.getRelatedArticles(newsStore.currentArticle, 3)
})

const categoryConfig: Record<string, { color: string; bgColor: string; label: string }> = {
  update: { color: 'text-blue-400', bgColor: 'bg-blue-500/20', label: 'Update' },
  event: { color: 'text-purple-400', bgColor: 'bg-purple-500/20', label: 'Event' },
  patch: { color: 'text-green-400', bgColor: 'bg-green-500/20', label: 'Patch Notes' },
  announcement: { color: 'text-orange-400', bgColor: 'bg-orange-500/20', label: 'Announcement' }
}

const categoryStyle = computed(() => {
  if (!newsStore.currentArticle) return categoryConfig.announcement
  return categoryConfig[newsStore.currentArticle.category] || categoryConfig.announcement
})

const formattedDate = computed(() => {
  if (!newsStore.currentArticle) return ''
  const date = new Date(newsStore.currentArticle.publishedAt)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
})

const formattedUpdateDate = computed(() => {
  if (!newsStore.currentArticle?.updatedAt) return null
  const date = new Date(newsStore.currentArticle.updatedAt)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
})

// Configure marked for safe rendering
marked.setOptions({
  breaks: true,
  gfm: true
})

const renderedContent = computed(() => {
  if (!newsStore.currentArticle) return ''
  const html = marked.parse(newsStore.currentArticle.content) as string
  return DOMPurify.sanitize(html)
})

async function loadArticle() {
  await newsStore.fetchArticle(slug.value)
}

onMounted(() => {
  loadArticle()
})

watch(slug, () => {
  loadArticle()
  window.scrollTo({ top: 0, behavior: 'smooth' })
})
</script>

<template>
  <div class="min-h-screen pt-24 pb-16 bg-weenie-darker">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Back Link -->
      <RouterLink
        to="/news"
        class="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors mb-8"
      >
        <ArrowLeft class="w-5 h-5" />
        Back to News
      </RouterLink>

      <!-- Loading State -->
      <div v-if="newsStore.loading" class="flex flex-col items-center justify-center py-20">
        <Loader2 class="w-10 h-10 text-weenie-gold animate-spin mb-4" />
        <p class="text-gray-400">Loading article...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="newsStore.error" class="text-center py-20">
        <h2 class="text-2xl font-bold text-white mb-4">Article Not Found</h2>
        <p class="text-gray-400 mb-6">{{ newsStore.error }}</p>
        <RouterLink
          to="/news"
          class="px-6 py-3 bg-weenie-gradient text-white font-semibold rounded-lg hover:opacity-90 transition-opacity inline-block"
        >
          Browse All News
        </RouterLink>
      </div>

      <!-- Article Content -->
      <article v-else-if="newsStore.currentArticle" class="card">
        <!-- Featured Image -->
        <div v-if="newsStore.currentArticle.featuredImage" class="relative -mx-6 -mt-6 mb-8 overflow-hidden rounded-t-2xl">
          <img
            :src="newsStore.currentArticle.featuredImage"
            :alt="newsStore.currentArticle.title"
            class="w-full h-64 md:h-80 object-cover"
          />
        </div>

        <!-- Category Badge -->
        <div class="mb-4">
          <span
            class="inline-flex items-center gap-1.5 px-3 py-1 text-sm font-semibold rounded-full"
            :class="[categoryStyle.bgColor, categoryStyle.color]"
          >
            {{ categoryStyle.label }}
          </span>
        </div>

        <!-- Title -->
        <h1 class="text-3xl md:text-4xl font-bold text-white mb-6">
          {{ newsStore.currentArticle.title }}
        </h1>

        <!-- Meta Info -->
        <div class="flex flex-wrap items-center gap-4 mb-8 pb-6 border-b border-white/10">
          <span class="flex items-center gap-2 text-gray-400">
            <Calendar class="w-5 h-5" />
            {{ formattedDate }}
          </span>
          <span class="flex items-center gap-2 text-gray-400">
            <User class="w-5 h-5" />
            {{ newsStore.currentArticle.author }}
          </span>
          <span v-if="formattedUpdateDate" class="flex items-center gap-2 text-gray-500 text-sm">
            <Clock class="w-4 h-4" />
            Updated {{ formattedUpdateDate }}
          </span>
        </div>

        <!-- Article Content -->
        <div
          class="prose prose-invert prose-lg max-w-none"
          v-html="renderedContent"
        />

        <!-- Tags -->
        <div v-if="newsStore.currentArticle.tags.length > 0" class="mt-8 pt-6 border-t border-white/10">
          <div class="flex flex-wrap items-center gap-2">
            <Tag class="w-5 h-5 text-gray-500" />
            <span
              v-for="tag in newsStore.currentArticle.tags"
              :key="tag"
              class="px-3 py-1 text-sm bg-weenie-dark text-gray-400 rounded-full"
            >
              {{ tag }}
            </span>
          </div>
        </div>
      </article>

      <!-- Related Articles -->
      <section v-if="relatedArticles.length > 0" class="mt-16">
        <h2 class="text-2xl font-bold text-white mb-8">Related Articles</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <NewsCard
            v-for="article in relatedArticles"
            :key="article.id"
            :article="article"
          />
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
/* Prose styling for markdown content */
.prose :deep(h1) {
  @apply text-2xl font-bold text-white mt-8 mb-4;
}

.prose :deep(h2) {
  @apply text-xl font-bold text-white mt-8 mb-4;
}

.prose :deep(h3) {
  @apply text-lg font-semibold text-white mt-6 mb-3;
}

.prose :deep(p) {
  @apply text-gray-300 mb-4 leading-relaxed;
}

.prose :deep(ul),
.prose :deep(ol) {
  @apply text-gray-300 mb-4 pl-6;
}

.prose :deep(li) {
  @apply mb-2;
}

.prose :deep(ul) {
  @apply list-disc;
}

.prose :deep(ol) {
  @apply list-decimal;
}

.prose :deep(a) {
  @apply text-weenie-gold hover:underline;
}

.prose :deep(code) {
  @apply bg-weenie-dark px-2 py-0.5 rounded text-weenie-gold text-sm;
}

.prose :deep(pre) {
  @apply bg-weenie-dark p-4 rounded-lg overflow-x-auto mb-4;
}

.prose :deep(pre code) {
  @apply bg-transparent px-0 py-0;
}

.prose :deep(blockquote) {
  @apply border-l-4 border-weenie-gold pl-4 italic text-gray-400 my-4;
}

.prose :deep(strong) {
  @apply text-white font-semibold;
}

.prose :deep(hr) {
  @apply border-white/10 my-8;
}

.prose :deep(img) {
  @apply rounded-lg my-6;
}

.prose :deep(table) {
  @apply w-full border-collapse mb-4;
}

.prose :deep(th),
.prose :deep(td) {
  @apply border border-white/10 px-4 py-2 text-left;
}

.prose :deep(th) {
  @apply bg-weenie-dark text-white font-semibold;
}

.prose :deep(td) {
  @apply text-gray-300;
}
</style>
