import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUiStore = defineStore('ui', () => {
  const cartOpen = ref(false)

  function openCart() {
    cartOpen.value = true
  }

  function closeCart() {
    cartOpen.value = false
  }

  function toggleCart() {
    cartOpen.value = !cartOpen.value
  }

  return {
    cartOpen,
    openCart,
    closeCart,
    toggleCart
  }
})
