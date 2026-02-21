<script setup lang="ts">
import { onMounted, ref, computed, watch } from 'vue'
import { Trophy, Clock, DollarSign, MapPin, Layers, Vote, Briefcase } from 'lucide-vue-next'
import LeaderboardTabs from '@/components/leaderboards/LeaderboardTabs.vue'
import LeaderboardPodium from '@/components/leaderboards/LeaderboardPodium.vue'
import LeaderboardTable from '@/components/leaderboards/LeaderboardTable.vue'
import LeaderboardStatus from '@/components/leaderboards/LeaderboardStatus.vue'
import PlayerSearchBar from '@/components/leaderboards/PlayerSearchBar.vue'
import { useLeaderboardStore, type LeaderboardCategory } from '@/stores/leaderboardStore'

const leaderboardStore = useLeaderboardStore()

const categories = [
  { id: 'playtime' as LeaderboardCategory, name: 'Playtime', icon: Clock },
  { id: 'economy' as LeaderboardCategory, name: 'Economy', icon: DollarSign },
  { id: 'claims' as LeaderboardCategory, name: 'Claims', icon: MapPin },
  { id: 'chunks' as LeaderboardCategory, name: 'Chunks', icon: Layers },
  { id: 'jobs' as LeaderboardCategory, name: 'Jobs', icon: Briefcase },
  { id: 'voting' as LeaderboardCategory, name: 'Voting', icon: Vote }
]

const activeCategory = ref<LeaderboardCategory>('playtime')
const searchQuery = ref('')

const filteredEntries = computed(() => {
  if (!searchQuery.value.trim()) {
    return leaderboardStore.entries
  }
  const query = searchQuery.value.toLowerCase()
  return leaderboardStore.entries.filter(entry =>
    entry.username.toLowerCase().includes(query)
  )
})

async function handleCategoryChange(category: LeaderboardCategory) {
  activeCategory.value = category
  searchQuery.value = ''
  await leaderboardStore.fetchLeaderboard(category)
}

function handleSearch(query: string) {
  searchQuery.value = query
}

function handleClearSearch() {
  searchQuery.value = ''
}

// Prefetch adjacent categories for instant switching
watch(activeCategory, (newCategory) => {
  const currentIndex = categories.findIndex(c => c.id === newCategory)

  // Prefetch next category
  if (currentIndex < categories.length - 1) {
    leaderboardStore.prefetchCategory(categories[currentIndex + 1].id)
  }

  // Prefetch previous category
  if (currentIndex > 0) {
    leaderboardStore.prefetchCategory(categories[currentIndex - 1].id)
  }
})

onMounted(async () => {
  // Load initial category (uses cache if available)
  await leaderboardStore.fetchLeaderboard(activeCategory.value)

  // Prefetch the next category for instant switching
  if (categories.length > 1) {
    leaderboardStore.prefetchCategory(categories[1].id)
  }
})
</script>

<template>
  <div class="min-h-screen pt-24 pb-16 bg-weenie-darker">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="text-center mb-12">
        <Trophy class="w-12 h-12 text-weenie-gold mx-auto mb-4" />
        <h1 class="text-4xl md:text-5xl font-bold mb-4">
          <span class="gradient-text">Leaderboards</span>
        </h1>
        <p class="text-gray-400 max-w-xl mx-auto">
          See who's at the top of WeenieSMP. Compete with others and climb the ranks!
        </p>
      </div>

      <LeaderboardTabs
        :categories="categories"
        :active-category="activeCategory"
        @change="handleCategoryChange"
      />

      <LeaderboardPodium :entries="leaderboardStore.entries" />

      <LeaderboardStatus
        :updating="leaderboardStore.updating"
        :last-updated="leaderboardStore.lastUpdated"
        :cache-age="leaderboardStore.cacheAge"
      />

      <PlayerSearchBar
        @search="handleSearch"
        @clear="handleClearSearch"
      />

      <LeaderboardTable
        :entries="filteredEntries"
        :loading="leaderboardStore.loading"
        :initial-loading="leaderboardStore.initialLoading"
        :error="leaderboardStore.error"
        :category="activeCategory"
        @retry="leaderboardStore.refresh()"
      />
    </div>
  </div>
</template>
