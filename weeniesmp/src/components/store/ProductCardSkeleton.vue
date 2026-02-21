<script setup lang="ts">
import Skeleton from '@/components/common/Skeleton.vue'
</script>

<template>
  <div class="card flex flex-col" aria-hidden="true">
    <!-- Header: Image area and price -->
    <div class="flex items-start justify-between mb-4">
      <!-- Image placeholder -->
      <Skeleton width="9rem" height="5rem" rounded="lg" />
      <!-- Price placeholder -->
      <Skeleton width="4.5rem" height="2rem" rounded="md" />
    </div>

    <!-- Title -->
    <Skeleton width="70%" height="1.5rem" rounded="md" class="mb-2" />

    <!-- Description lines -->
    <div class="flex-1 mb-4 space-y-2">
      <Skeleton width="100%" height="0.875rem" rounded="sm" />
      <Skeleton width="90%" height="0.875rem" rounded="sm" />
      <Skeleton width="75%" height="0.875rem" rounded="sm" />
    </div>

    <!-- Button -->
    <Skeleton width="100%" height="2.75rem" rounded="lg" />
  </div>
</template>
