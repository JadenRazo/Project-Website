<script setup lang="ts">
import { computed } from 'vue'
import { Clock, Search, CheckCircle, XCircle } from 'lucide-vue-next'
import type { AppealStatus } from '@/stores/appealStore'

const props = defineProps<{
  status: AppealStatus
}>()

const statusConfig = computed(() => {
  switch (props.status) {
    case 'pending':
      return {
        label: 'Pending',
        icon: Clock,
        bgColor: 'bg-yellow-500/10',
        borderColor: 'border-yellow-500/30',
        textColor: 'text-yellow-400'
      }
    case 'under_review':
      return {
        label: 'Under Review',
        icon: Search,
        bgColor: 'bg-blue-500/10',
        borderColor: 'border-blue-500/30',
        textColor: 'text-blue-400'
      }
    case 'approved':
      return {
        label: 'Approved',
        icon: CheckCircle,
        bgColor: 'bg-green-500/10',
        borderColor: 'border-green-500/30',
        textColor: 'text-green-400'
      }
    case 'denied':
      return {
        label: 'Denied',
        icon: XCircle,
        bgColor: 'bg-red-500/10',
        borderColor: 'border-red-500/30',
        textColor: 'text-red-400'
      }
    default:
      return {
        label: 'Unknown',
        icon: Clock,
        bgColor: 'bg-gray-500/10',
        borderColor: 'border-gray-500/30',
        textColor: 'text-gray-400'
      }
  }
})
</script>

<template>
  <span
    class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-sm font-medium border"
    :class="[statusConfig.bgColor, statusConfig.borderColor, statusConfig.textColor]"
  >
    <component :is="statusConfig.icon" class="w-4 h-4" />
    {{ statusConfig.label }}
  </span>
</template>
