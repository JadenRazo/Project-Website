// Subdomain constants
export const SUBDOMAINS = {
  STORE: 'store.weeniesmp.net',
  MAIN: 'weeniesmp.net'
} as const

// URLs for cross-domain navigation
export const mainSiteUrl = 'https://weeniesmp.net'
export const storeUrl = 'https://store.weeniesmp.net'

// Cache for subdomain detection (computed once per page load)
let cachedHostname: string | null = null
let cachedIsStore: boolean | null = null

function getHostname(): string {
  if (cachedHostname === null) {
    cachedHostname = typeof window !== 'undefined' ? window.location.hostname : ''
  }
  return cachedHostname
}

export function isStoreSubdomain(): boolean {
  if (cachedIsStore === null) {
    cachedIsStore = getHostname() === SUBDOMAINS.STORE
  }
  return cachedIsStore
}

export function isMainSite(): boolean {
  const hostname = getHostname()
  return hostname === SUBDOMAINS.MAIN || hostname === 'www.weeniesmp.net'
}

export function getSubdomainUrl(subdomain: keyof typeof SUBDOMAINS): string {
  switch (subdomain) {
    case 'STORE':
      return storeUrl
    case 'MAIN':
      return mainSiteUrl
    default:
      return mainSiteUrl
  }
}
