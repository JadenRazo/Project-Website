let scrollPosition = 0
let isLocked = false
let touchMoveHandler: ((e: TouchEvent) => void) | null = null

export function lockScroll(): void {
  if (isLocked) return
  isLocked = true

  scrollPosition = window.scrollY

  const scrollbarWidth = window.innerWidth - document.documentElement.clientWidth

  document.body.style.position = 'fixed'
  document.body.style.top = `-${scrollPosition}px`
  document.body.style.left = '0'
  document.body.style.right = '0'
  document.body.style.width = '100%'
  document.body.style.overflow = 'hidden'

  if (scrollbarWidth > 0) {
    document.body.style.paddingRight = `${scrollbarWidth}px`
  }

  touchMoveHandler = (e: TouchEvent) => {
    const target = e.target as HTMLElement
    const isScrollable = target.closest('[data-scroll-lock-scrollable]')
    if (!isScrollable) {
      e.preventDefault()
    }
  }

  document.addEventListener('touchmove', touchMoveHandler, { passive: false })
}

export function unlockScroll(): void {
  if (!isLocked) return
  isLocked = false

  document.body.style.position = ''
  document.body.style.top = ''
  document.body.style.left = ''
  document.body.style.right = ''
  document.body.style.width = ''
  document.body.style.overflow = ''
  document.body.style.paddingRight = ''

  window.scrollTo(0, scrollPosition)

  if (touchMoveHandler) {
    document.removeEventListener('touchmove', touchMoveHandler)
    touchMoveHandler = null
  }
}

export function isScrollLocked(): boolean {
  return isLocked
}
