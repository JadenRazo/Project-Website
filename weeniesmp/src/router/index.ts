import { createRouter, createWebHistory } from 'vue-router'
import RootView from '@/views/RootView.vue'
import Store from '@/views/Store.vue'
import Rules from '@/views/Rules.vue'
import FAQ from '@/views/FAQ.vue'
import NotFound from '@/views/NotFound.vue'
import { isStoreSubdomain } from '@/utils/subdomain'

// Lazy-loaded views
const Leaderboards = () => import('@/views/Leaderboards.vue')
const News = () => import('@/views/News.vue')
const NewsDetail = () => import('@/views/NewsDetail.vue')
const Wiki = () => import('@/views/Wiki.vue')
const WikiPage = () => import('@/views/WikiPage.vue')
const Dashboard = () => import('@/views/Dashboard.vue')
const PlayerSearch = () => import('@/views/PlayerSearch.vue')
const PlayerProfile = () => import('@/views/PlayerProfile.vue')
const ServerMap = () => import('@/views/ServerMap.vue')
const Gallery = () => import('@/views/Gallery.vue')
const Appeals = () => import('@/views/Appeals.vue')

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Root',
      component: RootView
    },
    {
      path: '/store',
      name: 'Store',
      component: Store,
      beforeEnter: (_to, _from, next) => {
        // On store subdomain, redirect /store to / for clean URL
        if (isStoreSubdomain()) {
          next({ path: '/', replace: true })
        } else {
          next()
        }
      }
    },
    {
      path: '/rules',
      name: 'Rules',
      component: Rules
    },
    {
      path: '/faq',
      name: 'FAQ',
      component: FAQ
    },
    {
      path: '/leaderboards',
      name: 'Leaderboards',
      component: Leaderboards
    },
    {
      path: '/news',
      name: 'News',
      component: News
    },
    {
      path: '/news/:slug',
      name: 'NewsDetail',
      component: NewsDetail
    },
    {
      path: '/wiki',
      name: 'Wiki',
      component: Wiki
    },
    {
      path: '/wiki/:slug',
      name: 'WikiPage',
      component: WikiPage
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: Dashboard
    },
    {
      path: '/players',
      name: 'PlayerSearch',
      component: PlayerSearch
    },
    {
      path: '/players/:username',
      name: 'PlayerProfile',
      component: PlayerProfile
    },
    {
      path: '/map',
      name: 'ServerMap',
      component: ServerMap
    },
    {
      path: '/gallery',
      name: 'Gallery',
      component: Gallery
    },
    {
      path: '/appeals',
      name: 'Appeals',
      component: Appeals
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: NotFound
    }
  ],
  scrollBehavior(_to, _from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    return { top: 0 }
  }
})

// Update canonical URL on route change for better SEO
router.afterEach((to) => {
  let canonical = document.querySelector('link[rel="canonical"]')
  if (!canonical) {
    canonical = document.createElement('link')
    canonical.setAttribute('rel', 'canonical')
    document.head.appendChild(canonical)
  }
  canonical.setAttribute('href', `https://weeniesmp.net${to.path}`)
})

export default router
