<script setup lang="ts">
import Navbar from '@/components/common/Navbar.vue'
import Footer from '@/components/common/Footer.vue'
import CartDrawer from '@/components/store/CartDrawer.vue'
import ToastContainer from '@/components/common/ToastContainer.vue'
import { isStoreSubdomain } from '@/utils/subdomain'

// Check once at setup - hostname doesn't change during session
const isStoreDomain = isStoreSubdomain()
</script>

<template>
  <div class="min-h-screen flex flex-col">
    <Navbar :is-store="isStoreDomain" />
    <main class="flex-1">
      <router-view />
    </main>
    <Footer />

    <!-- Global Cart Drawer (only on store subdomain) -->
    <CartDrawer v-if="isStoreDomain" />

    <!-- Toast Notifications -->
    <ToastContainer />
  </div>
</template>
