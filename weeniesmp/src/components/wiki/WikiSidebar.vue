<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import type { WikiCategory } from '@/stores/wikiStore'
import { BookOpen, Star, ChevronDown, ChevronRight, FileText, Home } from 'lucide-vue-next'

const props = defineProps<{
  categories: WikiCategory[]
  currentSlug: string
}>()

const emit = defineEmits<{
  close: []
}>()

// Track expanded categories
const expandedCategories = ref<Set<string>>(new Set(props.categories.map(c => c.id)))

function toggleCategory(categoryId: string) {
  if (expandedCategories.value.has(categoryId)) {
    expandedCategories.value.delete(categoryId)
  } else {
    expandedCategories.value.add(categoryId)
  }
}

function isExpanded(categoryId: string): boolean {
  return expandedCategories.value.has(categoryId)
}

function getCategoryIcon(iconName: string) {
  switch (iconName) {
    case 'book-open':
      return BookOpen
    case 'star':
      return Star
    default:
      return FileText
  }
}

// Check if current page is in a category
function categoryHasCurrentPage(category: WikiCategory): boolean {
  return category.pages.some(page => page.slug === props.currentSlug)
}
</script>

<template>
  <div class="p-4">
    <!-- Wiki Home Link -->
    <RouterLink
      to="/wiki"
      class="flex items-center gap-2 px-3 py-2 mb-4 rounded-lg hover:bg-white/5 transition-colors text-gray-400 hover:text-white"
      @click="emit('close')"
    >
      <Home class="w-5 h-5" />
      <span class="font-medium">Wiki Home</span>
    </RouterLink>

    <!-- Categories -->
    <nav class="space-y-2">
      <div
        v-for="category in categories"
        :key="category.id"
        class="rounded-lg overflow-hidden"
      >
        <!-- Category Header -->
        <button
          @click="toggleCategory(category.id)"
          class="w-full flex items-center justify-between px-3 py-2 text-left hover:bg-white/5 rounded-lg transition-colors"
          :class="categoryHasCurrentPage(category) ? 'text-weenie-gold' : 'text-gray-300'"
        >
          <div class="flex items-center gap-2">
            <component :is="getCategoryIcon(category.icon)" class="w-5 h-5" />
            <span class="font-medium">{{ category.name }}</span>
          </div>
          <ChevronDown
            v-if="isExpanded(category.id)"
            class="w-4 h-4 text-gray-500"
          />
          <ChevronRight
            v-else
            class="w-4 h-4 text-gray-500"
          />
        </button>

        <!-- Category Pages -->
        <Transition
          enter-active-class="transition-all duration-200 ease-out"
          enter-from-class="opacity-0 max-h-0"
          enter-to-class="opacity-100 max-h-96"
          leave-active-class="transition-all duration-150 ease-in"
          leave-from-class="opacity-100 max-h-96"
          leave-to-class="opacity-0 max-h-0"
        >
          <ul v-if="isExpanded(category.id)" class="overflow-hidden ml-4 mt-1 space-y-1">
            <li v-for="page in category.pages" :key="page.slug">
              <RouterLink
                :to="`/wiki/${page.slug}`"
                class="flex items-center gap-2 px-3 py-2 rounded-lg transition-colors text-sm"
                :class="page.slug === currentSlug
                  ? 'bg-weenie-gold/10 text-weenie-gold border-l-2 border-weenie-gold'
                  : 'text-gray-400 hover:text-white hover:bg-white/5'"
                @click="emit('close')"
              >
                <FileText class="w-4 h-4 flex-shrink-0" />
                <span>{{ page.title }}</span>
              </RouterLink>
            </li>
          </ul>
        </Transition>
      </div>
    </nav>

    <!-- Footer Links -->
    <div class="mt-8 pt-4 border-t border-white/10">
      <a
        href="https://discord.com/invite/weeniesmp"
        target="_blank"
        rel="noopener noreferrer"
        class="flex items-center gap-2 px-3 py-2 text-sm text-gray-400 hover:text-white hover:bg-white/5 rounded-lg transition-colors"
      >
        Need help? Join Discord
      </a>
    </div>
  </div>
</template>
