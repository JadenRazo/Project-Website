<script setup lang="ts">
import { ref } from 'vue'
import { User, Loader2, ExternalLink } from 'lucide-vue-next'
import { useUserStore } from '@/stores/userStore'

const userStore = useUserStore()
const usernameInput = ref('')
const submitting = ref(false)

async function handleSubmit() {
  if (!usernameInput.value.trim()) return

  submitting.value = true
  const success = await userStore.login(usernameInput.value.trim())
  submitting.value = false

  if (success) {
    usernameInput.value = ''
  }
}

function isValidFormat(name: string): boolean {
  return /^[a-zA-Z0-9_]{3,16}$/.test(name)
}
</script>

<template>
  <div class="min-h-[60vh] flex items-center justify-center px-4">
    <div class="w-full max-w-md">
      <div class="card p-8 text-center">
        <div class="w-16 h-16 bg-weenie-gradient rounded-full flex items-center justify-center mx-auto mb-6">
          <User class="w-8 h-8 text-white" />
        </div>

        <h1 class="text-2xl font-bold text-white mb-2">Sign In to Dashboard</h1>
        <p class="text-gray-400 mb-8">
          Enter your Minecraft username to view your purchase history and account details.
        </p>

        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div class="relative">
            <input
              v-model="usernameInput"
              type="text"
              placeholder="Your Minecraft username"
              class="w-full px-4 py-3 bg-weenie-darker border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-weenie-gold transition-colors"
              :class="{
                'border-red-500': usernameInput && !isValidFormat(usernameInput),
                'border-green-500': usernameInput && isValidFormat(usernameInput)
              }"
              :disabled="submitting"
            />
          </div>

          <p
            v-if="usernameInput && !isValidFormat(usernameInput)"
            class="text-red-400 text-sm text-left"
          >
            Username must be 3-16 characters (letters, numbers, underscores only)
          </p>

          <p v-if="userStore.error" class="text-red-400 text-sm">
            {{ userStore.error }}
          </p>

          <button
            type="submit"
            :disabled="!usernameInput.trim() || !isValidFormat(usernameInput) || submitting"
            class="w-full py-3 bg-weenie-gradient text-white font-semibold rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <Loader2 v-if="submitting" class="w-5 h-5 animate-spin" />
            <span>{{ submitting ? 'Signing in...' : 'Sign In' }}</span>
          </button>
        </form>

        <div class="mt-8 pt-6 border-t border-gray-700">
          <p class="text-gray-500 text-sm mb-3">
            New to Minecraft?
          </p>
          <a
            href="https://www.minecraft.net/get-minecraft"
            target="_blank"
            rel="noopener noreferrer"
            class="inline-flex items-center gap-2 text-weenie-gold hover:underline text-sm"
          >
            Get Minecraft
            <ExternalLink class="w-4 h-4" />
          </a>
        </div>
      </div>
    </div>
  </div>
</template>
