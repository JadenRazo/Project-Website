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
  } else {
    document.body.style.overflow = 'hidden'
  }
}

export function unlockScroll(): void {
  if (!isLocked) return
  isLocked = false

  if (lenisInstance) {
    lenisInstance.start()
  } else {
    document.body.style.overflow = ''
  }
}

export function isScrollLocked(): boolean {
  return isLocked
}
