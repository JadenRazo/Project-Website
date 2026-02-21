<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { LayoutDashboard, LogOut, Settings, RefreshCw, Loader2 } from 'lucide-vue-next'
import { useUserStore } from '@/stores/userStore'
import UserProfile from '@/components/dashboard/UserProfile.vue'
import PurchaseHistory from '@/components/dashboard/PurchaseHistory.vue'
import LoginPrompt from '@/components/dashboard/LoginPrompt.vue'

const router = useRouter()
const userStore = useUserStore()

const isAuthenticated = computed(() => userStore.isAuthenticated)

function handleLogout() {
  userStore.logout()
  router.push('/')
}

async function refreshPurchases() {
  await userStore.fetchPurchases()
}
</script>

<template>
  <div class="min-h-screen pt-24 pb-16 bg-weenie-darker">
    <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="text-center mb-12">
        <LayoutDashboard class="w-12 h-12 text-weenie-gold mx-auto mb-4" />
        <h1 class="text-4xl md:text-5xl font-bold mb-4">
          <span class="gradient-text">Dashboard</span>
        </h1>
        <p class="text-gray-400 max-w-xl mx-auto">
          View your account details and purchase history.
        </p>
      </div>

      <!-- Login Prompt (if not authenticated) -->
      <LoginPrompt v-if="!isAuthenticated" />

      <!-- Dashboard Content (if authenticated) -->
      <template v-else>
        <!-- User Profile -->
        <UserProfile
          :user="userStore.user!"
          :total-spent="userStore.totalSpent"
          :purchase-count="userStore.purchaseCount"
          class="mb-8"
        />

        <!-- Actions Bar -->
        <div class="flex flex-wrap gap-4 mb-8">
          <button
            @click="refreshPurchases"
            :disabled="userStore.loading"
            class="flex items-center gap-2 px-4 py-2 bg-weenie-dark border border-gray-700 text-gray-300 rounded-lg hover:text-white hover:border-weenie-gold transition-all disabled:opacity-50"
          >
            <RefreshCw v-if="!userStore.loading" class="w-4 h-4" />
            <Loader2 v-else class="w-4 h-4 animate-spin" />
            Refresh
          </button>

          <RouterLink
            to="/store"
            class="flex items-center gap-2 px-4 py-2 bg-weenie-gradient text-white font-semibold rounded-lg hover:opacity-90 transition-opacity"
          >
            Visit Store
          </RouterLink>

          <button
            @click="handleLogout"
            class="flex items-center gap-2 px-4 py-2 text-red-400 border border-red-400/30 rounded-lg hover:bg-red-400/10 transition-all ml-auto"
          >
            <LogOut class="w-4 h-4" />
            Sign Out
          </button>
        </div>

        <!-- Purchase History -->
        <PurchaseHistory :purchases="userStore.purchases" class="mb-8" />

        <!-- Account Settings -->
        <div class="card p-6">
          <div class="flex items-center gap-3 mb-6">
            <Settings class="w-5 h-5 text-weenie-gold" />
            <h3 class="text-xl font-bold text-white">Account Settings</h3>
          </div>

          <div class="space-y-4">
            <div class="flex flex-col sm:flex-row sm:items-center justify-between p-4 bg-weenie-darker rounded-lg">
              <div>
                <p class="text-white font-medium">Minecraft Username</p>
                <p class="text-gray-400 text-sm">{{ userStore.user?.username }}</p>
              </div>
              <button
                @click="handleLogout"
                class="mt-3 sm:mt-0 text-sm text-weenie-gold hover:underline"
              >
                Change Account
              </button>
            </div>

            <div class="flex flex-col sm:flex-row sm:items-center justify-between p-4 bg-weenie-darker rounded-lg">
              <div>
                <p class="text-white font-medium">UUID</p>
                <p class="text-gray-400 text-sm font-mono">{{ userStore.user?.uuid || 'Unknown' }}</p>
              </div>
            </div>
          </div>

          <div class="mt-6 pt-6 border-t border-gray-700">
            <p class="text-gray-500 text-sm">
              Need help? Contact us on
              <a
                href="https://discord.com/invite/weeniesmp"
                target="_blank"
                rel="noopener noreferrer"
                class="text-weenie-gold hover:underline"
              >
                Discord
              </a>
              for support.
            </p>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>
