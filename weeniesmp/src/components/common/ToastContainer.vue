<script setup lang="ts">
import { X, CheckCircle, AlertCircle, Info } from 'lucide-vue-next'
import { useToastStore } from '@/stores/toastStore'

const toastStore = useToastStore()

const icons = {
  success: CheckCircle,
  error: AlertCircle,
  info: Info
}

const colors = {
  success: 'bg-green-500/90 border-green-400',
  error: 'bg-red-500/90 border-red-400',
  info: 'bg-blue-500/90 border-blue-400'
}
</script>

<template>
  <Teleport to="body">
    <div
      role="status"
      aria-live="polite"
      aria-atomic="true"
      class="fixed top-20 right-4 z-[100] flex flex-col gap-2 max-w-sm"
    >
      <TransitionGroup
        enter-active-class="transition duration-300 ease-out"
        enter-from-class="opacity-0 translate-x-full"
        enter-to-class="opacity-100 translate-x-0"
        leave-active-class="transition duration-200 ease-in"
        leave-from-class="opacity-100 translate-x-0"
        leave-to-class="opacity-0 translate-x-full"
      >
        <div
          v-for="toast in toastStore.toasts"
          :key="toast.id"
          :class="[
            'flex items-center gap-3 px-4 py-3 rounded-lg border backdrop-blur-sm shadow-lg',
            colors[toast.type]
          ]"
        >
          <component :is="icons[toast.type]" class="w-5 h-5 text-white flex-shrink-0" />
          <span class="text-white text-sm font-medium flex-1">{{ toast.message }}</span>
          <button
            @click="toastStore.removeToast(toast.id)"
            class="p-1 text-white/70 hover:text-white transition-colors flex-shrink-0"
          >
            <X class="w-4 h-4" />
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>
