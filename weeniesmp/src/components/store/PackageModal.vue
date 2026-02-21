<script setup lang="ts">
import { ref, watch, computed, toRef } from 'vue'
import { X, ShoppingCart, Check, Loader2, Minus, Plus, Crown, Sparkles, Package, Coins, Tag } from 'lucide-vue-next'
import DOMPurify from 'dompurify'
import { useTebexStore, type TebexPackage } from '@/stores/tebexStore'
import { useUiStore } from '@/stores/uiStore'
import { formatCurrency } from '@/utils/currency'
import { useFocusTrap } from '@/composables/useFocusTrap'

const props = defineProps<{
  package: TebexPackage | null
  open: boolean
}>()

const modalRef = ref<HTMLElement | null>(null)
const isOpen = toRef(props, 'open')
const { handleEscape } = useFocusTrap(modalRef, isOpen)

const emit = defineEmits<{
  close: []
}>()

const tebexStore = useTebexStore()
const uiStore = useUiStore()
const quantity = ref(1)
const adding = ref(false)
const added = ref(false)

// Sanitized description to prevent XSS
const sanitizedDescription = computed(() => {
  if (!props.package?.description) return ''
  return DOMPurify.sanitize(props.package.description, {
    ALLOWED_TAGS: ['p', 'br', 'strong', 'b', 'em', 'i', 'u', 'ul', 'ol', 'li', 'span', 'div', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6'],
    ALLOWED_ATTR: ['class', 'style']
  })
})

// Reset quantity when package changes OR when modal opens
watch([() => props.package, () => props.open], ([newPkg, newOpen], [oldPkg, oldOpen]) => {
  // Reset when modal opens (even for same package)
  if (newOpen && !oldOpen) {
    quantity.value = 1
    added.value = false
  }
  // Also reset when package changes
  if (newPkg !== oldPkg) {
    quantity.value = 1
    added.value = false
  }
})

const getIcon = (categoryName: string) => {
  const lower = categoryName?.toLowerCase() || ''
  if (lower.includes('rank')) return Crown
  if (lower.includes('cosmetic')) return Sparkles
  if (lower.includes('crate') || lower.includes('key')) return Package
  if (lower.includes('coin') || lower.includes('currency')) return Coins
  return Tag
}

function incrementQty() {
  quantity.value++
}

function decrementQty() {
  if (quantity.value > 1) {
    quantity.value--
  }
}

async function addToCart() {
  if (!props.package) return
  adding.value = true
  const success = await tebexStore.addToBasket(props.package.id, quantity.value)
  adding.value = false

  if (success) {
    added.value = true
    setTimeout(() => {
      emit('close')
      uiStore.openCart()
    }, 800)
  }
}

function handleBackdropClick(e: MouseEvent) {
  if (e.target === e.currentTarget) {
    emit('close')
  }
}

function handleKeyDown(e: KeyboardEvent) {
  if (handleEscape(e)) {
    emit('close')
  }
}
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="open && package"
        ref="modalRef"
        role="dialog"
        aria-modal="true"
        :aria-labelledby="'modal-title-' + package.id"
        class="fixed inset-0 z-[90] flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm"
        @click="handleBackdropClick"
        @keydown="handleKeyDown"
      >
        <Transition
          enter-active-class="transition duration-200 ease-out"
          enter-from-class="opacity-0 scale-95"
          enter-to-class="opacity-100 scale-100"
          leave-active-class="transition duration-150 ease-in"
          leave-from-class="opacity-100 scale-100"
          leave-to-class="opacity-0 scale-95"
        >
          <div
            v-if="open"
            class="relative w-full max-w-2xl max-h-[90vh] overflow-y-auto bg-weenie-dark border border-gray-700 rounded-2xl shadow-2xl"
          >
            <!-- Close Button -->
            <button
              @click="emit('close')"
              class="absolute top-4 right-4 p-2 text-gray-400 hover:text-white hover:bg-white/10 rounded-lg transition-colors z-10"
            >
              <X class="w-5 h-5" />
            </button>

            <!-- Package Image -->
            <div class="relative h-48 bg-weenie-gradient flex items-center justify-center overflow-hidden rounded-t-2xl">
              <img
                v-if="package.image"
                :src="package.image"
                :alt="package.name"
                class="w-full h-full object-cover"
              />
              <component
                v-else
                :is="getIcon(package.category?.name || '')"
                class="w-20 h-20 text-white/80"
              />
              <!-- Category Badge -->
              <div
                v-if="package.category"
                class="absolute bottom-4 left-4 px-3 py-1 bg-black/50 backdrop-blur-sm rounded-full text-sm text-white"
              >
                {{ package.category.name }}
              </div>
            </div>

            <!-- Content -->
            <div class="p-6">
              <!-- Header -->
              <div class="flex items-start justify-between mb-4">
                <h2 :id="'modal-title-' + package.id" class="text-2xl font-bold text-white pr-8">{{ package.name }}</h2>
                <div class="text-right flex-shrink-0">
                  <span
                    v-if="package.sales_price && package.sales_price < package.base_price"
                    class="text-sm text-gray-500 line-through block"
                  >
                    {{ formatCurrency(package.currency) }}{{ package.base_price.toFixed(2) }}
                  </span>
                  <span class="text-3xl font-bold text-weenie-gold">
                    {{ formatCurrency(package.currency) }}{{ package.total_price.toFixed(2) }}
                  </span>
                </div>
              </div>

              <!-- Description (sanitized for XSS protection) -->
              <div
                class="prose prose-invert prose-sm max-w-none mb-6 text-gray-300"
                v-html="sanitizedDescription"
              />

              <!-- Quantity & Add to Cart -->
              <div class="flex items-center gap-4 pt-4 border-t border-gray-700">
                <!-- Quantity Selector -->
                <div class="flex items-center gap-2">
                  <span class="text-gray-400 text-sm">Qty:</span>
                  <div class="flex items-center bg-weenie-darker rounded-lg">
                    <button
                      @click="decrementQty"
                      :disabled="quantity <= 1"
                      class="p-2 text-gray-400 hover:text-white disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                    >
                      <Minus class="w-4 h-4" />
                    </button>
                    <span class="w-10 text-center text-white font-semibold">{{ quantity }}</span>
                    <button
                      @click="incrementQty"
                      class="p-2 text-gray-400 hover:text-white transition-colors"
                    >
                      <Plus class="w-4 h-4" />
                    </button>
                  </div>
                </div>

                <!-- Add to Cart Button -->
                <button
                  @click="addToCart"
                  :disabled="adding"
                  class="flex-1 py-3 rounded-lg font-semibold transition-all duration-300 flex items-center justify-center gap-2 disabled:opacity-50"
                  :class="
                    added
                      ? 'bg-green-500 text-white'
                      : 'bg-weenie-gradient text-white hover:opacity-90'
                  "
                >
                  <Loader2 v-if="adding" class="w-5 h-5 animate-spin" />
                  <Check v-else-if="added" class="w-5 h-5" />
                  <ShoppingCart v-else class="w-5 h-5" />
                  {{ adding ? 'Adding...' : added ? 'Added!' : `Add to Cart - ${formatCurrency(package.currency)}${(package.total_price * quantity).toFixed(2)}` }}
                </button>
              </div>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
/* Custom prose styles for description */
.prose :deep(ul) {
  list-style-type: disc;
  padding-left: 1.5rem;
}

.prose :deep(li) {
  margin-bottom: 0.25rem;
}

.prose :deep(p) {
  margin-bottom: 0.75rem;
}

.prose :deep(strong) {
  color: white;
}
</style>
