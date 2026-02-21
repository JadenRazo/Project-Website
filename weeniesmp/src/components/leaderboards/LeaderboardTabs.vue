<script setup lang="ts">
import type { Component } from 'vue'
import type { LeaderboardCategory } from '@/stores/leaderboardStore'

interface Category {
  id: LeaderboardCategory
  name: string
  icon: Component
}

defineProps<{
  categories: Category[]
  activeCategory: LeaderboardCategory
}>()

const emit = defineEmits<{
  (e: 'change', category: LeaderboardCategory): void
}>()
</script>

<template>
  <div class="flex flex-wrap justify-center gap-2 mb-8">
    <button
      v-for="category in categories"
      :key="category.id"
      @click="emit('change', category.id)"
      class="flex items-center gap-2 px-3 sm:px-6 py-2 sm:py-3 rounded-lg font-medium transition-all duration-300"
      :class="
        activeCategory === category.id
          ? 'bg-weenie-gradient text-white shadow-lg shadow-weenie-red/30'
          : 'bg-weenie-dark/50 text-gray-400 hover:text-white hover:bg-weenie-dark'
      "
    >
      <component :is="category.icon" class="w-4 h-4 sm:w-5 sm:h-5" />
      {{ category.name }}
    </button>
  </div>
</template>
