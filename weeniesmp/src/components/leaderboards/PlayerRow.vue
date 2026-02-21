<script setup lang="ts">
import { computed } from 'vue'
import type { LeaderboardEntry } from '@/stores/leaderboardStore'

const props = defineProps<{
  entry: LeaderboardEntry
  rank: number
}>()

const avatarUrl = computed(() => {
  return `https://mc-heads.net/avatar/${props.entry.uuid}/32`
})

const rankClass = computed(() => {
  switch (props.rank) {
    case 1:
      return 'bg-gradient-to-r from-yellow-500 to-amber-400 text-black font-bold'
    case 2:
      return 'bg-gradient-to-r from-gray-300 to-gray-400 text-black font-bold'
    case 3:
      return 'bg-gradient-to-r from-amber-600 to-amber-700 text-white font-bold'
    default:
      return 'bg-weenie-dark text-gray-400'
  }
})

const rowClass = computed(() => {
  switch (props.rank) {
    case 1:
      return 'bg-yellow-500/5 hover:bg-yellow-500/10'
    case 2:
      return 'bg-gray-300/5 hover:bg-gray-300/10'
    case 3:
      return 'bg-amber-600/5 hover:bg-amber-600/10'
    default:
      return 'hover:bg-white/5'
  }
})
</script>

<template>
  <div
    class="flex items-center gap-2 sm:gap-4 p-2 sm:p-4 transition-colors rounded-lg"
    :class="rowClass"
  >
    <!-- Rank Badge -->
    <div
      class="w-8 h-8 sm:w-10 sm:h-10 rounded-lg flex items-center justify-center text-xs sm:text-sm font-semibold shrink-0"
      :class="rankClass"
    >
      {{ rank }}
    </div>

    <!-- Player Avatar -->
    <img
      :src="avatarUrl"
      :alt="entry.username"
      class="w-8 h-8 sm:w-10 sm:h-10 rounded bg-weenie-dark"
      loading="lazy"
    />

    <!-- Player Info -->
    <div class="flex-1 min-w-0">
      <p class="text-sm sm:text-base text-white font-medium truncate">{{ entry.username }}</p>
      <p v-if="entry.metadata?.job" class="text-gray-500 text-sm">{{ entry.metadata.job }}</p>
    </div>

    <!-- Value -->
    <div class="text-right shrink-0">
      <span
        class="text-base sm:text-lg font-semibold"
        :class="rank <= 3 ? 'text-weenie-gold' : 'text-white'"
      >
        {{ entry.metadata?.formatted || entry.value.toLocaleString() }}
      </span>
    </div>
  </div>
</template>
