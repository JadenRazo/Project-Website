<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Users, Clock, Blocks, Sword, TrendingUp, Activity } from 'lucide-vue-next'
import { apiClient, ApiError } from '@/utils/api'

interface ServerStats {
  totalPlayers: number
  totalPlaytime: number
  blocksPlaced: number
  mobsKilled: number
  uniquePlayersToday: number
  peakOnlineToday: number
}

const stats = ref<ServerStats | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const displayValues = ref<Record<string, number>>({})

const statItems = [
  { key: 'totalPlayers', label: 'Total Players', icon: Users, color: 'text-blue-400', format: (v: number) => v.toLocaleString() },
  { key: 'totalPlaytime', label: 'Hours Played', icon: Clock, color: 'text-green-400', format: (v: number) => v.toLocaleString() },
  { key: 'blocksPlaced', label: 'Blocks Placed', icon: Blocks, color: 'text-purple-400', format: (v: number) => formatLargeNumber(v) },
  { key: 'mobsKilled', label: 'Mobs Killed', icon: Sword, color: 'text-red-400', format: (v: number) => formatLargeNumber(v) },
  { key: 'uniquePlayersToday', label: 'Players Today', icon: TrendingUp, color: 'text-yellow-400', format: (v: number) => v.toLocaleString() },
  { key: 'peakOnlineToday', label: 'Peak Online', icon: Activity, color: 'text-cyan-400', format: (v: number) => v.toLocaleString() }
]

function formatLargeNumber(num: number): string {
  if (num >= 1000000000) {
    return (num / 1000000000).toFixed(1) + 'B'
  }
  if (num >= 1000000) {
    return (num / 1000000).toFixed(1) + 'M'
  }
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K'
  }
  return num.toLocaleString()
}

function animateNumber(key: string, target: number, duration = 1500) {
  const start = displayValues.value[key] || 0
  const startTime = performance.now()

  function update(currentTime: number) {
    const elapsed = currentTime - startTime
    const progress = Math.min(elapsed / duration, 1)

    // Easing function (ease-out cubic)
    const eased = 1 - Math.pow(1 - progress, 3)

    displayValues.value[key] = Math.round(start + (target - start) * eased)

    if (progress < 1) {
      requestAnimationFrame(update)
    }
  }

  requestAnimationFrame(update)
}

async function fetchStats() {
  loading.value = true
  error.value = null

  try {
    // Fetch from MC Stats API (Go backend)
    const statsApiUrl = import.meta.env.VITE_STATS_API_URL || 'https://weeniesmp.net/api/mc'
    const data = await apiClient.get<ServerStats>(
      `${statsApiUrl}/stats`,
      { retries: 3 }
    )

    stats.value = data

    // Animate each stat value
    for (const item of statItems) {
      const value = data[item.key as keyof ServerStats] ?? 0
      animateNumber(item.key, value)
    }
  } catch (e) {
    if (e instanceof ApiError) {
      error.value = e.message
    } else {
      error.value = 'Failed to load server stats'
    }
    console.error('Stats fetch error:', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  // Initialize display values
  for (const item of statItems) {
    displayValues.value[item.key] = 0
  }
  fetchStats()
})
</script>

<template>
  <section class="py-6 bg-weenie-dark">
    <div class="max-w-4xl mx-auto px-6">
      <h2 class="text-xl font-semibold text-center text-white mb-5">Server Stats</h2>

      <!-- Loading skeleton -->
      <div v-if="loading" class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4">
        <div
          v-for="i in 6"
          :key="i"
          class="flex flex-col items-center gap-2 p-4 rounded-lg bg-white/[0.02] border border-white/5"
        >
          <div class="w-8 h-8 rounded-full bg-white/5 animate-pulse"></div>
          <div class="w-16 h-6 rounded bg-white/5 animate-pulse"></div>
          <div class="w-12 h-4 rounded bg-white/5 animate-pulse"></div>
        </div>
      </div>

      <!-- Error state -->
      <div v-else-if="error" class="text-center py-8">
        <p class="text-gray-500 text-sm">{{ error }}</p>
        <button
          @click="fetchStats"
          class="mt-3 px-4 py-2 text-sm text-gray-400 hover:text-white bg-white/5 hover:bg-white/10 rounded-lg transition-colors"
        >
          Try Again
        </button>
      </div>

      <!-- Stats grid -->
      <div v-else class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4">
        <div
          v-for="item in statItems"
          :key="item.key"
          class="flex flex-col items-center gap-2 p-4 rounded-lg bg-white/[0.02] border border-white/5 hover:border-white/10 transition-colors"
        >
          <component
            :is="item.icon"
            class="w-8 h-8"
            :class="item.color"
          />
          <span class="text-lg font-semibold text-white tabular-nums">
            {{ item.format(displayValues[item.key] ?? 0) }}
          </span>
          <span class="text-xs text-gray-500 text-center">{{ item.label }}</span>
        </div>
      </div>
    </div>
  </section>
</template>
