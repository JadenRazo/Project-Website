<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { ShoppingBag, User, LogOut, Search, X, ArrowUpDown } from 'lucide-vue-next'
import CategoryTabs from '@/components/store/CategoryTabs.vue'
import ProductCard from '@/components/store/ProductCard.vue'
import ProductCardSkeleton from '@/components/store/ProductCardSkeleton.vue'
import PackageModal from '@/components/store/PackageModal.vue'
import { useTebexStore, type TebexPackage } from '@/stores/tebexStore'

const tebexStore = useTebexStore()
const selectedPackage = ref<TebexPackage | null>(null)
const modalOpen = ref(false)
const searchQuery = ref('')
const sortOption = ref<'default' | 'price-asc' | 'price-desc' | 'name'>('default')

function openPackageModal(pkg: TebexPackage) {
  selectedPackage.value = pkg
  modalOpen.value = true
}

function closeModal() {
  modalOpen.value = false
}
const activeCategory = ref<number | null>(null)
const usernameInput = ref('')

function handleLogin() {
  if (usernameInput.value.trim()) {
    const success = tebexStore.setUsername(usernameInput.value.trim())
    if (success) {
      usernameInput.value = ''
    }
  }
}

const filteredPackages = computed(() => {
  let result = tebexStore.packages

  // Filter by category
  if (activeCategory.value) {
    result = result.filter(p => p.category?.id === activeCategory.value)
  }

  // Filter by search query
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(p =>
      p.name.toLowerCase().includes(query) ||
      p.description.toLowerCase().includes(query)
    )
  }

  return result
})

const sortedPackages = computed(() => {
  const result = [...filteredPackages.value]

  switch (sortOption.value) {
    case 'price-asc':
      return result.sort((a, b) => a.total_price - b.total_price)
    case 'price-desc':
      return result.sort((a, b) => b.total_price - a.total_price)
    case 'name':
      return result.sort((a, b) => a.name.localeCompare(b.name))
    default:
      return result
  }
})

onMounted(async () => {
  await tebexStore.fetchCategories()
  if (tebexStore.categories.length > 0) {
    activeCategory.value = tebexStore.categories[0].id
  }
})
</script>

<template>
  <div class="min-h-screen pt-24 pb-16 bg-weenie-darker">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="text-center mb-12">
        <ShoppingBag class="w-12 h-12 text-weenie-gold mx-auto mb-4" />
        <h1 class="text-4xl md:text-5xl font-bold mb-4">
          <span class="gradient-text">Store</span>
        </h1>
        <p class="text-gray-400 max-w-xl mx-auto">
          Support the server and get awesome perks! All purchases are delivered instantly.
        </p>
      </div>

      <!-- Username Login -->
      <div class="max-w-md mx-auto mb-8">
        <div v-if="!tebexStore.username" class="card p-6">
          <div class="flex items-center gap-3 mb-4">
            <User class="w-6 h-6 text-weenie-gold" />
            <h3 class="text-lg font-semibold text-white">Enter your Minecraft username</h3>
          </div>
          <p class="text-gray-400 text-sm mb-2">Required to make purchases</p>
          <p class="text-yellow-400/80 text-xs mb-4">You must be online on the server with this exact username to receive purchases.</p>
          <form @submit.prevent="handleLogin" class="flex gap-2">
            <input
              v-model="usernameInput"
              type="text"
              placeholder="Your Minecraft username"
              class="flex-1 px-4 py-2 bg-weenie-darker border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-weenie-gold"
            />
            <button
              type="submit"
              class="px-6 py-2 bg-weenie-gradient text-white font-semibold rounded-lg hover:opacity-90 transition-opacity"
            >
              Continue
            </button>
          </form>
        </div>
        <div v-else class="flex items-center justify-center gap-4 p-4 bg-weenie-dark rounded-lg">
          <User class="w-5 h-5 text-weenie-gold" />
          <span class="text-white">Shopping as <strong class="text-weenie-gold">{{ tebexStore.username }}</strong></span>
          <button
            @click="tebexStore.logout()"
            class="flex items-center gap-1 text-gray-400 hover:text-red-400 transition-colors text-sm"
          >
            <LogOut class="w-4 h-4" />
            Change
          </button>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="tebexStore.loading">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <ProductCardSkeleton v-for="i in 6" :key="i" />
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="tebexStore.error" class="text-center py-20">
        <p class="text-red-400 mb-4">{{ tebexStore.error }}</p>
        <button
          @click="tebexStore.fetchCategories()"
          class="px-6 py-2 bg-weenie-red text-white rounded-lg hover:bg-weenie-red/80 transition-colors"
        >
          Retry
        </button>
      </div>

      <!-- Store Content -->
      <template v-else>
        <!-- Search & Sort Controls -->
        <div class="flex flex-col sm:flex-row gap-4 mb-6">
          <!-- Search Box -->
          <div class="relative flex-1 max-w-md">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-500" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search packages..."
              class="w-full pl-10 pr-10 py-3 bg-weenie-dark border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-weenie-gold transition-colors"
            />
            <button
              v-if="searchQuery"
              @click="searchQuery = ''"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-white transition-colors"
            >
              <X class="w-5 h-5" />
            </button>
          </div>

          <!-- Sort Dropdown -->
          <div class="relative">
            <ArrowUpDown class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500 pointer-events-none" />
            <select
              v-model="sortOption"
              class="pl-10 pr-8 py-3 bg-weenie-dark border border-gray-700 rounded-lg text-white appearance-none cursor-pointer focus:outline-none focus:border-weenie-gold transition-colors"
            >
              <option value="default">Default</option>
              <option value="price-asc">Price: Low to High</option>
              <option value="price-desc">Price: High to Low</option>
              <option value="name">Name: A-Z</option>
            </select>
          </div>
        </div>

        <CategoryTabs
          v-if="tebexStore.categories.length > 0"
          :categories="tebexStore.categories"
          :active-category="activeCategory"
          @change="activeCategory = $event"
        />

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <ProductCard
            v-for="pkg in sortedPackages"
            :key="pkg.id"
            :package="pkg"
            @select="openPackageModal"
          />
        </div>

        <div v-if="sortedPackages.length === 0" class="text-center py-12">
          <p class="text-gray-400">
            {{ searchQuery ? `No packages found matching "${searchQuery}"` : 'No packages available in this category.' }}
          </p>
          <button
            v-if="searchQuery"
            @click="searchQuery = ''"
            class="mt-4 px-4 py-2 text-weenie-gold border border-weenie-gold rounded-lg hover:bg-weenie-gold/10 transition-colors"
          >
            Clear search
          </button>
        </div>
      </template>

      <div class="mt-12 text-center text-sm text-gray-500">
        <p>All purchases are processed securely through Tebex.</p>
        <p class="mt-1">Items are delivered automatically within minutes.</p>
      </div>
    </div>

    <!-- Package Detail Modal -->
    <PackageModal
      :package="selectedPackage"
      :open="modalOpen"
      @close="closeModal"
    />
  </div>
</template>
