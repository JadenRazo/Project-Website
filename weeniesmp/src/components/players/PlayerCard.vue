<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { formatDate } from '@/utils/formatters'
import type { PlayerSearchResult } from '@/stores/playerStore'

interface Props {
  player: PlayerSearchResult
}

defineProps<Props>()

function getCraftarAvatar(uuid: string): string {
  // Remove dashes from UUID for Crafatar
  const cleanUuid = uuid.replace(/-/g, '')
  return `https://crafatar.com/avatars/${cleanUuid}?size=64&overlay`
}
</script>

<template>
  <RouterLink
    :to="`/players/${player.username}`"
    class="card group flex items-center gap-4 hover:bg-white/[0.06]"
  >
    <!-- Avatar -->
    <div class="relative flex-shrink-0">
      <img
        :src="getCraftarAvatar(player.uuid)"
        :alt="player.username"
        class="w-14 h-14 rounded-lg"
        loading="lazy"
      />
      <!-- Online indicator -->
      <div
        v-if="player.isOnline"
        class="absolute -bottom-1 -right-1 w-4 h-4 bg-green-500 rounded-full border-2 border-weenie-darker"
        title="Online"
      ></div>
    </div>

    <!-- Info -->
    <div class="flex-1 min-w-0">
      <h3 class="font-semibold text-white group-hover:text-weenie-gold transition-colors truncate">
        {{ player.username }}
      </h3>
      <p class="text-sm text-gray-400">
        <template v-if="player.isOnline">
          <span class="text-green-400">Online now</span>
        </template>
        <template v-else>
          Last seen {{ formatDate(player.lastSeen) }}
        </template>
      </p>
    </div>

    <!-- Arrow -->
    <svg
      class="w-5 h-5 text-gray-500 group-hover:text-weenie-gold transition-colors flex-shrink-0"
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
    </svg>
  </RouterLink>
</template>
