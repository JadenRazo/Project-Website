let isLocked = false

export function lockScroll(): void {
  if (isLocked) return
  isLocked = true
  document.body.style.overflow = 'hidden'
}

export function unlockScroll(): void {
  if (!isLocked) return
  isLocked = false

  document.body.style.overflow = ''

}

export function isScrollLocked(): boolean {
  return isLocked
}
