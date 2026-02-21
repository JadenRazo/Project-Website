<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Map, Maximize2, Minimize2, ExternalLink, ZoomIn, ZoomOut, Crosshair } from 'lucide-vue-next'
import MapEmbed from '@/components/map/MapEmbed.vue'

const MAP_URL = 'https://map.weeniesmp.net'

const isFullscreen = ref(false)
const mapContainer = ref<HTMLElement | null>(null)
const viewportWidth = ref(window.innerWidth)

const mapHeight = computed(() => {
  if (isFullscreen.value) return '100vh'
  return viewportWidth.value < 768 ? 'calc(100vh - 100px)' : 'calc(100vh - 180px)'
})

function toggleFullscreen() {
  if (!mapContainer.value) return

  if (!isFullscreen.value) {
    if (mapContainer.value.requestFullscreen) {
      mapContainer.value.requestFullscreen()
    }
  } else {
    if (document.exitFullscreen) {
      document.exitFullscreen()
    }
  }
}

function handleFullscreenChange() {
  isFullscreen.value = !!document.fullscreenElement
}

function openInNewTab() {
  window.open(MAP_URL, '_blank', 'noopener,noreferrer')
}

function handleResize() {
  viewportWidth.value = window.innerWidth
}

onMounted(() => {
  document.addEventListener('fullscreenchange', handleFullscreenChange)
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  document.removeEventListener('fullscreenchange', handleFullscreenChange)
  window.removeEventListener('resize', handleResize)
})
</script>

<template>
  <div class="min-h-screen pt-20 pb-8 bg-weenie-darker">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
        <div class="flex items-center gap-3">
          <Map class="w-8 h-8 text-weenie-gold" />
          <div>
            <h1 class="text-2xl md:text-3xl font-bold text-white">Server Map</h1>
            <p class="text-gray-400 text-sm">Explore the WeenieSMP world in real-time</p>
          </div>
        </div>

        <!-- Controls -->
        <div class="flex items-center gap-2">
          <button
            @click="toggleFullscreen"
            class="inline-flex items-center gap-2 px-4 py-2 bg-white/5 hover:bg-white/10 text-gray-300 hover:text-white rounded-lg transition-colors border border-white/10"
            :title="isFullscreen ? 'Exit fullscreen' : 'Enter fullscreen'"
          >
            <Maximize2 v-if="!isFullscreen" class="w-4 h-4" />
            <Minimize2 v-else class="w-4 h-4" />
            <span class="hidden sm:inline">{{ isFullscreen ? 'Exit' : 'Fullscreen' }}</span>
          </button>
          <button
            @click="openInNewTab"
            class="inline-flex items-center gap-2 px-4 py-2 bg-weenie-gold hover:bg-weenie-gold/90 text-black font-medium rounded-lg transition-colors"
          >
            <ExternalLink class="w-4 h-4" />
            <span class="hidden sm:inline">Open in New Tab</span>
          </button>
        </div>
      </div>

      <!-- Map Container -->
      <div ref="mapContainer" class="relative">
        <MapEmbed :map-url="MAP_URL" :height="mapHeight" />

        <!-- Map Controls Overlay -->
        <div class="absolute bottom-2 left-2 sm:bottom-4 sm:left-4 flex flex-col gap-2 z-20">
          <div class="bg-black/70 backdrop-blur-sm rounded-lg border border-white/10 p-0.5 sm:p-1">
            <button
              class="p-1.5 sm:p-2 text-gray-400 hover:text-white hover:bg-white/10 rounded transition-colors block"
              title="Zoom In (use scroll wheel in map)"
            >
              <ZoomIn class="w-4 h-4 sm:w-5 sm:h-5" />
            </button>
            <button
              class="p-1.5 sm:p-2 text-gray-400 hover:text-white hover:bg-white/10 rounded transition-colors block"
              title="Zoom Out (use scroll wheel in map)"
            >
              <ZoomOut class="w-4 h-4 sm:w-5 sm:h-5" />
            </button>
            <button
              class="p-1.5 sm:p-2 text-gray-400 hover:text-white hover:bg-white/10 rounded transition-colors block"
              title="Coordinates (click on map)"
            >
              <Crosshair class="w-4 h-4 sm:w-5 sm:h-5" />
            </button>
          </div>
        </div>

        <!-- Info Badge -->
        <div class="absolute top-2 right-2 sm:top-4 sm:right-4 z-20">
          <div class="bg-black/70 backdrop-blur-sm rounded-lg border border-white/10 px-2 py-1 sm:px-3 sm:py-2">
            <p class="text-[10px] sm:text-xs text-gray-400">
              <span class="text-weenie-gold font-medium">Tip:</span>
              <span class="hidden sm:inline">Use mouse to pan, scroll to zoom</span>
              <span class="sm:hidden">Drag & pinch</span>
            </p>
          </div>
        </div>
      </div>

      <!-- Map Legend / Info -->
      <div class="mt-6 grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div class="card p-4">
          <h3 class="text-sm font-medium text-white mb-2">Real-time Updates</h3>
          <p class="text-xs text-gray-400">
            Player positions and world changes update automatically.
          </p>
        </div>
        <div class="card p-4">
          <h3 class="text-sm font-medium text-white mb-2">Multiple Layers</h3>
          <p class="text-xs text-gray-400">
            Switch between surface, cave, and nether views.
          </p>
        </div>
        <div class="card p-4">
          <h3 class="text-sm font-medium text-white mb-2">Player Markers</h3>
          <p class="text-xs text-gray-400">
            See online players' locations on the map.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
