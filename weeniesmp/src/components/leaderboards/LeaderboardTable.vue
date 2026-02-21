<script setup lang="ts">
import { Loader2, AlertCircle, RefreshCw } from 'lucide-vue-next'
import PlayerRow from './PlayerRow.vue'
import type { LeaderboardEntry, LeaderboardCategory } from '@/stores/leaderboardStore'

defineProps<{
  entries: LeaderboardEntry[]
  loading: boolean
  initialLoading: boolean
  error: string | null
  category: LeaderboardCategory
}>()

const emit = defineEmits<{
  (e: 'retry'): void
}>()
</script>

<template>
  <div class="card relative">
    <!-- Initial Loading State (only when no cached data) -->
    <div v-if="initialLoading" class="py-12">
      <div class="flex flex-col items-center justify-center">
        <Loader2 class="w-10 h-10 text-weenie-red animate-spin mb-4" />
        <p class="text-gray-400">Loading leaderboard...</p>
      </div>
      <!-- Skeleton Rows -->
      <div class="mt-8 space-y-3">
        <div
          v-for="i in 10"
          :key="i"
          class="flex items-center gap-4 p-4 bg-weenie-dark/30 rounded-lg animate-pulse"
        >
          <div class="w-8 h-8 bg-gray-700 rounded-full"></div>
          <div class="w-10 h-10 bg-gray-700 rounded"></div>
          <div class="flex-1">
            <div class="h-4 bg-gray-700 rounded w-32 mb-2"></div>
            <div class="h-3 bg-gray-700 rounded w-20"></div>
          </div>
          <div class="h-5 bg-gray-700 rounded w-16"></div>
        </div>
      </div>
    </div>

    <!-- Error State (only if no cached data) -->
    <div v-else-if="error && entries.length === 0" class="py-12 text-center">
      <AlertCircle class="w-12 h-12 text-red-400 mx-auto mb-4" />
      <p class="text-red-400 mb-4">{{ error }}</p>
      <button
        @click="emit('retry')"
        class="inline-flex items-center gap-2 px-6 py-3 bg-weenie-red text-white rounded-lg hover:bg-weenie-red/80 transition-colors"
      >
        <RefreshCw class="w-5 h-5" />
        Retry
      </button>
    </div>

    <!-- Empty State -->
    <div v-else-if="entries.length === 0" class="py-12 text-center">
      <p class="text-gray-400">No data available for this leaderboard yet.</p>
    </div>

    <!-- Leaderboard Entries (with smooth transitions) -->
    <div v-else class="divide-y divide-white/5">
      <transition-group
        name="list"
        tag="div"
      >
        <PlayerRow
          v-for="(entry, index) in entries"
          :key="entry.uuid"
          :entry="entry"
          :rank="index + 1"
        />
      </transition-group>
    </div>
  </div>
</template>

<style scoped>
/* Smooth list transitions */
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}

.list-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.list-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

.list-move {
  transition: transform 0.3s ease;
}
</style>
