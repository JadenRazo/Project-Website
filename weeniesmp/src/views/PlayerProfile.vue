<script setup lang="ts">
import { watch, onMounted } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { ArrowLeft, Calendar, Clock, Wifi, WifiOff } from 'lucide-vue-next'
import { usePlayerStore } from '@/stores/playerStore'
import { formatPlaytime, formatDate, formatFullDate } from '@/utils/formatters'
import PlayerSkinViewer from '@/components/players/PlayerSkinViewer.vue'
import PlayerStats from '@/components/players/PlayerStats.vue'
import AchievementGrid from '@/components/players/AchievementGrid.vue'

const route = useRoute()
const playerStore = usePlayerStore()

// Fetch profile on mount and when username changes
onMounted(() => {
  const username = route.params.username as string
  if (username) {
    playerStore.fetchProfile(username)
  }
})

watch(
  () => route.params.username,
  (newUsername) => {
    if (newUsername && typeof newUsername === 'string') {
      playerStore.fetchProfile(newUsername)
    }
  }
)
</script>

<template>
  <div class="min-h-screen pt-20 pb-12">
    <div class="max-w-5xl mx-auto px-4 sm:px-6">
      <!-- Back button -->
      <RouterLink
        to="/players"
        class="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors mb-6"
      >
        <ArrowLeft class="w-4 h-4" />
        Back to Players
      </RouterLink>

      <!-- Loading state -->
      <div v-if="playerStore.loading" class="text-center py-20">
        <div class="inline-block w-12 h-12 border-4 border-weenie-gold/30 border-t-weenie-gold rounded-full animate-spin"></div>
        <p class="mt-4 text-gray-400">Loading player profile...</p>
      </div>

      <!-- Error / Not found state -->
      <div v-else-if="playerStore.error" class="text-center py-20">
        <div class="w-24 h-24 mx-auto mb-6 rounded-full bg-weenie-darker/50 flex items-center justify-center">
          <svg class="w-12 h-12 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <h1 class="text-2xl font-bold text-white mb-2">Player Not Found</h1>
        <p class="text-gray-400 mb-6">{{ playerStore.error }}</p>
        <RouterLink to="/players" class="btn-primary">
          Search Players
        </RouterLink>
      </div>

      <!-- Profile content -->
      <div v-else-if="playerStore.profile" class="space-y-8">
        <!-- Header -->
        <div class="card !p-0 overflow-hidden">
          <div class="bg-weenie-gradient/10 p-8">
            <div class="flex flex-col md:flex-row items-center gap-8">
              <!-- Skin viewer -->
              <PlayerSkinViewer
                :uuid="playerStore.profile.uuid"
                :username="playerStore.profile.username"
                size="xl"
              />

              <!-- Player info -->
              <div class="flex-1 text-center md:text-left">
                <div class="flex items-center justify-center md:justify-start gap-3 mb-2">
                  <h1 class="text-4xl font-bold text-white">
                    {{ playerStore.profile.username }}
                  </h1>
                  <!-- Online status -->
                  <span
                    v-if="playerStore.profile.isOnline"
                    class="inline-flex items-center gap-1.5 px-3 py-1 bg-green-500/20 text-green-400 text-sm font-medium rounded-full"
                  >
                    <Wifi class="w-4 h-4" />
                    Online
                  </span>
                  <span
                    v-else
                    class="inline-flex items-center gap-1.5 px-3 py-1 bg-gray-500/20 text-gray-400 text-sm font-medium rounded-full"
                  >
                    <WifiOff class="w-4 h-4" />
                    Offline
                  </span>
                </div>

                <!-- Quick stats -->
                <div class="flex flex-wrap items-center justify-center md:justify-start gap-4 text-gray-400 mt-4">
                  <div class="flex items-center gap-2">
                    <Calendar class="w-4 h-4" />
                    <span>Joined {{ formatFullDate(playerStore.profile.firstJoin) }}</span>
                  </div>
                  <div class="flex items-center gap-2">
                    <Clock class="w-4 h-4" />
                    <span>{{ formatPlaytime(playerStore.profile.playtime) }} playtime</span>
                  </div>
                </div>

                <!-- Last seen (if offline) -->
                <p v-if="!playerStore.profile.isOnline" class="text-sm text-gray-500 mt-2">
                  Last seen {{ formatDate(playerStore.profile.lastSeen) }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Stats section -->
        <section>
          <h2 class="text-xl font-bold text-white mb-4">Statistics</h2>
          <PlayerStats :stats="playerStore.profile.stats" />
        </section>

        <!-- Achievements section -->
        <section>
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-xl font-bold text-white">Achievements</h2>
            <span class="text-sm text-gray-400">
              {{ playerStore.profile.achievements.filter(a => a.unlockedAt).length }} / {{ playerStore.profile.achievements.length }} unlocked
            </span>
          </div>
          <AchievementGrid :achievements="playerStore.profile.achievements" />
        </section>
      </div>
    </div>
  </div>
</template>
