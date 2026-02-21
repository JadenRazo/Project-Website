<script setup lang="ts">
import { ref } from 'vue'
import { Scale, FileText, Search } from 'lucide-vue-next'
import AppealForm from '@/components/appeals/AppealForm.vue'
import AppealStatusCheck from '@/components/appeals/AppealStatusCheck.vue'

type Tab = 'submit' | 'status'
const activeTab = ref<Tab>('submit')

const tabs = [
  { id: 'submit' as Tab, name: 'Submit Appeal', icon: FileText },
  { id: 'status' as Tab, name: 'Check Status', icon: Search }
]
</script>

<template>
  <div class="min-h-screen pt-24 pb-16 bg-weenie-darker">
    <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="text-center mb-12">
        <Scale class="w-12 h-12 text-weenie-gold mx-auto mb-4" />
        <h1 class="text-4xl md:text-5xl font-bold mb-4">
          <span class="gradient-text">Ban Appeals</span>
        </h1>
        <p class="text-gray-400 max-w-xl mx-auto">
          If you believe you were banned unfairly, you can submit an appeal for review by our moderation team.
        </p>
      </div>

      <!-- Important Notice -->
      <div class="bg-yellow-500/10 border border-yellow-500/30 rounded-xl p-4 mb-8">
        <h3 class="text-yellow-400 font-semibold mb-2">Before You Appeal</h3>
        <ul class="text-sm text-yellow-200/80 space-y-1 list-disc list-inside">
          <li>Appeals are reviewed within 3-5 business days</li>
          <li>Submitting multiple appeals will delay your review</li>
          <li>Be honest and provide accurate information</li>
          <li>False information will result in automatic denial</li>
        </ul>
      </div>

      <!-- Tabs -->
      <div class="flex gap-2 mb-8">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          class="flex-1 flex items-center justify-center gap-2 px-4 py-3 rounded-lg font-medium transition-all"
          :class="activeTab === tab.id
            ? 'bg-weenie-gradient text-white'
            : 'bg-weenie-dark text-gray-400 hover:text-white hover:bg-weenie-dark/80'"
        >
          <component :is="tab.icon" class="w-5 h-5" />
          {{ tab.name }}
        </button>
      </div>

      <!-- Tab Content -->
      <div class="card">
        <AppealForm v-if="activeTab === 'submit'" />
        <AppealStatusCheck v-else />
      </div>

      <!-- FAQ Link -->
      <div class="text-center mt-8">
        <p class="text-gray-500 text-sm">
          Have questions about the appeal process?
          <RouterLink to="/faq" class="text-weenie-gold hover:underline">
            Check our FAQ
          </RouterLink>
        </p>
      </div>
    </div>
  </div>
</template>
