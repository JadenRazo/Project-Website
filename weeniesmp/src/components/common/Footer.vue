<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { isStoreSubdomain, mainSiteUrl, storeUrl } from '@/utils/subdomain'

const onStoreSubdomain = isStoreSubdomain()
</script>

<template>
  <footer class="py-8 border-t border-white/5 bg-weenie-darker">
    <div class="max-w-6xl mx-auto px-6">
      <!-- Link Sections -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-8 mb-8">
        <!-- Community -->
        <div>
          <h4 class="text-white font-semibold mb-3">Community</h4>
          <ul class="space-y-2 text-sm text-gray-500">
            <li>
              <a v-if="onStoreSubdomain" :href="`${mainSiteUrl}/players`" class="hover:text-white transition-colors">Players</a>
              <RouterLink v-else to="/players" class="hover:text-white transition-colors">Players</RouterLink>
            </li>
            <li>
              <a v-if="onStoreSubdomain" :href="`${mainSiteUrl}/leaderboards`" class="hover:text-white transition-colors">Leaderboards</a>
              <RouterLink v-else to="/leaderboards" class="hover:text-white transition-colors">Leaderboards</RouterLink>
            </li>
            <li>
              <a v-if="onStoreSubdomain" :href="`${mainSiteUrl}/gallery`" class="hover:text-white transition-colors">Gallery</a>
              <RouterLink v-else to="/gallery" class="hover:text-white transition-colors">Gallery</RouterLink>
            </li>
            <li>
              <a href="https://discord.com/invite/weeniesmp" target="_blank" rel="noopener noreferrer" class="hover:text-white transition-colors">Discord</a>
            </li>
          </ul>
        </div>

        <!-- Resources -->
        <div>
          <h4 class="text-white font-semibold mb-3">Resources</h4>
          <ul class="space-y-2 text-sm text-gray-500">
            <li>
              <a v-if="onStoreSubdomain" :href="`${mainSiteUrl}/wiki`" class="hover:text-white transition-colors">Wiki</a>
              <RouterLink v-else to="/wiki" class="hover:text-white transition-colors">Wiki</RouterLink>
            </li>
            <li>
              <a v-if="onStoreSubdomain" :href="`${mainSiteUrl}/rules`" class="hover:text-white transition-colors">Rules</a>
              <RouterLink v-else to="/rules" class="hover:text-white transition-colors">Rules</RouterLink>
            </li>
            <li>
              <a v-if="onStoreSubdomain" :href="`${mainSiteUrl}/faq`" class="hover:text-white transition-colors">FAQ</a>
              <RouterLink v-else to="/faq" class="hover:text-white transition-colors">FAQ</RouterLink>
            </li>
            <li>
              <a v-if="onStoreSubdomain" :href="`${mainSiteUrl}/map`" class="hover:text-white transition-colors">Server Map</a>
              <RouterLink v-else to="/map" class="hover:text-white transition-colors">Server Map</RouterLink>
            </li>
          </ul>
        </div>

        <!-- Support -->
        <div>
          <h4 class="text-white font-semibold mb-3">Support</h4>
          <ul class="space-y-2 text-sm text-gray-500">
            <li>
              <a v-if="onStoreSubdomain" :href="`${mainSiteUrl}/appeals`" class="hover:text-white transition-colors">Ban Appeals</a>
              <RouterLink v-else to="/appeals" class="hover:text-white transition-colors">Ban Appeals</RouterLink>
            </li>
            <li>
              <a href="https://discord.com/invite/weeniesmp" target="_blank" rel="noopener noreferrer" class="hover:text-white transition-colors">Get Help</a>
            </li>
            <li>
              <a v-if="onStoreSubdomain" :href="`${mainSiteUrl}/news`" class="hover:text-white transition-colors">News</a>
              <RouterLink v-else to="/news" class="hover:text-white transition-colors">News</RouterLink>
            </li>
          </ul>
        </div>

        <!-- Store -->
        <div>
          <h4 class="text-white font-semibold mb-3">Store</h4>
          <ul class="space-y-2 text-sm text-gray-500">
            <li>
              <a :href="storeUrl" class="hover:text-white transition-colors">Browse Store</a>
            </li>
            <li>
              <a v-if="onStoreSubdomain" :href="`${mainSiteUrl}/dashboard`" class="hover:text-white transition-colors">My Account</a>
              <RouterLink v-else to="/dashboard" class="hover:text-white transition-colors">My Account</RouterLink>
            </li>
          </ul>
        </div>
      </div>

      <!-- Bottom Bar -->
      <div class="pt-6 border-t border-white/5 flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-gray-500">
        <p>&copy; {{ new Date().getFullYear() }} WeenieSMP. All rights reserved.</p>
        <p class="text-gray-600">play.weeniesmp.net</p>
      </div>
    </div>
  </footer>
</template>
