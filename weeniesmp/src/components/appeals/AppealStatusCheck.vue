<script setup lang="ts">
import { ref } from 'vue'
import { Search, Loader2, ArrowLeft, MessageSquare, User, Calendar } from 'lucide-vue-next'
import { useAppealStore } from '@/stores/appealStore'
import AppealStatusBadge from './AppealStatusBadge.vue'
import AppealTimeline from './AppealTimeline.vue'

const appealStore = useAppealStore()

const appealId = ref('')
const email = ref('')

async function handleCheck() {
  if (!appealId.value.trim() || !email.value.trim() || appealStore.loading) return
  await appealStore.checkAppealStatus(appealId.value.trim(), email.value.trim())
}

function goBack() {
  appealStore.clearCurrentAppeal()
  appealId.value = ''
  email.value = ''
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function getBanReasonLabel(reason: string): string {
  const reasons: Record<string, string> = {
    cheating: 'Cheating / Hacks',
    griefing: 'Griefing',
    harassment: 'Harassment / Toxic Behavior',
    scamming: 'Scamming',
    exploiting: 'Exploiting Bugs',
    advertising: 'Advertising',
    other: 'Other / Unknown'
  }
  return reasons[reason] || reason
}
</script>

<template>
  <!-- Appeal Details View -->
  <div v-if="appealStore.currentAppeal">
    <button
      @click="goBack"
      class="flex items-center gap-2 text-gray-400 hover:text-white transition-colors mb-6"
    >
      <ArrowLeft class="w-4 h-4" />
      Check another appeal
    </button>

    <!-- Appeal Header -->
    <div class="bg-weenie-dark rounded-lg p-6 mb-6">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-4">
        <div>
          <p class="text-sm text-gray-400 mb-1">Appeal ID</p>
          <code class="text-lg font-mono text-weenie-gold">{{ appealStore.currentAppeal.id }}</code>
        </div>
        <AppealStatusBadge :status="appealStore.currentAppeal.status" />
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 pt-4 border-t border-gray-700">
        <div class="flex items-center gap-3">
          <User class="w-5 h-5 text-gray-500" />
          <div>
            <p class="text-xs text-gray-500">Username</p>
            <p class="text-white font-medium">{{ appealStore.currentAppeal.username }}</p>
          </div>
        </div>
        <div class="flex items-center gap-3">
          <MessageSquare class="w-5 h-5 text-gray-500" />
          <div>
            <p class="text-xs text-gray-500">Ban Reason</p>
            <p class="text-white font-medium">{{ getBanReasonLabel(appealStore.currentAppeal.banReason) }}</p>
          </div>
        </div>
        <div class="flex items-center gap-3">
          <Calendar class="w-5 h-5 text-gray-500" />
          <div>
            <p class="text-xs text-gray-500">Submitted</p>
            <p class="text-white font-medium text-sm">{{ formatDate(appealStore.currentAppeal.createdAt) }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Timeline -->
    <div class="bg-weenie-dark rounded-lg p-6 mb-6">
      <h3 class="text-lg font-semibold text-white mb-6">Appeal Progress</h3>
      <AppealTimeline :appeal="appealStore.currentAppeal" />
    </div>

    <!-- Your Appeal Text -->
    <div class="bg-weenie-dark rounded-lg p-6 mb-6">
      <h3 class="text-lg font-semibold text-white mb-4">Your Appeal</h3>
      <p class="text-gray-300 whitespace-pre-wrap">{{ appealStore.currentAppeal.appealText }}</p>
    </div>

    <!-- Staff Response -->
    <div
      v-if="appealStore.currentAppeal.staffResponse"
      class="bg-weenie-dark rounded-lg p-6 border-l-4"
      :class="appealStore.currentAppeal.status === 'approved' ? 'border-green-500' : appealStore.currentAppeal.status === 'denied' ? 'border-red-500' : 'border-blue-500'"
    >
      <h3 class="text-lg font-semibold text-white mb-4">Staff Response</h3>
      <p class="text-gray-300 whitespace-pre-wrap">{{ appealStore.currentAppeal.staffResponse }}</p>
    </div>

    <!-- No Response Yet -->
    <div
      v-else-if="appealStore.currentAppeal.status === 'pending' || appealStore.currentAppeal.status === 'under_review'"
      class="bg-weenie-dark rounded-lg p-6 text-center"
    >
      <MessageSquare class="w-10 h-10 text-gray-600 mx-auto mb-3" />
      <p class="text-gray-400">No staff response yet. Please check back later.</p>
      <p class="text-sm text-gray-500 mt-2">Appeals are typically reviewed within 24-48 hours.</p>
    </div>
  </div>

  <!-- Search Form -->
  <form v-else @submit.prevent="handleCheck" class="space-y-6">
    <p class="text-gray-400 mb-6">
      Enter your appeal ID and the email address you used when submitting your appeal.
    </p>

    <!-- Appeal ID Field -->
    <div>
      <label for="checkAppealId" class="block text-sm font-medium text-gray-300 mb-2">
        Appeal ID
      </label>
      <input
        id="checkAppealId"
        v-model="appealId"
        type="text"
        placeholder="e.g., WS-ABC123-XYZ9"
        class="w-full px-4 py-3 bg-weenie-dark border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-weenie-gold/50 focus:border-weenie-gold transition-all font-mono"
      />
    </div>

    <!-- Email Field -->
    <div>
      <label for="checkEmail" class="block text-sm font-medium text-gray-300 mb-2">
        Email Address
      </label>
      <input
        id="checkEmail"
        v-model="email"
        type="email"
        placeholder="your@email.com"
        class="w-full px-4 py-3 bg-weenie-dark border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-weenie-gold/50 focus:border-weenie-gold transition-all"
      />
      <p class="mt-1.5 text-sm text-gray-500">The email you used when submitting your appeal</p>
    </div>

    <!-- Error Message -->
    <div
      v-if="appealStore.error"
      class="bg-red-500/10 border border-red-500/30 rounded-lg p-4 text-red-400"
    >
      {{ appealStore.error }}
    </div>

    <!-- Submit Button -->
    <button
      type="submit"
      :disabled="!appealId.trim() || !email.trim() || appealStore.loading"
      class="btn-primary w-full flex items-center justify-center gap-2"
    >
      <Loader2 v-if="appealStore.loading" class="w-5 h-5 animate-spin" />
      <Search v-else class="w-5 h-5" />
      {{ appealStore.loading ? 'Checking...' : 'Check Status' }}
    </button>
  </form>
</template>
