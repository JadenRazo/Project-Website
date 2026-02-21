<script setup lang="ts">
import { ref, watch } from 'vue'
import { Copy, Check } from 'lucide-vue-next'

const props = defineProps<{
  text: string
}>()

const copied = ref(false)
const announcement = ref('')

// Announce copy success to screen readers
watch(copied, (isCopied) => {
  announcement.value = isCopied ? 'Copied to clipboard' : ''
})

async function copyToClipboard() {
  try {
    await navigator.clipboard.writeText(props.text)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch {
    const textArea = document.createElement('textarea')
    textArea.value = props.text
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  }
}
</script>

<template>
  <button
    @click="copyToClipboard"
    aria-label="Copy to clipboard"
    class="group inline-flex items-center gap-3 px-5 py-3 bg-black/40 border border-white/10 rounded-lg transition-all duration-200 hover:bg-black/60 hover:border-white/20"
    :class="{ 'border-green-500/50 bg-green-500/10': copied }"
  >
    <span class="font-mono text-lg text-white tracking-wide">{{ text }}</span>
    <span
      class="flex items-center justify-center w-8 h-8 rounded-md transition-all"
      :class="copied ? 'bg-green-500/20 text-green-400' : 'bg-white/5 text-gray-400 group-hover:text-white group-hover:bg-white/10'"
    >
      <Check v-if="copied" class="w-4 h-4" />
      <Copy v-else class="w-4 h-4" />
    </span>
    <!-- Screen reader announcement -->
    <span role="status" aria-live="polite" class="sr-only">{{ announcement }}</span>
  </button>
</template>
