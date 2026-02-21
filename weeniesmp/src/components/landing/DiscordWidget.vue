<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ExternalLink, Users } from 'lucide-vue-next'

interface DiscordWidget {
  id: string
  name: string
  instant_invite: string
  presence_count: number
  members: Array<{ id: string; username: string; avatar_url: string }>
}

const discordData = ref<DiscordWidget | null>(null)
const loading = ref(true)

async function fetchDiscordWidget() {
  try {
    const res = await fetch('https://discord.com/api/guilds/1223815321912082492/widget.json')
    if (res.ok) {
      discordData.value = await res.json()
    }
  } catch (e) {
    console.error('Failed to fetch Discord widget:', e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchDiscordWidget)
</script>

<template>
  <section class="py-6">
    <div class="max-w-4xl mx-auto px-6">
      <div class="flex flex-col sm:flex-row items-center justify-between gap-4 p-5 rounded-xl bg-[#5865F2]/10 border border-[#5865F2]/20">
        <div>
          <h3 class="text-white font-medium">Join Our Discord</h3>
          <p class="text-sm text-gray-400">Get support, updates & connect with the community</p>
          <div v-if="discordData" class="flex items-center gap-2 mt-2">
            <span class="relative flex h-2 w-2">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-green-500"></span>
            </span>
            <span class="text-sm text-green-400 font-medium">
              {{ discordData.presence_count }} online
            </span>
          </div>
        </div>
        <a
          href="https://discord.com/invite/weeniesmp"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex items-center gap-2 px-5 py-2.5 text-sm font-medium text-white bg-[#5865F2] hover:bg-[#4752C4] rounded-lg transition-colors"
        >
          <Users class="w-4 h-4" />
          Join Server
          <ExternalLink class="w-3.5 h-3.5" />
        </a>
      </div>
    </div>
  </section>
</template>
