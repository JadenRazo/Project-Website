<script setup lang="ts">
import { computed } from 'vue'
import { Trophy, Medal, Award } from 'lucide-vue-next'
import type { LeaderboardEntry } from '@/stores/leaderboardStore'

const props = defineProps<{
  entries: LeaderboardEntry[]
}>()

const topThree = computed(() => {
  if (props.entries.length === 0) return []

  // Get top 3 entries
  const top = props.entries.slice(0, 3)

  // Reorder to show 2nd, 1st, 3rd for podium layout
  if (top.length === 3) {
    return [top[1], top[0], top[2]]
  } else if (top.length === 2) {
    return [undefined, top[0], top[1]]
  } else {
    return [undefined, top[0], undefined]
  }
})

function getAvatarUrl(uuid: string) {
  return `https://mc-heads.net/avatar/${uuid}/64`
}

function getPodiumClass(rank: number) {
  switch (rank) {
    case 1:
      return 'h-48 bg-gradient-to-br from-yellow-500/20 to-amber-600/20 border-yellow-500/50'
    case 2:
      return 'h-40 bg-gradient-to-br from-gray-300/20 to-gray-400/20 border-gray-300/50'
    case 3:
      return 'h-36 bg-gradient-to-br from-amber-600/20 to-amber-700/20 border-amber-600/50'
    default:
      return 'h-32 bg-weenie-dark/50 border-white/10'
  }
}

function getTrophyColor(rank: number) {
  switch (rank) {
    case 1:
      return 'text-yellow-400'
    case 2:
      return 'text-gray-300'
    case 3:
      return 'text-amber-600'
    default:
      return 'text-gray-500'
  }
}

function getFormattedValue(entry: LeaderboardEntry) {
  return entry.metadata?.formatted || entry.value.toLocaleString()
}
</script>

<template>
  <div v-if="entries.length >= 3" class="mb-12">
    <h2 class="text-2xl font-bold text-center mb-8 text-white">Top 3 Players</h2>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 items-end max-w-3xl mx-auto">
      <!-- 2nd Place -->
      <div v-if="topThree[0]" class="flex flex-col items-center md:order-1">
        <div class="relative mb-4">
          <img
            :src="getAvatarUrl(topThree[0].uuid)"
            :alt="topThree[0].username"
            class="w-12 h-12 sm:w-16 sm:h-16 rounded-lg border-2 border-gray-300 shadow-lg"
          />
          <div class="absolute -bottom-2 -right-2 w-8 h-8 rounded-full bg-gray-300 flex items-center justify-center text-black font-bold text-sm">
            2
          </div>
        </div>

        <div :class="getPodiumClass(2)" class="w-full rounded-t-lg border-2 flex flex-col items-center justify-center p-4">
          <Medal :class="getTrophyColor(2)" class="w-8 h-8 mb-2" />
          <p class="text-white font-semibold truncate w-full text-center">{{ topThree[0].username }}</p>
          <p class="text-gray-300 text-sm">{{ getFormattedValue(topThree[0]) }}</p>
        </div>
      </div>

      <!-- 1st Place -->
      <div v-if="topThree[1]" class="flex flex-col items-center md:order-2">
        <div class="relative mb-4">
          <img
            :src="getAvatarUrl(topThree[1].uuid)"
            :alt="topThree[1].username"
            class="w-16 h-16 sm:w-20 sm:h-20 rounded-lg border-2 border-yellow-400 shadow-xl shadow-yellow-400/50"
          />
          <div class="absolute -bottom-2 -right-2 w-10 h-10 rounded-full bg-yellow-400 flex items-center justify-center text-black font-bold">
            1
          </div>
        </div>

        <div :class="getPodiumClass(1)" class="w-full rounded-t-lg border-2 flex flex-col items-center justify-center p-4 relative overflow-hidden">
          <div class="absolute inset-0 bg-gradient-to-t from-yellow-500/10 to-transparent"></div>
          <Trophy :class="getTrophyColor(1)" class="w-10 h-10 mb-2 relative z-10" />
          <p class="text-white font-bold text-lg truncate w-full text-center relative z-10">{{ topThree[1].username }}</p>
          <p class="text-yellow-300 font-semibold relative z-10">{{ getFormattedValue(topThree[1]) }}</p>
        </div>
      </div>

      <!-- 3rd Place -->
      <div v-if="topThree[2]" class="flex flex-col items-center md:order-3">
        <div class="relative mb-4">
          <img
            :src="getAvatarUrl(topThree[2].uuid)"
            :alt="topThree[2].username"
            class="w-12 h-12 sm:w-16 sm:h-16 rounded-lg border-2 border-amber-600 shadow-lg"
          />
          <div class="absolute -bottom-2 -right-2 w-8 h-8 rounded-full bg-amber-600 flex items-center justify-center text-white font-bold text-sm">
            3
          </div>
        </div>

        <div :class="getPodiumClass(3)" class="w-full rounded-t-lg border-2 flex flex-col items-center justify-center p-4">
          <Award :class="getTrophyColor(3)" class="w-8 h-8 mb-2" />
          <p class="text-white font-semibold truncate w-full text-center">{{ topThree[2].username }}</p>
          <p class="text-gray-300 text-sm">{{ getFormattedValue(topThree[2]) }}</p>
        </div>
      </div>
    </div>
  </div>
</template>
