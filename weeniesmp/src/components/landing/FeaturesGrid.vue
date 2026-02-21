<script setup lang="ts">
import {
  Coins,
  Briefcase,
  MapPin,
  Store,
  Backpack,
  Castle,
  Mountain,
  Wheat,
  Swords
} from 'lucide-vue-next'

const features = [
  { icon: Coins, title: 'Economy', color: 'text-yellow-400' },
  { icon: Briefcase, title: 'Jobs', color: 'text-blue-400' },
  { icon: MapPin, title: 'Claims', color: 'text-green-400' },
  { icon: Store, title: 'Auctions', color: 'text-purple-400' },
  { icon: Backpack, title: 'Backpacks', color: 'text-orange-400' },
  { icon: Castle, title: 'Dungeons', color: 'text-red-400' },
  { icon: Mountain, title: 'Mines', color: 'text-amber-400' },
  { icon: Wheat, title: 'Farms', color: 'text-lime-400' },
  { icon: Swords, title: 'PvP Arena', color: 'text-pink-400' }
]
</script>

<template>
  <section class="py-6 bg-weenie-darker">
    <div class="max-w-4xl mx-auto px-6">
      <h2 class="text-xl font-semibold text-center text-white mb-5">Features</h2>

      <div class="grid grid-cols-3 sm:grid-cols-5 lg:grid-cols-9 gap-4">
        <div
          v-for="feature in features"
          :key="feature.title"
          class="flex flex-col items-center gap-2 p-3 rounded-lg bg-white/[0.02] border border-white/5 hover:border-white/10 transition-colors"
        >
          <component
            :is="feature.icon"
            class="w-6 h-6"
            :class="feature.color"
          />
          <span class="text-xs text-gray-400 text-center">{{ feature.title }}</span>
        </div>
      </div>
    </div>
  </section>
</template>
