<script setup lang="ts">
import { ref } from 'vue'
import { Star, Sword, Gem, Clock, Hammer, Box, Lock } from 'lucide-vue-next'
import { formatDate } from '@/utils/formatters'
import type { Achievement } from '@/stores/playerStore'

interface Props {
  achievements: Achievement[]
}

defineProps<Props>()

const hoveredAchievement = ref<string | null>(null)

const iconMap: Record<string, typeof Star> = {
  star: Star,
  sword: Sword,
  gem: Gem,
  clock: Clock,
  hammer: Hammer,
  blocks: Box
}

function getIcon(iconName: string) {
  return iconMap[iconName] || Star
}
</script>

<template>
  <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
    <div
      v-for="achievement in achievements"
      :key="achievement.id"
      class="relative"
      @mouseenter="hoveredAchievement = achievement.id"
      @mouseleave="hoveredAchievement = null"
    >
      <!-- Achievement Badge -->
      <div
        class="flex flex-col items-center p-4 rounded-xl border transition-all duration-200"
        :class="achievement.unlockedAt
          ? 'bg-weenie-darker/50 border-weenie-gold/30 hover:border-weenie-gold/60'
          : 'bg-weenie-darker/30 border-white/5 opacity-50'"
      >
        <!-- Icon -->
        <div
          class="w-12 h-12 rounded-full flex items-center justify-center mb-3"
          :class="achievement.unlockedAt
            ? 'bg-weenie-gradient'
            : 'bg-gray-700'"
        >
          <component
            v-if="achievement.unlockedAt"
            :is="getIcon(achievement.icon)"
            class="w-6 h-6 text-white"
          />
          <Lock v-else class="w-6 h-6 text-gray-400" />
        </div>

        <!-- Name -->
        <h4
          class="text-sm font-medium text-center"
          :class="achievement.unlockedAt ? 'text-white' : 'text-gray-500'"
        >
          {{ achievement.name }}
        </h4>

        <!-- Locked indicator -->
        <span
          v-if="!achievement.unlockedAt"
          class="text-xs text-gray-500 mt-1"
        >
          Locked
        </span>
      </div>

      <!-- Tooltip -->
      <Transition
        enter-active-class="transition duration-150 ease-out"
        enter-from-class="opacity-0 translate-y-1"
        enter-to-class="opacity-100 translate-y-0"
        leave-active-class="transition duration-100 ease-in"
        leave-from-class="opacity-100 translate-y-0"
        leave-to-class="opacity-0 translate-y-1"
      >
        <div
          v-if="hoveredAchievement === achievement.id"
          class="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 z-10 w-48 p-3 bg-weenie-dark rounded-lg border border-white/10 shadow-xl"
        >
          <h5 class="font-semibold text-white text-sm mb-1">{{ achievement.name }}</h5>
          <p class="text-xs text-gray-400 mb-2">{{ achievement.description }}</p>
          <p v-if="achievement.unlockedAt" class="text-xs text-weenie-gold">
            Unlocked {{ formatDate(achievement.unlockedAt) }}
          </p>
          <p v-else class="text-xs text-gray-500">
            Not yet unlocked
          </p>
          <!-- Tooltip arrow -->
          <div class="absolute top-full left-1/2 -translate-x-1/2 w-0 h-0 border-l-8 border-r-8 border-t-8 border-transparent border-t-weenie-dark"></div>
        </div>
      </Transition>
    </div>
  </div>
</template>
