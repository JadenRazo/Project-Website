let isLocked = false
let lenisInstance: any = null

export function setLenisInstance(lenis: any): void {
  lenisInstance = lenis
}

export function resetScrollLock(): void {
  isLocked = false
  document.body.style.overflow = ''
}

export function lockScroll(): void {
  if (isLocked) return
  isLocked = true

  if (lenisInstance) {
    lenisInstance.stop()
  } else {
    document.body.style.overflow = 'hidden'
  }
}

export function unlockScroll(): void {
  if (!isLocked) return
  isLocked = false

  document.body.style.overflow = ''

  if (lenisInstance) {
    lenisInstance.start()
  }
}

export function isScrollLocked(): boolean {
  return isLocked
}
