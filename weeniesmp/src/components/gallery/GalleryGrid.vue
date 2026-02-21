<script setup lang="ts">
import type { GalleryImage } from '@/stores/galleryStore'
import ImageCard from './ImageCard.vue'

interface Props {
  images: GalleryImage[]
}

defineProps<Props>()
const emit = defineEmits<{
  imageClick: [image: GalleryImage]
}>()

function handleImageClick(image: GalleryImage) {
  emit('imageClick', image)
}
</script>

<template>
  <div class="gallery-grid">
    <ImageCard
      v-for="image in images"
      :key="image.id"
      :image="image"
      @click="handleImageClick(image)"
    />
  </div>
</template>

<style scoped>
.gallery-grid {
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  gap: 1rem;
}

@media (min-width: 640px) {
  .gallery-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1024px) {
  .gallery-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (min-width: 1280px) {
  .gallery-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}
</style>
