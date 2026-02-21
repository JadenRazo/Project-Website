<script setup lang="ts">
import { ref, computed } from 'vue'
import { Send, Loader2, CheckCircle, Copy, AlertCircle } from 'lucide-vue-next'
import { useAppealStore } from '@/stores/appealStore'
import { useToastStore } from '@/stores/toastStore'
import type { Appeal } from '@/stores/appealStore'

const appealStore = useAppealStore()
const toastStore = useToastStore()

const username = ref('')
const email = ref('')
const banReason = ref('')
const appealText = ref('')
const agreedToTerms = ref(false)

const submittedAppeal = ref<Appeal | null>(null)

const banReasons = [
  { value: 'cheating', label: 'Cheating / Hacks' },
  { value: 'griefing', label: 'Griefing' },
  { value: 'harassment', label: 'Harassment / Toxic Behavior' },
  { value: 'scamming', label: 'Scamming' },
  { value: 'exploiting', label: 'Exploiting Bugs' },
  { value: 'advertising', label: 'Advertising' },
  { value: 'other', label: 'Other / Unknown' }
]

// Validation
const usernameError = computed(() => {
  if (!username.value) return ''
  if (username.value.length < 3) return 'Username must be at least 3 characters'
  if (username.value.length > 16) return 'Username must be at most 16 characters'
  if (!/^[a-zA-Z0-9_]+$/.test(username.value)) return 'Username can only contain letters, numbers, and underscores'
  return ''
})

const emailError = computed(() => {
  if (!email.value) return ''
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(email.value)) return 'Please enter a valid email address'
  return ''
})

const appealTextLength = computed(() => appealText.value.length)

const appealTextError = computed(() => {
  if (!appealText.value) return ''
  if (appealTextLength.value < 100) return `Please write at least 100 characters (${100 - appealTextLength.value} more needed)`
  if (appealTextLength.value > 2000) return `Maximum 2000 characters (${appealTextLength.value - 2000} over)`
  return ''
})

const isFormValid = computed(() => {
  return (
    username.value.length >= 3 &&
    username.value.length <= 16 &&
    /^[a-zA-Z0-9_]+$/.test(username.value) &&
    /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value) &&
    banReason.value !== '' &&
    appealTextLength.value >= 100 &&
    appealTextLength.value <= 2000 &&
    agreedToTerms.value
  )
})

async function handleSubmit() {
  if (!isFormValid.value || appealStore.loading) return

  const result = await appealStore.submitAppeal({
    username: username.value,
    email: email.value,
    banReason: banReason.value,
    appealText: appealText.value
  })

  if (result) {
    submittedAppeal.value = result
  }
}

function copyAppealId() {
  if (submittedAppeal.value) {
    navigator.clipboard.writeText(submittedAppeal.value.id)
    toastStore.success('Appeal ID copied to clipboard!')
  }
}

function resetForm() {
  username.value = ''
  email.value = ''
  banReason.value = ''
  appealText.value = ''
  agreedToTerms.value = false
  submittedAppeal.value = null
  appealStore.clearCurrentAppeal()
}
</script>

<template>
  <!-- Success State -->
  <div v-if="submittedAppeal" class="text-center">
    <div class="w-20 h-20 mx-auto mb-6 rounded-full bg-green-500/10 flex items-center justify-center">
      <CheckCircle class="w-10 h-10 text-green-400" />
    </div>
    <h3 class="text-2xl font-bold text-white mb-2">Appeal Submitted!</h3>
    <p class="text-gray-400 mb-6">
      Your appeal has been submitted successfully. Save your appeal ID to check the status later.
    </p>

    <div class="bg-weenie-dark rounded-lg p-4 mb-6">
      <p class="text-sm text-gray-400 mb-2">Your Appeal ID</p>
      <div class="flex items-center justify-center gap-3">
        <code class="text-xl font-mono text-weenie-gold">{{ submittedAppeal.id }}</code>
        <button
          @click="copyAppealId"
          class="p-2 text-gray-400 hover:text-white hover:bg-white/10 rounded-lg transition-all"
          title="Copy to clipboard"
        >
          <Copy class="w-5 h-5" />
        </button>
      </div>
    </div>

    <div class="bg-yellow-500/10 border border-yellow-500/30 rounded-lg p-4 mb-6 text-left">
      <div class="flex items-start gap-3">
        <AlertCircle class="w-5 h-5 text-yellow-400 flex-shrink-0 mt-0.5" />
        <div>
          <p class="text-yellow-400 font-medium mb-1">Important</p>
          <p class="text-sm text-gray-300">
            Please save your Appeal ID and the email you used. You will need both to check your appeal status.
            Appeals are typically reviewed within 24-48 hours.
          </p>
        </div>
      </div>
    </div>

    <button @click="resetForm" class="btn-secondary">
      Submit Another Appeal
    </button>
  </div>

  <!-- Form State -->
  <form v-else @submit.prevent="handleSubmit" class="space-y-6">
    <!-- Username Field -->
    <div>
      <label for="username" class="block text-sm font-medium text-gray-300 mb-2">
        Minecraft Username
      </label>
      <input
        id="username"
        v-model="username"
        type="text"
        placeholder="Enter your username"
        class="w-full px-4 py-3 bg-weenie-dark border rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 transition-all"
        :class="usernameError ? 'border-red-500 focus:ring-red-500/50' : 'border-gray-700 focus:ring-weenie-gold/50 focus:border-weenie-gold'"
      />
      <p v-if="usernameError" class="mt-1.5 text-sm text-red-400">{{ usernameError }}</p>
      <p v-else class="mt-1.5 text-sm text-gray-500">3-16 characters, letters, numbers, and underscores only</p>
    </div>

    <!-- Email Field -->
    <div>
      <label for="email" class="block text-sm font-medium text-gray-300 mb-2">
        Email Address
      </label>
      <input
        id="email"
        v-model="email"
        type="email"
        placeholder="your@email.com"
        class="w-full px-4 py-3 bg-weenie-dark border rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 transition-all"
        :class="emailError ? 'border-red-500 focus:ring-red-500/50' : 'border-gray-700 focus:ring-weenie-gold/50 focus:border-weenie-gold'"
      />
      <p v-if="emailError" class="mt-1.5 text-sm text-red-400">{{ emailError }}</p>
      <p v-else class="mt-1.5 text-sm text-gray-500">We'll use this to notify you about your appeal</p>
    </div>

    <!-- Ban Reason Dropdown -->
    <div>
      <label for="banReason" class="block text-sm font-medium text-gray-300 mb-2">
        Reason for Ban
      </label>
      <select
        id="banReason"
        v-model="banReason"
        class="w-full px-4 py-3 bg-weenie-dark border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-weenie-gold/50 focus:border-weenie-gold transition-all appearance-none cursor-pointer"
      >
        <option value="" disabled>Select the reason you were banned</option>
        <option v-for="reason in banReasons" :key="reason.value" :value="reason.value">
          {{ reason.label }}
        </option>
      </select>
    </div>

    <!-- Appeal Text -->
    <div>
      <label for="appealText" class="block text-sm font-medium text-gray-300 mb-2">
        Your Appeal
      </label>
      <textarea
        id="appealText"
        v-model="appealText"
        rows="6"
        placeholder="Explain why you believe your ban should be lifted. Be honest and provide any relevant context..."
        class="w-full px-4 py-3 bg-weenie-dark border rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 transition-all resize-none"
        :class="appealTextError ? 'border-red-500 focus:ring-red-500/50' : 'border-gray-700 focus:ring-weenie-gold/50 focus:border-weenie-gold'"
      ></textarea>
      <div class="flex items-center justify-between mt-1.5">
        <p v-if="appealTextError" class="text-sm text-red-400">{{ appealTextError }}</p>
        <p v-else class="text-sm text-gray-500">Minimum 100 characters</p>
        <span
          class="text-sm"
          :class="appealTextLength > 2000 ? 'text-red-400' : appealTextLength >= 100 ? 'text-green-400' : 'text-gray-500'"
        >
          {{ appealTextLength }} / 2000
        </span>
      </div>
    </div>

    <!-- Terms Checkbox -->
    <div class="flex items-start gap-3">
      <input
        id="terms"
        v-model="agreedToTerms"
        type="checkbox"
        class="mt-1 w-4 h-4 rounded border-gray-700 bg-weenie-dark text-weenie-gold focus:ring-weenie-gold/50 focus:ring-offset-0 cursor-pointer"
      />
      <label for="terms" class="text-sm text-gray-400 cursor-pointer">
        I confirm that all information provided is accurate. I understand that providing false information may result in my appeal being denied and additional penalties.
      </label>
    </div>

    <!-- Submit Button -->
    <button
      type="submit"
      :disabled="!isFormValid || appealStore.loading"
      class="btn-primary w-full flex items-center justify-center gap-2"
    >
      <Loader2 v-if="appealStore.loading" class="w-5 h-5 animate-spin" />
      <Send v-else class="w-5 h-5" />
      {{ appealStore.loading ? 'Submitting...' : 'Submit Appeal' }}
    </button>
  </form>
</template>
