import { ref, onMounted, onUnmounted, watch, type Ref } from 'vue'

const FOCUSABLE_SELECTOR = [
  'a[href]',
  'button:not([disabled])',
  'input:not([disabled])',
  'select:not([disabled])',
  'textarea:not([disabled])',
  '[tabindex]:not([tabindex="-1"])'
].join(', ')

export function useFocusTrap(
  containerRef: Ref<HTMLElement | null>,
  isActive: Ref<boolean>
) {
  const previouslyFocused = ref<HTMLElement | null>(null)

  function getFocusableElements(): HTMLElement[] {
    if (!containerRef.value) return []
    return Array.from(containerRef.value.querySelectorAll<HTMLElement>(FOCUSABLE_SELECTOR))
  }

  function handleKeyDown(event: KeyboardEvent) {
    if (event.key !== 'Tab' || !isActive.value) return

    const focusable = getFocusableElements()
    if (focusable.length === 0) return

    const first = focusable[0]
    const last = focusable[focusable.length - 1]

    if (event.shiftKey && document.activeElement === first) {
      event.preventDefault()
      last.focus()
    } else if (!event.shiftKey && document.activeElement === last) {
      event.preventDefault()
      first.focus()
    }
  }

  function handleEscape(event: KeyboardEvent) {
    if (event.key === 'Escape' && isActive.value) {
      return true
    }
    return false
  }

  function activate() {
    previouslyFocused.value = document.activeElement as HTMLElement
    const focusable = getFocusableElements()
    if (focusable.length > 0) {
      focusable[0].focus()
    }
  }

  function deactivate() {
    if (previouslyFocused.value && typeof previouslyFocused.value.focus === 'function') {
      previouslyFocused.value.focus()
    }
  }

  watch(isActive, (active) => {
    if (active) {
      activate()
      document.addEventListener('keydown', handleKeyDown)
    } else {
      deactivate()
      document.removeEventListener('keydown', handleKeyDown)
    }
  })

  onMounted(() => {
    if (isActive.value) {
      activate()
      document.addEventListener('keydown', handleKeyDown)
    }
  })

  onUnmounted(() => {
    document.removeEventListener('keydown', handleKeyDown)
    deactivate()
  })

  return {
    handleEscape
  }
}
