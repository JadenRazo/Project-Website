<script setup lang="ts">
import { ExternalLink } from 'lucide-vue-next'

const voteLinks = [
  { name: 'MC Server List', url: 'https://minecraft-server-list.com/server/weeniesmp' },
  { name: 'MC Servers', url: 'https://minecraftservers.org/server/weeniesmp' }
]
</script>

<template>
  <section class="py-4">
    <div class="max-w-4xl mx-auto px-6">
      <div class="flex flex-col sm:flex-row items-center justify-between gap-4 p-5 rounded-xl bg-white/[0.02] border border-white/5">
        <div>
          <h3 class="text-white font-medium">Vote for Rewards</h3>
          <p class="text-sm text-gray-500">Each vote earns you a reward key</p>
        </div>
        <div class="flex gap-2">
          <a
            v-for="link in voteLinks"
            :key="link.name"
            :href="link.url"
            target="_blank"
            rel="noopener noreferrer"
            class="inline-flex items-center gap-1.5 px-4 py-2 text-sm text-gray-400 hover:text-white bg-white/5 hover:bg-white/10 rounded-lg transition-colors"
          >
            {{ link.name }}
            <ExternalLink class="w-3 h-3" />
          </a>
        </div>
      </div>
    </div>
  </section>
</template>
