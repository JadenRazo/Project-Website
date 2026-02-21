<script setup lang="ts">
import { ShoppingCart, Check, Crown, Sparkles, Package, Coins, Tag, Loader2 } from 'lucide-vue-next'
import { ref } from 'vue'
import { useTebexStore, type TebexPackage } from '@/stores/tebexStore'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
  package: TebexPackage
}>()

const emit = defineEmits<{
  select: [pkg: TebexPackage]
}>()

const tebexStore = useTebexStore()
const adding = ref(false)
const added = ref(false)

const getIcon = (categoryName: string) => {
  const lower = categoryName?.toLowerCase() || ''
  if (lower.includes('rank')) return Crown
  if (lower.includes('cosmetic')) return Sparkles
  if (lower.includes('crate') || lower.includes('key')) return Package
  if (lower.includes('coin') || lower.includes('currency')) return Coins
  return Tag
}

async function addToCart() {
  adding.value = true
  const success = await tebexStore.addToBasket(props.package.id)
  adding.value = false

  if (success) {
    added.value = true
    setTimeout(() => {
      added.value = false
    }, 1500)
  }
}

function stripHtml(html: string) {
  const doc = new DOMParser().parseFromString(html, 'text/html')
  return doc.body.textContent || ''
}
</script>

<template>
  <div class="card group hover:scale-[1.02] flex flex-col cursor-pointer relative" @click="emit('select', package)">
    <!-- Sale Badge -->
    <div
      v-if="package.sales_price && package.sales_price < package.base_price"
      class="absolute -top-2 -right-2 px-2 py-1 bg-weenie-red text-white text-xs font-bold rounded-full z-10"
    >
      SALE
    </div>

    <div class="flex items-start justify-between mb-4">
      <div
        class="w-36 h-20 rounded-xl bg-weenie-gradient flex items-center justify-center overflow-hidden"
      >
        <img
          v-if="package.image"
          :src="package.image"
          :alt="package.name"
          class="w-full h-full object-cover"
        />
        <component
          v-else
          :is="getIcon(package.category?.name || '')"
          class="w-10 h-10 text-white"
        />
      </div>
      <div class="text-right">
        <span
          v-if="package.sales_price && package.sales_price < package.base_price"
          class="text-sm text-gray-500 line-through block"
        >
          {{ formatCurrency(package.currency) }}{{ package.base_price.toFixed(2) }}
        </span>
        <span class="text-2xl font-bold text-weenie-gold">
          {{ formatCurrency(package.currency) }}{{ package.total_price.toFixed(2) }}
        </span>
      </div>
    </div>

    <h3 class="text-xl font-bold text-white mb-2">{{ package.name }}</h3>
    <p class="text-gray-400 text-sm mb-4 flex-1">
      {{ stripHtml(package.description).slice(0, 150) }}{{ stripHtml(package.description).length > 150 ? '...' : '' }}
    </p>

    <button
      @click.stop="addToCart"
      :disabled="adding"
      class="w-full py-3 rounded-lg font-semibold transition-all duration-300 flex items-center justify-center gap-2 disabled:opacity-50"
      :class="
        added
          ? 'bg-green-500 text-white'
          : 'bg-weenie-dark hover:bg-weenie-gradient text-gray-300 hover:text-white'
      "
    >
      <Loader2 v-if="adding" class="w-5 h-5 animate-spin" />
      <Check v-else-if="added" class="w-5 h-5" />
      <ShoppingCart v-else class="w-5 h-5" />
      {{ adding ? 'Adding...' : added ? 'Added!' : 'Add to Cart' }}
    </button>
  </div>
</template>
