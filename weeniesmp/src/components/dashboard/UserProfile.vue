<script setup lang="ts">
import { computed } from 'vue'
import { Calendar, ShoppingBag, DollarSign } from 'lucide-vue-next'
import type { User } from '@/stores/userStore'
import { formatPrice } from '@/utils/currency'

const props = defineProps<{
  user: User
  totalSpent: number
  purchaseCount: number
}>()

const avatarUrl = computed(() => {
  if (!props.user.uuid) {
    return `https://crafatar.com/avatars/MHF_Steve?size=128&overlay`
  }
  return `https://crafatar.com/avatars/${props.user.uuid}?size=128&overlay`
})

const bodyUrl = computed(() => {
  if (!props.user.uuid) {
    return `https://crafatar.com/renders/body/MHF_Steve?size=128&overlay`
  }
  return `https://crafatar.com/renders/body/${props.user.uuid}?size=128&overlay`
})

const formattedJoinDate = computed(() => {
  if (!props.user.joinDate) return 'Unknown'
  return new Date(props.user.joinDate).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
})
</script>

<template>
  <div class="card p-6 md:p-8">
    <div class="flex flex-col md:flex-row items-center md:items-start gap-6">
      <!-- Avatar -->
      <div class="relative group">
        <div class="w-32 h-32 rounded-xl overflow-hidden bg-weenie-darker border-2 border-weenie-gold/30 group-hover:border-weenie-gold transition-colors">
          <img
            :src="avatarUrl"
            :alt="user.username"
            class="w-full h-full object-cover"
            loading="lazy"
          />
        </div>
        <div class="absolute -bottom-2 -right-2 w-8 h-8 bg-green-500 rounded-full border-4 border-weenie-dark" title="Online"></div>
      </div>

      <!-- Info -->
      <div class="flex-1 text-center md:text-left">
        <h2 class="text-3xl font-bold text-white mb-2">{{ user.username }}</h2>

        <div class="flex flex-wrap justify-center md:justify-start gap-4 text-gray-400 text-sm mb-6">
          <div v-if="user.joinDate" class="flex items-center gap-2">
            <Calendar class="w-4 h-4" />
            <span>Joined {{ formattedJoinDate }}</span>
          </div>
        </div>

        <!-- Stats -->
        <div class="grid grid-cols-2 gap-4">
          <div class="bg-weenie-darker rounded-lg p-4">
            <div class="flex items-center gap-2 text-gray-400 text-sm mb-1">
              <DollarSign class="w-4 h-4" />
              <span>Total Spent</span>
            </div>
            <p class="text-2xl font-bold text-weenie-gold">
              {{ formatPrice(totalSpent, 'USD') }}
            </p>
          </div>

          <div class="bg-weenie-darker rounded-lg p-4">
            <div class="flex items-center gap-2 text-gray-400 text-sm mb-1">
              <ShoppingBag class="w-4 h-4" />
              <span>Purchases</span>
            </div>
            <p class="text-2xl font-bold text-white">
              {{ purchaseCount }}
            </p>
          </div>
        </div>
      </div>

      <!-- 3D Body Render (hidden on mobile) -->
      <div class="hidden lg:block">
        <img
          :src="bodyUrl"
          :alt="user.username"
          class="h-40 object-contain opacity-80 hover:opacity-100 transition-opacity"
          loading="lazy"
        />
      </div>
    </div>
  </div>
</template>
