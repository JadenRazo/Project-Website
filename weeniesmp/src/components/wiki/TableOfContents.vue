<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { List } from 'lucide-vue-next'

const props = defineProps<{
  content: string
}>()

interface TocItem {
  id: string
  text: string
  level: number
}

const tocItems = ref<TocItem[]>([])
const activeId = ref<string>('')

// Parse headings from content
function parseHeadings(html: string): TocItem[] {
  const parser = new DOMParser()
  const doc = parser.parseFromString(html, 'text/html')
  const headings = doc.querySelectorAll('h1, h2, h3')

  return Array.from(headings).map((heading, index) => {
    const level = parseInt(heading.tagName[1])
    const text = heading.textContent || ''
    const id = heading.id || `heading-${index}`
    return { id, text, level }
  })
}

// Update TOC when content changes
watch(() => props.content, (newContent) => {
  tocItems.value = parseHeadings(newContent)
}, { immediate: true })

// Track active heading on scroll
function handleScroll() {
  const headings = document.querySelectorAll('h1[id], h2[id], h3[id]')
  let currentActive = ''

  for (const heading of headings) {
    const rect = heading.getBoundingClientRect()
    if (rect.top <= 100) {
      currentActive = heading.id
    }
  }

  activeId.value = currentActive
}

function scrollToHeading(id: string) {
  const element = document.getElementById(id)
  if (element) {
    const offset = 100
    const elementPosition = element.getBoundingClientRect().top
    const offsetPosition = elementPosition + window.pageYOffset - offset

    window.scrollTo({
      top: offsetPosition,
      behavior: 'smooth'
    })
  }
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll, { passive: true })
  handleScroll()
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<template>
  <nav v-if="tocItems.length > 0" class="sticky top-24">
    <div class="flex items-center gap-2 mb-4 text-sm font-semibold text-gray-300">
      <List class="w-4 h-4" />
      On This Page
    </div>

    <ul class="space-y-1 border-l border-white/10">
      <li v-for="item in tocItems" :key="item.id">
        <button
          @click="scrollToHeading(item.id)"
          class="block w-full text-left text-sm py-1.5 transition-colors"
          :class="[
            item.level === 1 ? 'pl-4' : item.level === 2 ? 'pl-6' : 'pl-8',
            activeId === item.id
              ? 'text-weenie-gold border-l-2 border-weenie-gold -ml-px'
              : 'text-gray-400 hover:text-white'
          ]"
        >
          {{ item.text }}
        </button>
      </li>
    </ul>
  </nav>
</template>
