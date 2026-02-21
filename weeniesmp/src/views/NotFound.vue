<script setup lang="ts">
import { Home, ArrowLeft } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<template>
  <div class="min-h-screen flex items-center justify-center pt-16 bg-weenie-darker">
    <div class="text-center px-4">
      <h1 class="text-8xl md:text-9xl font-bold gradient-text mb-4">404</h1>
      <h2 class="text-2xl md:text-3xl font-semibold text-white mb-4">Page Not Found</h2>
      <p class="text-gray-400 mb-8 max-w-md mx-auto">
        Looks like you wandered into the void. This page doesn't exist or has been moved.
      </p>

      <div class="flex flex-col sm:flex-row items-center justify-center gap-4">
        <RouterLink to="/" class="btn-primary flex items-center gap-2">
          <Home class="w-5 h-5" />
          Go Home
        </RouterLink>
        <button
          @click="router.back()"
          class="btn-secondary flex items-center gap-2"
        >
          <ArrowLeft class="w-5 h-5" />
          Go Back
        </button>
      </div>
    </div>
  </div>
</template>
