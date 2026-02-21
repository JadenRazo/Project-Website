<script setup lang="ts">
import { Crown, Sparkles, Package, Coins, Tag } from 'lucide-vue-next'
import type { TebexCategory } from '@/stores/tebexStore'

defineProps<{
  categories: TebexCategory[]
  activeCategory: number | null
}>()

const emit = defineEmits<{
  (e: 'change', categoryId: number): void
}>()

const getIcon = (name: string) => {
  const lower = name.toLowerCase()
  if (lower.includes('rank')) return Crown
  if (lower.includes('cosmetic')) return Sparkles
  if (lower.includes('crate') || lower.includes('key')) return Package
  if (lower.includes('coin') || lower.includes('currency')) return Coins
  return Tag
}
</script>

<template>
  <div class="flex flex-wrap justify-center gap-2 mb-8">
    <button
      v-for="category in categories"
      :key="category.id"
      @click="emit('change', category.id)"
      class="flex items-center gap-2 px-6 py-3 rounded-lg font-medium transition-all duration-300"
      :class="
        activeCategory === category.id
          ? 'bg-weenie-gradient text-white shadow-lg shadow-weenie-red/30'
          : 'bg-weenie-dark/50 text-gray-400 hover:text-white hover:bg-weenie-dark'
      "
    >
      <component :is="getIcon(category.name)" class="w-5 h-5" />
      {{ category.name }}
    </button>
  </div>
</template>
