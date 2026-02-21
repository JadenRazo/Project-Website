<script setup lang="ts">
import { onMounted } from 'vue'
import { Users } from 'lucide-vue-next'
import CopyButton from '@/components/common/CopyButton.vue'
import { useServerStore } from '@/stores/serverStore'

// Optimized hero background images
import heroWebp640 from '@/assets/images/bg-hero.jpg?w=640&format=webp'
import heroWebp1024 from '@/assets/images/bg-hero.jpg?w=1024&format=webp'
import heroWebp1920 from '@/assets/images/bg-hero.jpg?w=1920&format=webp'
import heroJpg640 from '@/assets/images/bg-hero.jpg?w=640&format=jpg'
import heroJpg1024 from '@/assets/images/bg-hero.jpg?w=1024&format=jpg'
import heroJpg1920 from '@/assets/images/bg-hero.jpg?w=1920&format=jpg'

const serverStore = useServerStore()

onMounted(() => {
  serverStore.fetchStatus()
})
</script>

<template>
  <section class="relative pt-24 pb-10 overflow-hidden min-h-[85vh] flex items-center">
    <!-- Background Image -->
    <div class="absolute inset-0">
      <picture>
        <source
          type="image/webp"
          :srcset="`${heroWebp640} 640w, ${heroWebp1024} 1024w, ${heroWebp1920} 1920w`"
          sizes="100vw"
        />
        <img
          :src="heroJpg1920"
          :srcset="`${heroJpg640} 640w, ${heroJpg1024} 1024w, ${heroJpg1920} 1920w`"
          sizes="100vw"
          alt="WeenieSMP Minecraft server landscape"
          class="w-full h-full object-cover object-center"
          loading="eager"
          fetchpriority="high"
        />
      </picture>
    </div>
    <!-- Dark Overlay for readability -->
    <div class="absolute inset-0 bg-gradient-to-b from-black/70 via-black/50 to-weenie-darker"></div>
    <!-- Subtle red accent glow -->
    <div class="absolute inset-0 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-weenie-red/10 via-transparent to-transparent"></div>

    <div class="relative z-10 w-full max-w-2xl mx-auto px-6 text-center">
      <!-- Logo -->
      <img
        src="/logo-v2.png"
        alt="WeenieSMP"
        class="h-16 sm:h-20 md:h-24 w-auto mx-auto mb-2 drop-shadow-2xl"
      />

      <!-- Tagline -->
      <p class="text-lg sm:text-xl text-gray-400 mb-4">
        Your Ultimate Survival Adventure
      </p>

      <!-- Server IP -->
      <div class="flex flex-col items-center gap-2 mb-4">
        <div class="flex items-center gap-2">
          <span class="text-xs text-gray-500 uppercase tracking-wide">Java</span>
          <CopyButton text="play.weeniesmp.net" />
        </div>
        <div class="flex items-center gap-2">
          <span class="text-xs text-gray-500 uppercase tracking-wide">Bedrock</span>
          <CopyButton text="play.weeniesmp.net" />
          <span class="text-xs text-gray-400">Port: 19011</span>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex justify-center gap-3 mb-5">
        <a
          href="https://discord.com/invite/weeniesmp"
          target="_blank"
          rel="noopener noreferrer"
          class="px-6 py-2.5 bg-[#5865F2] hover:bg-[#4752C4] text-white text-sm font-medium rounded-lg transition-colors"
        >
          Join Discord
        </a>
        <a
          href="https://store.weeniesmp.net"
          class="px-6 py-2.5 bg-white/10 hover:bg-white/15 text-white text-sm font-medium rounded-lg transition-colors border border-white/10"
        >
          Visit Store
        </a>
      </div>

      <!-- Status Row -->
      <div class="flex flex-wrap justify-center gap-2 text-sm text-gray-500">
        <div class="flex items-center gap-2 px-3 py-1.5 bg-black/30 rounded-full">
          <span
            class="w-2 h-2 rounded-full"
            :class="serverStore.status.online ? 'bg-green-500' : 'bg-red-500'"
          ></span>
          {{ serverStore.status.online ? 'Online' : 'Offline' }}
        </div>
        <div class="flex items-center gap-2 px-3 py-1.5 bg-black/30 rounded-full">
          <Users class="w-3.5 h-3.5" />
          {{ serverStore.status.players.online }} players
        </div>
        <div class="px-3 py-1.5 bg-black/30 rounded-full">
          Java & Bedrock
        </div>
        <div class="px-3 py-1.5 bg-black/30 rounded-full">
          {{ serverStore.status.version }}
        </div>
      </div>
    </div>
  </section>
</template>
