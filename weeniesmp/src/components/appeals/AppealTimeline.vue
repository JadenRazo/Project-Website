<script setup lang="ts">
import { computed } from 'vue'
import { Clock, Search, CheckCircle, XCircle, FileText } from 'lucide-vue-next'
import type { Appeal, AppealStatus } from '@/stores/appealStore'

const props = defineProps<{
  appeal: Appeal
}>()

interface TimelineStep {
  status: AppealStatus | 'submitted'
  label: string
  description: string
  icon: typeof Clock
  completed: boolean
  current: boolean
  date?: string
}

const statusOrder: (AppealStatus | 'submitted')[] = ['submitted', 'pending', 'under_review', 'approved']

function getStatusIndex(status: AppealStatus | 'submitted'): number {
  if (status === 'denied') return 2 // Same level as under_review
  return statusOrder.indexOf(status)
}

const timelineSteps = computed((): TimelineStep[] => {
  const appeal = props.appeal
  const currentIndex = getStatusIndex(appeal.status)
  const isDenied = appeal.status === 'denied'

  const steps: TimelineStep[] = [
    {
      status: 'submitted',
      label: 'Submitted',
      description: 'Your appeal has been submitted',
      icon: FileText,
      completed: true,
      current: false,
      date: formatDate(appeal.createdAt)
    },
    {
      status: 'pending',
      label: 'Pending Review',
      description: 'Waiting for staff to review',
      icon: Clock,
      completed: currentIndex > 1,
      current: currentIndex === 1,
      date: currentIndex >= 1 ? formatDate(appeal.createdAt) : undefined
    },
    {
      status: 'under_review',
      label: 'Under Review',
      description: 'Staff is reviewing your appeal',
      icon: Search,
      completed: currentIndex > 2 || isDenied,
      current: currentIndex === 2 && !isDenied,
      date: currentIndex >= 2 || isDenied ? formatDate(appeal.updatedAt) : undefined
    }
  ]

  if (isDenied) {
    steps.push({
      status: 'denied',
      label: 'Denied',
      description: 'Your appeal was not approved',
      icon: XCircle,
      completed: true,
      current: true,
      date: formatDate(appeal.updatedAt)
    })
  } else {
    steps.push({
      status: 'approved',
      label: 'Approved',
      description: 'Your ban has been lifted',
      icon: CheckCircle,
      completed: appeal.status === 'approved',
      current: appeal.status === 'approved',
      date: appeal.status === 'approved' ? formatDate(appeal.updatedAt) : undefined
    })
  }

  return steps
})

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function getStepStyles(step: TimelineStep) {
  if (step.current && step.status === 'denied') {
    return {
      circle: 'bg-red-500 border-red-500',
      icon: 'text-white',
      line: 'bg-red-500'
    }
  }
  if (step.current && step.status === 'approved') {
    return {
      circle: 'bg-green-500 border-green-500',
      icon: 'text-white',
      line: 'bg-green-500'
    }
  }
  if (step.completed) {
    return {
      circle: 'bg-weenie-gold border-weenie-gold',
      icon: 'text-black',
      line: 'bg-weenie-gold'
    }
  }
  if (step.current) {
    return {
      circle: 'bg-blue-500 border-blue-500 animate-pulse',
      icon: 'text-white',
      line: 'bg-gray-700'
    }
  }
  return {
    circle: 'bg-weenie-dark border-gray-600',
    icon: 'text-gray-500',
    line: 'bg-gray-700'
  }
}
</script>

<template>
  <div class="relative">
    <div class="space-y-8">
      <div
        v-for="(step, index) in timelineSteps"
        :key="step.status"
        class="relative flex items-start gap-4"
      >
        <!-- Vertical line -->
        <div
          v-if="index < timelineSteps.length - 1"
          class="absolute left-5 top-10 w-0.5 h-full -translate-x-1/2"
          :class="getStepStyles(step).line"
        ></div>

        <!-- Icon circle -->
        <div
          class="relative z-10 flex-shrink-0 w-10 h-10 rounded-full border-2 flex items-center justify-center"
          :class="getStepStyles(step).circle"
        >
          <component :is="step.icon" class="w-5 h-5" :class="getStepStyles(step).icon" />
        </div>

        <!-- Content -->
        <div class="flex-1 min-w-0 pb-8">
          <div class="flex items-center justify-between">
            <h4
              class="font-semibold"
              :class="step.completed || step.current ? 'text-white' : 'text-gray-500'"
            >
              {{ step.label }}
            </h4>
            <span v-if="step.date" class="text-sm text-gray-500">
              {{ step.date }}
            </span>
          </div>
          <p class="mt-1 text-sm" :class="step.completed || step.current ? 'text-gray-400' : 'text-gray-600'">
            {{ step.description }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
