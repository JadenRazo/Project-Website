<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { Users, Server, Wifi, WifiOff } from 'lucide-vue-next'
import { useServerStore } from '@/stores/serverStore'

const serverStore = useServerStore()
let interval: number | undefined

onMounted(() => {
  serverStore.fetchStatus()
  interval = setInterval(() => {
    serverStore.fetchStatus()
  }, 60000) as unknown as number
})

onUnmounted(() => {
  if (interval) {
    clearInterval(interval)
  }
})
</script>

<template>
  <section class="py-10 bg-weenie-dark/30">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="card max-w-xl mx-auto">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-xl font-bold text-white flex items-center gap-2">
            <Server class="w-5 h-5 text-weenie-red" />
            Server Status
          </h3>
          <div class="flex items-center gap-2">
            <template v-if="serverStore.loading">
              <div class="w-3 h-3 bg-yellow-500 rounded-full animate-pulse"></div>
              <span class="text-sm text-yellow-400">Checking...</span>
            </template>
            <template v-else-if="serverStore.status.online">
              <Wifi class="w-4 h-4 text-green-400" />
              <span class="text-sm text-green-400">Online</span>
            </template>
            <template v-else>
              <WifiOff class="w-4 h-4 text-red-400" />
              <span class="text-sm text-red-400">Offline</span>
            </template>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div class="bg-weenie-darker/50 rounded-lg p-4 text-center">
            <Users class="w-8 h-8 text-weenie-gold mx-auto mb-2" />
            <div class="text-2xl font-bold text-white">
              {{ serverStore.status.players.online }}
              <span class="text-gray-500 text-sm font-normal">/ {{ serverStore.status.players.max }}</span>
            </div>
            <div class="text-sm text-gray-400">Players Online</div>
          </div>

          <div class="bg-weenie-darker/50 rounded-lg p-4 text-center">
            <Server class="w-8 h-8 text-weenie-red mx-auto mb-2" />
            <div class="text-2xl font-bold text-white">{{ serverStore.status.version }}</div>
            <div class="text-sm text-gray-400">Version</div>
          </div>
        </div>

        <p class="mt-4 text-center text-sm text-gray-500">
          Status updates every minute
        </p>
      </div>
    </div>
  </section>
</template>
