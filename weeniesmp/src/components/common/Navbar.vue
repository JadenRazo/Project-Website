<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { Menu, X, ShoppingCart } from 'lucide-vue-next'
import { isStoreSubdomain, mainSiteUrl, storeUrl } from '@/utils/subdomain'
import { useTebexStore } from '@/stores/tebexStore'
import { useUiStore } from '@/stores/uiStore'

const route = useRoute()
const mobileMenuOpen = ref(false)
const tebexStore = useTebexStore()
const uiStore = useUiStore()

// Check once at component setup time (client-side)
const onStoreSubdomain = isStoreSubdomain()

interface NavLink {
  name: string
  routerPath?: string
  href?: string
  external?: boolean
}

// Static links based on which subdomain we're on
const navLinks: NavLink[] = onStoreSubdomain
  ? [
      { name: 'Home', href: mainSiteUrl },
      { name: 'Store', routerPath: '/' },
      { name: 'News', href: `${mainSiteUrl}/news` },
      { name: 'Leaderboards', href: `${mainSiteUrl}/leaderboards` },
      { name: 'Wiki', href: `${mainSiteUrl}/wiki` },
      { name: 'Discord', href: 'https://discord.com/invite/weeniesmp', external: true },
    ]
  : [
      { name: 'Home', routerPath: '/' },
      { name: 'Store', href: storeUrl },
      { name: 'News', routerPath: '/news' },
      { name: 'Leaderboards', routerPath: '/leaderboards' },
      { name: 'Wiki', routerPath: '/wiki' },
      { name: 'Discord', href: 'https://discord.com/invite/weeniesmp', external: true },
    ]

function isActive(link: NavLink): boolean {
  if (link.routerPath) {
    return route.path === link.routerPath
  }
  return false
}
</script>

<template>
  <nav aria-label="Main navigation" class="fixed top-0 left-0 right-0 z-50 bg-black/80 backdrop-blur-lg border-b border-white/5">
    <div class="max-w-6xl mx-auto px-4 sm:px-6">
      <div class="flex items-center justify-between h-16">
        <!-- Logo -->
        <a v-if="onStoreSubdomain" :href="mainSiteUrl" class="flex-shrink-0">
          <img
            src="/logo-v2.png"
            alt="WeenieSMP"
            class="h-12 w-auto hover:opacity-90 transition-opacity"
          />
        </a>
        <RouterLink v-else to="/" class="flex-shrink-0">
          <img
            src="/logo-v2.png"
            alt="WeenieSMP"
            class="h-12 w-auto hover:opacity-90 transition-opacity"
          />
        </RouterLink>

        <!-- Desktop Nav -->
        <div class="hidden md:flex items-center gap-1">
          <template v-for="link in navLinks" :key="link.name">
            <!-- External or cross-domain link -->
            <a
              v-if="link.href"
              :href="link.href"
              :target="link.external ? '_blank' : '_self'"
              :rel="link.external ? 'noopener noreferrer' : undefined"
              class="px-4 py-2 text-sm font-medium rounded-lg transition-all text-gray-400 hover:text-white hover:bg-white/5"
            >
              {{ link.name }}
            </a>
            <!-- Internal router link -->
            <RouterLink
              v-else-if="link.routerPath"
              :to="link.routerPath"
              class="px-4 py-2 text-sm font-medium rounded-lg transition-all"
              :class="isActive(link)
                ? 'text-white bg-white/10'
                : 'text-gray-400 hover:text-white hover:bg-white/5'"
            >
              {{ link.name }}
            </RouterLink>
          </template>
        </div>

        <!-- Cart & Mobile Menu Buttons -->
        <div class="flex items-center gap-2">
          <!-- Cart Button (only on store subdomain) -->
          <button
            v-if="onStoreSubdomain"
            @click="uiStore.toggleCart"
            class="relative p-2.5 text-gray-400 hover:text-white hover:bg-white/5 rounded-lg transition-all"
          >
            <ShoppingCart class="w-5 h-5" />
            <span
              v-if="tebexStore.basketItemCount > 0"
              class="absolute -top-0.5 -right-0.5 w-5 h-5 bg-weenie-red text-white text-xs font-bold rounded-full flex items-center justify-center"
            >
              {{ tebexStore.basketItemCount }}
            </span>
          </button>

          <button
            @click="mobileMenuOpen = !mobileMenuOpen"
            :aria-expanded="mobileMenuOpen"
            aria-controls="mobile-menu"
            :aria-label="mobileMenuOpen ? 'Close menu' : 'Open menu'"
            class="md:hidden p-2.5 text-gray-400 hover:text-white hover:bg-white/5 rounded-lg transition-all"
          >
            <Menu v-if="!mobileMenuOpen" class="w-5 h-5" />
            <X v-else class="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>

    <!-- Mobile Menu -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div v-if="mobileMenuOpen" id="mobile-menu" class="md:hidden bg-black/95 border-b border-white/5">
        <div class="px-4 py-3 space-y-1">
          <template v-for="link in navLinks" :key="link.name">
            <a
              v-if="link.href"
              :href="link.href"
              :target="link.external ? '_blank' : '_self'"
              :rel="link.external ? 'noopener noreferrer' : undefined"
              class="block px-4 py-2.5 rounded-lg transition-all text-gray-400 hover:text-white hover:bg-white/5"
              @click="mobileMenuOpen = false"
            >
              {{ link.name }}
            </a>
            <RouterLink
              v-else-if="link.routerPath"
              :to="link.routerPath"
              @click="mobileMenuOpen = false"
              class="block px-4 py-2.5 rounded-lg transition-all"
              :class="isActive(link)
                ? 'text-white bg-white/10'
                : 'text-gray-400 hover:text-white hover:bg-white/5'"
            >
              {{ link.name }}
            </RouterLink>
          </template>
        </div>
      </div>
    </Transition>
  </nav>
</template>
