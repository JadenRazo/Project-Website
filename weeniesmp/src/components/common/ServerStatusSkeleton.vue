<script setup lang="ts">
import Skeleton from '@/components/common/Skeleton.vue'
</script>

<template>
  <section class="py-10 bg-weenie-dark/30" aria-hidden="true">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="card max-w-xl mx-auto">
        <!-- Header -->
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center gap-2">
            <Skeleton width="1.25rem" height="1.25rem" rounded="sm" />
            <Skeleton width="7rem" height="1.5rem" rounded="md" />
          </div>
          <div class="flex items-center gap-2">
            <Skeleton width="0.75rem" height="0.75rem" rounded="full" />
            <Skeleton width="4rem" height="1rem" rounded="md" />
          </div>
        </div>

        <!-- Stats grid -->
        <div class="grid grid-cols-2 gap-4">
          <!-- Players stat -->
          <div class="bg-weenie-darker/50 rounded-lg p-4 text-center">
            <Skeleton width="2rem" height="2rem" rounded="md" class="mx-auto mb-2" />
            <Skeleton width="5rem" height="2rem" rounded="md" class="mx-auto mb-2" />
            <Skeleton width="6rem" height="0.875rem" rounded="sm" class="mx-auto" />
          </div>

          <!-- Version stat -->
          <div class="bg-weenie-darker/50 rounded-lg p-4 text-center">
            <Skeleton width="2rem" height="2rem" rounded="md" class="mx-auto mb-2" />
            <Skeleton width="4rem" height="2rem" rounded="md" class="mx-auto mb-2" />
            <Skeleton width="4rem" height="0.875rem" rounded="sm" class="mx-auto" />
          </div>
        </div>

        <!-- Footer text -->
        <div class="mt-4 flex justify-center">
          <Skeleton width="10rem" height="0.875rem" rounded="sm" />
        </div>
      </div>
    </div>
  </section>
</template>
