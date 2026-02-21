<script setup lang="ts">
import { ref, computed } from 'vue'
import { ShoppingBag, ChevronLeft, ChevronRight } from 'lucide-vue-next'
import type { Purchase } from '@/stores/userStore'
import { formatPrice } from '@/utils/currency'

const props = defineProps<{
  purchases: Purchase[]
}>()

const currentPage = ref(1)
const itemsPerPage = 10

const totalPages = computed(() => Math.ceil(props.purchases.length / itemsPerPage))

const paginatedPurchases = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage
  const end = start + itemsPerPage
  return props.purchases.slice(start, end)
})

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

function getStatusClasses(status: Purchase['status']): string {
  switch (status) {
    case 'completed':
      return 'bg-green-500/20 text-green-400 border-green-500/30'
    case 'pending':
      return 'bg-yellow-500/20 text-yellow-400 border-yellow-500/30'
    case 'refunded':
      return 'bg-red-500/20 text-red-400 border-red-500/30'
    default:
      return 'bg-gray-500/20 text-gray-400 border-gray-500/30'
  }
}

function prevPage() {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}
</script>

<template>
  <div class="card">
    <div class="p-6 border-b border-gray-700">
      <h3 class="text-xl font-bold text-white flex items-center gap-2">
        <ShoppingBag class="w-5 h-5 text-weenie-gold" />
        Purchase History
      </h3>
    </div>

    <!-- Empty State -->
    <div v-if="purchases.length === 0" class="p-12 text-center">
      <ShoppingBag class="w-12 h-12 text-gray-600 mx-auto mb-4" />
      <h4 class="text-lg font-semibold text-white mb-2">No purchases yet</h4>
      <p class="text-gray-400 mb-6">
        You haven't made any purchases. Visit our store to support the server!
      </p>
      <RouterLink
        to="/store"
        class="inline-flex items-center gap-2 px-6 py-3 bg-weenie-gradient text-white font-semibold rounded-lg hover:opacity-90 transition-opacity"
      >
        <ShoppingBag class="w-5 h-5" />
        Browse Store
      </RouterLink>
    </div>

    <!-- Purchase Table -->
    <div v-else class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b border-gray-700">
            <th class="text-left text-gray-400 text-sm font-medium px-6 py-4">Date</th>
            <th class="text-left text-gray-400 text-sm font-medium px-6 py-4">Package</th>
            <th class="text-right text-gray-400 text-sm font-medium px-6 py-4">Price</th>
            <th class="text-center text-gray-400 text-sm font-medium px-6 py-4">Status</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="purchase in paginatedPurchases"
            :key="purchase.id"
            class="border-b border-gray-700/50 hover:bg-white/5 transition-colors"
          >
            <td class="px-6 py-4 text-gray-300 text-sm">
              {{ formatDate(purchase.date) }}
            </td>
            <td class="px-6 py-4 text-white font-medium">
              {{ purchase.packageName }}
            </td>
            <td class="px-6 py-4 text-right text-weenie-gold font-semibold">
              {{ formatPrice(purchase.price, 'USD') }}
            </td>
            <td class="px-6 py-4 text-center">
              <span
                :class="getStatusClasses(purchase.status)"
                class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium border capitalize"
              >
                {{ purchase.status }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Pagination -->
      <div
        v-if="totalPages > 1"
        class="flex items-center justify-between px-6 py-4 border-t border-gray-700"
      >
        <p class="text-gray-400 text-sm">
          Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to
          {{ Math.min(currentPage * itemsPerPage, purchases.length) }} of
          {{ purchases.length }} purchases
        </p>

        <div class="flex items-center gap-2">
          <button
            @click="prevPage"
            :disabled="currentPage === 1"
            class="p-2 rounded-lg text-gray-400 hover:text-white hover:bg-white/5 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
          >
            <ChevronLeft class="w-5 h-5" />
          </button>

          <span class="text-white px-3">
            {{ currentPage }} / {{ totalPages }}
          </span>

          <button
            @click="nextPage"
            :disabled="currentPage === totalPages"
            class="p-2 rounded-lg text-gray-400 hover:text-white hover:bg-white/5 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
          >
            <ChevronRight class="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
