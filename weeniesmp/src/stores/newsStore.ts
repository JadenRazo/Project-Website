import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface NewsArticle {
  id: string
  slug: string
  title: string
  excerpt: string
  content: string
  category: 'update' | 'event' | 'patch' | 'announcement'
  author: string
  publishedAt: string
  updatedAt?: string
  featuredImage?: string
  tags: string[]
}

export interface NewsFilters {
  category?: string
  search?: string
  page?: number
  limit?: number
}

// News articles will be fetched from the API when available
const mockArticles: NewsArticle[] = []

export const useNewsStore = defineStore('news', () => {
  const articles = ref<NewsArticle[]>([])
  const currentArticle = ref<NewsArticle | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const totalPages = ref(1)
  const currentPage = ref(1)
  const totalArticles = ref(0)

  const categories = computed(() => [
    { id: 'all', name: 'All News' },
    { id: 'update', name: 'Updates' },
    { id: 'event', name: 'Events' },
    { id: 'patch', name: 'Patch Notes' },
    { id: 'announcement', name: 'Announcements' }
  ])

  async function fetchArticles(filters: NewsFilters = {}) {
    loading.value = true
    error.value = null

    try {
      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 300))

      let filtered = [...mockArticles]

      // Filter by category
      if (filters.category && filters.category !== 'all') {
        filtered = filtered.filter(a => a.category === filters.category)
      }

      // Filter by search
      if (filters.search) {
        const searchLower = filters.search.toLowerCase()
        filtered = filtered.filter(a =>
          a.title.toLowerCase().includes(searchLower) ||
          a.excerpt.toLowerCase().includes(searchLower) ||
          a.tags.some(t => t.toLowerCase().includes(searchLower))
        )
      }

      // Sort by date (newest first)
      filtered.sort((a, b) =>
        new Date(b.publishedAt).getTime() - new Date(a.publishedAt).getTime()
      )

      // Pagination
      const page = filters.page || 1
      const limit = filters.limit || 6
      const startIndex = (page - 1) * limit
      const endIndex = startIndex + limit

      totalArticles.value = filtered.length
      totalPages.value = Math.ceil(filtered.length / limit)
      currentPage.value = page

      articles.value = filtered.slice(startIndex, endIndex)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch articles'
      console.error('News fetch error:', e)
    } finally {
      loading.value = false
    }
  }

  async function fetchArticle(slug: string) {
    loading.value = true
    error.value = null
    currentArticle.value = null

    try {
      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 200))

      const article = mockArticles.find(a => a.slug === slug)

      if (!article) {
        throw new Error('Article not found')
      }

      currentArticle.value = article
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch article'
      console.error('News fetch error:', e)
    } finally {
      loading.value = false
    }
  }

  function getRelatedArticles(article: NewsArticle, limit = 3): NewsArticle[] {
    return mockArticles
      .filter(a => a.id !== article.id)
      .filter(a =>
        a.category === article.category ||
        a.tags.some(t => article.tags.includes(t))
      )
      .slice(0, limit)
  }

  return {
    articles,
    currentArticle,
    loading,
    error,
    totalPages,
    currentPage,
    totalArticles,
    categories,
    fetchArticles,
    fetchArticle,
    getRelatedArticles
  }
})
