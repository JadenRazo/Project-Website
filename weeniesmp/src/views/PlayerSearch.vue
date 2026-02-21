<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { Search, Users, Clock, TrendingUp } from 'lucide-vue-next'
import { usePlayerStore } from '@/stores/playerStore'
import PlayerCard from '@/components/players/PlayerCard.vue'

const playerStore = usePlayerStore()
const searchQuery = ref('')
const searchInputRef = ref<HTMLInputElement | null>(null)

// Debounce search
let searchTimeout: ReturnType<typeof setTimeout> | null = null

watch(searchQuery, (newQuery) => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  searchTimeout = setTimeout(() => {
    playerStore.searchPlayers(newQuery)
  }, 300)
})

onMounted(() => {
  playerStore.fetchRecentPlayers()
  playerStore.fetchPopularPlayers()
  // Focus search input on mount
  searchInputRef.value?.focus()
})
</script>

<template>
  <div class="min-h-screen pt-20 pb-12">
    <div class="max-w-4xl mx-auto px-4 sm:px-6">
      <!-- Header -->
      <div class="text-center mb-10">
        <h1 class="text-4xl font-bold text-white mb-3">
          <span class="gradient-text">Player</span> Directory
        </h1>
        <p class="text-gray-400 max-w-md mx-auto">
          Search for players, view their stats, achievements, and more.
        </p>
      </div>

      <!-- Search box -->
      <div class="relative max-w-xl mx-auto mb-12">
        <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
          <Search class="w-5 h-5 text-gray-400" />
        </div>
        <input
          ref="searchInputRef"
          v-model="searchQuery"
          type="text"
          placeholder="Search by username..."
          class="w-full pl-12 pr-4 py-4 bg-weenie-darker/50 border border-white/10 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:border-weenie-gold/50 focus:ring-2 focus:ring-weenie-gold/20 transition-all"
        />
        <!-- Loading indicator -->
        <div
          v-if="playerStore.searchLoading"
          class="absolute inset-y-0 right-0 pr-4 flex items-center"
        >
          <div class="w-5 h-5 border-2 border-weenie-gold/30 border-t-weenie-gold rounded-full animate-spin"></div>
        </div>
      </div>

      <!-- Search results -->
      <section v-if="searchQuery && !playerStore.searchLoading" class="mb-12">
        <div v-if="playerStore.searchResults.length > 0" class="space-y-3">
          <h2 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            <Search class="w-5 h-5 text-weenie-gold" />
            Search Results
          </h2>
          <div class="grid gap-3 sm:grid-cols-2">
            <PlayerCard
              v-for="player in playerStore.searchResults"
              :key="player.uuid"
              :player="player"
            />
          </div>
        </div>
        <div v-else class="text-center py-10">
          <Users class="w-12 h-12 text-gray-500 mx-auto mb-3" />
          <p class="text-gray-400">No players found matching "{{ searchQuery }}"</p>
          <p class="text-sm text-gray-500 mt-1">Try a different search term</p>
        </div>
      </section>

      <!-- Recent & Popular players (shown when not searching) -->
      <div v-if="!searchQuery" class="space-y-12">
        <!-- Online / Recent players -->
        <section v-if="playerStore.recentPlayers.length > 0">
          <h2 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            <Clock class="w-5 h-5 text-weenie-gold" />
            Recently Active
          </h2>
          <div class="grid gap-3 sm:grid-cols-2">
            <PlayerCard
              v-for="player in playerStore.recentPlayers"
              :key="player.uuid"
              :player="player"
            />
          </div>
        </section>

        <!-- Popular players -->
        <section v-if="playerStore.popularPlayers.length > 0">
          <h2 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            <TrendingUp class="w-5 h-5 text-weenie-gold" />
            Top Players
          </h2>
          <div class="grid gap-3 sm:grid-cols-2">
            <PlayerCard
              v-for="player in playerStore.popularPlayers"
              :key="player.uuid"
              :player="player"
            />
          </div>
        </section>

        <!-- Empty state -->
        <div
          v-if="playerStore.recentPlayers.length === 0 && playerStore.popularPlayers.length === 0"
          class="text-center py-16"
        >
          <Users class="w-16 h-16 text-gray-500 mx-auto mb-4" />
          <h3 class="text-lg font-medium text-white mb-2">No Players Yet</h3>
          <p class="text-gray-400">Be the first to join WeenieSMP!</p>
        </div>
      </div>
    </div>
  </div>
</template>
