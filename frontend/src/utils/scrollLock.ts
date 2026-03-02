let isLocked = false
let lenisInstance: any = null

export function setLenisInstance(lenis: any): void {
  lenisInstance = lenis
}

export function lockScroll(): void {
  if (isLocked) return
  isLocked = true

  if (lenisInstance) {
    lenisInstance.stop()
  }

  document.documentElement.style.overflow = 'hidden'
  document.body.style.overflow = 'hidden'
}

export function unlockScroll(): void {
  if (!isLocked) return
  isLocked = false

  document.documentElement.style.overflow = ''
  document.body.style.overflow = ''

  if (lenisInstance) {
    lenisInstance.scrollTo(0, { immediate: true })
    lenisInstance.start()
  }
}

export function isScrollLocked(): boolean {
  return isLocked
}
