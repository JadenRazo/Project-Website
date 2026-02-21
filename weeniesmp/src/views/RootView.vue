<script setup lang="ts">
import Home from './Home.vue'
import Store from './Store.vue'
import { isStoreSubdomain } from '@/utils/subdomain'

const onStoreSubdomain = isStoreSubdomain()
</script>

<template>
  <Store v-if="onStoreSubdomain" />
  <Home v-else />
</template>
