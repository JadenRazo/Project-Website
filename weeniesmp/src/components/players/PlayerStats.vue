<script setup lang="ts">
import { Sword, Shield, Box, Pickaxe, Bug, MapPin } from 'lucide-vue-next'
import { formatNumber } from '@/utils/formatters'
import type { PlayerStats } from '@/stores/playerStore'

interface Props {
  stats: PlayerStats
}

defineProps<Props>()

const statItems = [
  { key: 'kills', label: 'Player Kills', icon: Sword, color: 'text-red-400' },
  { key: 'deaths', label: 'Deaths', icon: Shield, color: 'text-gray-400' },
  { key: 'blocksPlaced', label: 'Blocks Placed', icon: Box, color: 'text-blue-400' },
  { key: 'blocksBroken', label: 'Blocks Broken', icon: Pickaxe, color: 'text-yellow-400' },
  { key: 'mobsKilled', label: 'Mobs Killed', icon: Bug, color: 'text-green-400' },
  { key: 'distanceTraveled', label: 'Distance Traveled', icon: MapPin, color: 'text-purple-400' }
] as const
</script>

<template>
  <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
    <div
      v-for="stat in statItems"
      :key="stat.key"
      class="bg-weenie-darker/50 rounded-xl p-4 border border-white/5"
    >
      <div class="flex items-center gap-2 mb-2">
        <component :is="stat.icon" class="w-5 h-5" :class="stat.color" />
        <span class="text-sm text-gray-400">{{ stat.label }}</span>
      </div>
      <p class="text-2xl font-bold text-white">
        {{ formatNumber(stats[stat.key]) }}
        <span v-if="stat.key === 'distanceTraveled'" class="text-sm font-normal text-gray-500">blocks</span>
      </p>
    </div>
  </div>
</template>
