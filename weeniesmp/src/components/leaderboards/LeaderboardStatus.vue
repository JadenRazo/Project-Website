<script setup lang="ts">
import { RefreshCw, Clock, CheckCircle } from 'lucide-vue-next'
import { computed } from 'vue'

const props = defineProps<{
  updating: boolean
  lastUpdated: number
  cacheAge: number | null
}>()

const lastUpdatedText = computed(() => {
  if (!props.lastUpdated) return 'Never'

  const seconds = Math.floor((Date.now() - props.lastUpdated) / 1000)

  if (seconds < 60) return 'Just now'
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`
  return `${Math.floor(seconds / 86400)}d ago`
})

const statusColor = computed(() => {
  if (props.updating) return 'text-blue-400'
  if (props.cacheAge && props.cacheAge < 60) return 'text-green-400'
  return 'text-gray-400'
})

const statusIcon = computed(() => {
  if (props.updating) return RefreshCw
  if (props.cacheAge && props.cacheAge < 60) return CheckCircle
  return Clock
})
</script>

<template>
  <div class="flex items-center justify-center gap-2 text-sm mt-4 mb-2">
    <component
      :is="statusIcon"
      :class="[
        'w-4 h-4',
        statusColor,
        { 'animate-spin': updating }
      ]"
    />
    <span :class="statusColor">
      <template v-if="updating">
        Updating...
      </template>
      <template v-else>
        Updated {{ lastUpdatedText }}
      </template>
    </span>
  </div>
</template>
