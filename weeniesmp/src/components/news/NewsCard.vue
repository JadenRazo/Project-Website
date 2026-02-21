<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { Calendar, User, ArrowRight } from 'lucide-vue-next'
import type { NewsArticle } from '@/stores/newsStore'

const props = defineProps<{
  article: NewsArticle
}>()

const categoryConfig: Record<string, { color: string; bgColor: string }> = {
  update: { color: 'text-blue-400', bgColor: 'bg-blue-500/20' },
  event: { color: 'text-purple-400', bgColor: 'bg-purple-500/20' },
  patch: { color: 'text-green-400', bgColor: 'bg-green-500/20' },
  announcement: { color: 'text-orange-400', bgColor: 'bg-orange-500/20' }
}

const categoryStyle = computed(() => {
  return categoryConfig[props.article.category] || categoryConfig.announcement
})

const categoryLabel = computed(() => {
  const labels: Record<string, string> = {
    update: 'Update',
    event: 'Event',
    patch: 'Patch Notes',
    announcement: 'Announcement'
  }
  return labels[props.article.category] || props.article.category
})

const formattedDate = computed(() => {
  const date = new Date(props.article.publishedAt)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
})

const fallbackImage = '/images/news/default.png'
</script>

<template>
  <RouterLink
    :to="`/news/${article.slug}`"
    class="card group hover:scale-[1.02] flex flex-col cursor-pointer overflow-hidden"
  >
    <!-- Featured Image -->
    <div class="relative w-full h-48 -mx-6 -mt-6 mb-4 overflow-hidden">
      <img
        :src="article.featuredImage || fallbackImage"
        :alt="article.title"
        class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
        @error="($event.target as HTMLImageElement).src = fallbackImage"
      />
      <!-- Category Badge -->
      <span
        class="absolute top-3 left-3 px-3 py-1 text-xs font-semibold rounded-full"
        :class="[categoryStyle.bgColor, categoryStyle.color]"
      >
        {{ categoryLabel }}
      </span>
    </div>

    <!-- Content -->
    <div class="flex-1 flex flex-col">
      <h3 class="text-xl font-bold text-white mb-2 group-hover:text-weenie-gold transition-colors line-clamp-2">
        {{ article.title }}
      </h3>

      <p class="text-gray-400 text-sm mb-4 flex-1 line-clamp-3">
        {{ article.excerpt }}
      </p>

      <!-- Meta Info -->
      <div class="flex items-center justify-between text-sm text-gray-500">
        <div class="flex items-center gap-4">
          <span class="flex items-center gap-1.5">
            <Calendar class="w-4 h-4" />
            {{ formattedDate }}
          </span>
          <span class="flex items-center gap-1.5">
            <User class="w-4 h-4" />
            {{ article.author }}
          </span>
        </div>
      </div>

      <!-- Read More -->
      <div class="mt-4 pt-4 border-t border-white/10">
        <span class="flex items-center gap-2 text-weenie-gold font-medium group-hover:gap-3 transition-all">
          Read More
          <ArrowRight class="w-4 h-4" />
        </span>
      </div>
    </div>
  </RouterLink>
</template>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
