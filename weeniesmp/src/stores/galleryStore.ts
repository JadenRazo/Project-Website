import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface GalleryImage {
  id: string
  url: string
  thumbnail: string
  title: string
  author: string
  description: string
  likes: number
  createdAt: string
  tags: string[]
}

export interface GalleryFilters {
  tags: string[]
  sortBy: 'newest' | 'popular' | 'random'
  search: string
}

export const useGalleryStore = defineStore('gallery', () => {
  const images = ref<GalleryImage[]>([])
  const loading = ref(false)
  const loadingMore = ref(false)
  const error = ref<string | null>(null)
  const filters = ref<GalleryFilters>({
    tags: [],
    sortBy: 'newest',
    search: ''
  })
  const currentPage = ref(1)
  const hasMore = ref(true)
  const totalImages = ref(0)
  const likedImages = ref<Set<string>>(new Set())

  // Available tags for filtering
  const availableTags = ref<string[]>([
    'Builds',
    'Bases',
    'Farms',
    'PvP',
    'Events',
    'Landscape',
    'Redstone',
    'Art',
    'Community'
  ])

  const filteredImages = computed(() => {
    let result = [...images.value]

    // Filter by search
    if (filters.value.search) {
      const searchLower = filters.value.search.toLowerCase()
      result = result.filter(
        img =>
          img.title.toLowerCase().includes(searchLower) ||
          img.author.toLowerCase().includes(searchLower) ||
          img.description.toLowerCase().includes(searchLower)
      )
    }

    // Filter by tags
    if (filters.value.tags.length > 0) {
      result = result.filter(img =>
        filters.value.tags.some(tag => img.tags.includes(tag))
      )
    }

    // Sort
    switch (filters.value.sortBy) {
      case 'newest':
        result.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
        break
      case 'popular':
        result.sort((a, b) => b.likes - a.likes)
        break
      case 'random':
        result.sort(() => Math.random() - 0.5)
        break
    }

    return result
  })

  async function fetchImages(page = 1, append = false) {
    loading.value = true
    error.value = null

    try {
      // Mock API call - replace with actual API endpoint when available
      await new Promise(resolve => setTimeout(resolve, 800))

      const mockImages: GalleryImage[] = generateMockImages(page)

      if (append) {
        images.value = [...images.value, ...mockImages]
      } else {
        images.value = mockImages
      }

      currentPage.value = page
      hasMore.value = page < 5 // Mock: 5 pages of content
      totalImages.value = 60 // Mock total

      // Load liked images from localStorage
      const stored = localStorage.getItem('weeniesmp_liked_images')
      if (stored) {
        likedImages.value = new Set(JSON.parse(stored))
      }
    } catch (e) {
      error.value = 'Failed to load gallery images'
      console.error('Gallery fetch error:', e)
    } finally {
      loading.value = false
    }
  }

  async function loadMore() {
    if (loading.value || loadingMore.value || !hasMore.value) return
    loadingMore.value = true
    await fetchImages(currentPage.value + 1, true)
    loadingMore.value = false
  }

  async function likeImage(id: string) {
    const image = images.value.find(img => img.id === id)
    if (!image) return

    const isLiked = likedImages.value.has(id)

    if (isLiked) {
      image.likes--
      likedImages.value.delete(id)
    } else {
      image.likes++
      likedImages.value.add(id)
    }

    // Persist to localStorage
    localStorage.setItem('weeniesmp_liked_images', JSON.stringify([...likedImages.value]))

    // In a real implementation, this would call an API endpoint
    // await api.post(`/gallery/${id}/like`, { liked: !isLiked })
  }

  function isLiked(id: string): boolean {
    return likedImages.value.has(id)
  }

  function setFilter(key: keyof GalleryFilters, value: unknown) {
    (filters.value as Record<string, unknown>)[key] = value
  }

  function toggleTag(tag: string) {
    const index = filters.value.tags.indexOf(tag)
    if (index === -1) {
      filters.value.tags.push(tag)
    } else {
      filters.value.tags.splice(index, 1)
    }
  }

  function clearFilters() {
    filters.value = {
      tags: [],
      sortBy: 'newest',
      search: ''
    }
  }

  return {
    images,
    loading,
    loadingMore,
    error,
    filters,
    currentPage,
    hasMore,
    totalImages,
    availableTags,
    filteredImages,
    fetchImages,
    loadMore,
    likeImage,
    isLiked,
    setFilter,
    toggleTag,
    clearFilters
  }
})

// Mock data generator - remove when API is implemented
function generateMockImages(page: number): GalleryImage[] {
  const tags = ['Builds', 'Bases', 'Farms', 'PvP', 'Events', 'Landscape', 'Redstone', 'Art', 'Community']
  const authors = ['Steve', 'Alex', 'Notch', 'Jeb', 'Dinnerbone', 'xXPlayer123Xx', 'MinecraftPro', 'BuildMaster']
  const titles = [
    'Epic Castle Build',
    'My Survival Base',
    'Automatic Farm Setup',
    'PvP Arena Design',
    'Community Event Screenshot',
    'Beautiful Sunset View',
    'Redstone Contraption',
    'Pixel Art Creation',
    'Town Center',
    'Mountain Base',
    'Underground Bunker',
    'Ocean Monument Drain'
  ]

  const count = 12
  const startId = (page - 1) * count

  return Array.from({ length: count }, (_, i) => {
    const id = startId + i
    const randomTags = tags
      .sort(() => Math.random() - 0.5)
      .slice(0, Math.floor(Math.random() * 3) + 1)

    // Generate placeholder image URLs with different sizes
    const width = 400 + Math.floor(Math.random() * 200)
    const height = 300 + Math.floor(Math.random() * 200)

    return {
      id: `img-${id}`,
      url: `https://picsum.photos/seed/${id}/${width * 2}/${height * 2}`,
      thumbnail: `https://picsum.photos/seed/${id}/${width}/${height}`,
      title: titles[id % titles.length],
      author: authors[id % authors.length],
      description: `An amazing screenshot from WeenieSMP showing off some incredible work by our community members. This was captured during gameplay and showcases the creativity of our players.`,
      likes: Math.floor(Math.random() * 500) + 10,
      createdAt: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
      tags: randomTags
    }
  })
}
