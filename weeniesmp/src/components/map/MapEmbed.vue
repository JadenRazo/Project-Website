<script setup lang="ts">
import { ref, computed } from 'vue'
import { Loader2, AlertTriangle, RefreshCw } from 'lucide-vue-next'

interface Props {
  mapUrl: string
  height?: string
}

const props = withDefaults(defineProps<Props>(), {
  height: '600px'
})

const loading = ref(true)
const hasError = ref(false)
const iframeRef = ref<HTMLIFrameElement | null>(null)

const containerStyle = computed(() => ({
  height: props.height,
  minHeight: '400px'
}))

function onIframeLoad() {
  loading.value = false
  hasError.value = false
}

function onIframeError() {
  loading.value = false
  hasError.value = true
}

function retry() {
  loading.value = true
  hasError.value = false
  if (iframeRef.value) {
    iframeRef.value.src = props.mapUrl
  }
}
</script>

<template>
  <div class="relative w-full rounded-xl overflow-hidden bg-weenie-dark border border-white/10" :style="containerStyle">
    <!-- Loading Overlay -->
    <Transition
      enter-active-class="transition-opacity duration-300"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-200"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="loading && !hasError"
        class="absolute inset-0 flex flex-col items-center justify-center bg-weenie-dark z-10"
      >
        <Loader2 class="w-10 h-10 text-weenie-gold animate-spin mb-4" />
        <p class="text-gray-400 text-sm">Loading map...</p>
      </div>
    </Transition>

    <!-- Error State -->
    <Transition
      enter-active-class="transition-opacity duration-300"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-200"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="hasError"
        class="absolute inset-0 flex flex-col items-center justify-center bg-weenie-dark z-10"
      >
        <AlertTriangle class="w-12 h-12 text-weenie-red mb-4" />
        <h3 class="text-lg font-semibold text-white mb-2">Map Unavailable</h3>
        <p class="text-gray-400 text-sm mb-6 text-center max-w-md px-4">
          The map server may be temporarily down or undergoing maintenance.
        </p>
        <button
          @click="retry"
          class="inline-flex items-center gap-2 px-4 py-2 bg-weenie-gold hover:bg-weenie-gold/90 text-black font-medium rounded-lg transition-colors"
        >
          <RefreshCw class="w-4 h-4" />
          Try Again
        </button>
      </div>
    </Transition>

    <!-- Map Iframe -->
    <iframe
      ref="iframeRef"
      :src="mapUrl"
      class="w-full h-full border-0"
      :class="{ 'invisible': loading || hasError }"
      allow="fullscreen"
      loading="lazy"
      @load="onIframeLoad"
      @error="onIframeError"
      title="WeenieSMP Server Map"
    />
  </div>
</template>
