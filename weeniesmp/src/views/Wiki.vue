<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useWikiStore } from '@/stores/wikiStore'
import WikiSearch from '@/components/wiki/WikiSearch.vue'
import { BookOpen, ChevronRight } from 'lucide-vue-next'

const router = useRouter()
const wikiStore = useWikiStore()

onMounted(() => {
  wikiStore.fetchCategories()
})

function navigateToPage(slug: string) {
  router.push(`/wiki/${slug}`)
}
</script>

<template>
  <div class="min-h-screen pt-24 pb-16 bg-weenie-darker">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="text-center mb-12">
        <BookOpen class="w-12 h-12 text-weenie-gold mx-auto mb-4" />
        <h1 class="text-4xl md:text-5xl font-bold mb-4">
          <span class="gradient-text">Server Wiki</span>
        </h1>
        <p class="text-gray-400 max-w-xl mx-auto mb-8">
          Everything you need to know about WeenieSMP. Browse guides, commands, and features.
        </p>

        <!-- Search -->
        <div class="max-w-md mx-auto">
          <WikiSearch @select="navigateToPage" />
        </div>
      </div>

      <!-- Category Grid -->
      <div class="grid md:grid-cols-2 gap-6">
        <div
          v-for="category in wikiStore.categories"
          :key="category.id"
          class="card p-6"
        >
          <h2 class="text-xl font-semibold text-white mb-4 flex items-center gap-2">
            <component
              :is="category.icon === 'book-open' ? BookOpen : BookOpen"
              class="w-5 h-5 text-weenie-gold"
            />
            {{ category.name }}
          </h2>

          <ul class="space-y-2">
            <li v-for="page in category.pages" :key="page.slug">
              <router-link
                :to="`/wiki/${page.slug}`"
                class="flex items-center justify-between p-3 rounded-lg hover:bg-white/5 transition-colors group"
              >
                <span class="text-gray-300 group-hover:text-white transition-colors">
                  {{ page.title }}
                </span>
                <ChevronRight class="w-4 h-4 text-gray-500 group-hover:text-weenie-gold transition-colors" />
              </router-link>
            </li>
          </ul>
        </div>
      </div>

      <!-- Quick Links -->
      <div class="mt-12 card p-6">
        <h2 class="text-xl font-semibold text-white mb-4">Quick Links</h2>
        <div class="grid sm:grid-cols-2 md:grid-cols-4 gap-4">
          <router-link
            to="/wiki/getting-started"
            class="p-4 rounded-lg bg-white/5 hover:bg-white/10 transition-colors text-center"
          >
            <span class="text-2xl mb-2 block">New Here?</span>
            <span class="text-gray-400">Getting Started Guide</span>
          </router-link>
          <router-link
            to="/wiki/commands"
            class="p-4 rounded-lg bg-white/5 hover:bg-white/10 transition-colors text-center"
          >
            <span class="text-2xl mb-2 block">Commands</span>
            <span class="text-gray-400">Full Command List</span>
          </router-link>
          <router-link
            to="/wiki/claims"
            class="p-4 rounded-lg bg-white/5 hover:bg-white/10 transition-colors text-center"
          >
            <span class="text-2xl mb-2 block">Protection</span>
            <span class="text-gray-400">Claim Your Land</span>
          </router-link>
          <router-link
            to="/wiki/economy"
            class="p-4 rounded-lg bg-white/5 hover:bg-white/10 transition-colors text-center"
          >
            <span class="text-2xl mb-2 block">Economy</span>
            <span class="text-gray-400">Make Money</span>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>
