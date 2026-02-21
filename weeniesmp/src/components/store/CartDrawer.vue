<script setup lang="ts">
import { ref, watch, computed, toRef } from 'vue'
import { X, Minus, Plus, Trash2, ShoppingBag, ExternalLink, Loader2, Tag, Sparkles, ChevronDown, ChevronUp } from 'lucide-vue-next'
import { useTebexStore } from '@/stores/tebexStore'
import { useUiStore } from '@/stores/uiStore'
import { formatCurrency } from '@/utils/currency'
import { useFocusTrap } from '@/composables/useFocusTrap'

const tebexStore = useTebexStore()
const uiStore = useUiStore()

const drawerRef = ref<HTMLElement | null>(null)
const isOpen = toRef(uiStore, 'cartOpen')
const { handleEscape } = useFocusTrap(drawerRef, isOpen)

function handleKeyDown(e: KeyboardEvent) {
  if (handleEscape(e)) {
    uiStore.toggleCart()
  }
}
const updating = ref<number | null>(null)

// Coupon/Creator code state
const couponCode = ref('')
const creatorCode = ref('')
const applyingCoupon = ref(false)
const applyingCreator = ref(false)
const showDiscountSection = ref(false)

const hasAppliedCoupon = computed(() => tebexStore.basket?.coupons && tebexStore.basket.coupons.length > 0)
const hasAppliedCreatorCode = computed(() => !!tebexStore.basket?.creator_code)

async function handleApplyCoupon() {
  if (!couponCode.value.trim()) return
  applyingCoupon.value = true
  const success = await tebexStore.applyCoupon(couponCode.value.trim())
  applyingCoupon.value = false
  if (success) {
    couponCode.value = ''
  }
}

async function handleRemoveCoupon(couponId: string) {
  await tebexStore.removeCoupon(couponId)
}

async function handleApplyCreatorCode() {
  if (!creatorCode.value.trim()) return
  applyingCreator.value = true
  const success = await tebexStore.applyCreatorCode(creatorCode.value.trim())
  applyingCreator.value = false
  if (success) {
    creatorCode.value = ''
  }
}

async function handleRemoveCreatorCode() {
  await tebexStore.removeCreatorCode()
}

async function updateQuantity(packageId: number, newQuantity: number) {
  if (newQuantity <= 0) {
    await removeItem(packageId)
    return
  }
  updating.value = packageId
  await tebexStore.updateBasketQuantity(packageId, newQuantity)
  updating.value = null
}

async function removeItem(packageId: number) {
  updating.value = packageId
  await tebexStore.removeFromBasket(packageId)
  updating.value = null
}

function checkout() {
  tebexStore.checkout()
}

watch(() => tebexStore.basketItemCount, (newCount, oldCount) => {
  if (newCount > oldCount) {
    uiStore.openCart()
  }
})
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="uiStore.cartOpen"
        class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50"
        @click="uiStore.toggleCart"
      ></div>
    </Transition>

    <Transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="translate-x-full"
      enter-to-class="translate-x-0"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="translate-x-0"
      leave-to-class="translate-x-full"
    >
      <div
        v-if="uiStore.cartOpen"
        ref="drawerRef"
        role="dialog"
        aria-modal="true"
        aria-label="Shopping cart"
        class="fixed right-0 top-0 bottom-0 w-full max-w-md bg-weenie-dark border-l border-white/10 z-50 flex flex-col"
        @keydown="handleKeyDown"
      >
        <div class="flex items-center justify-between p-6 border-b border-white/10">
          <h2 class="text-xl font-bold text-white flex items-center gap-2">
            <ShoppingBag class="w-6 h-6 text-weenie-red" />
            Cart ({{ tebexStore.basketItemCount }})
          </h2>
          <button
            @click="uiStore.toggleCart"
            class="p-2 text-gray-400 hover:text-white transition-colors"
          >
            <X class="w-6 h-6" />
          </button>
        </div>

        <div class="flex-1 overflow-y-auto p-6">
          <div v-if="!tebexStore.basket?.packages?.length" class="text-center py-12">
            <ShoppingBag class="w-16 h-16 text-gray-600 mx-auto mb-4" />
            <p class="text-gray-400">Your cart is empty</p>
          </div>

          <div v-else class="space-y-4">
            <div
              v-for="item in tebexStore.basket.packages"
              :key="item.id"
              class="bg-weenie-darker/50 rounded-lg p-4"
            >
              <div class="flex items-start justify-between mb-2">
                <h3 class="font-semibold text-white">{{ item.name }}</h3>
                <button
                  @click="removeItem(item.id)"
                  :disabled="updating === item.id"
                  class="p-1 text-gray-500 hover:text-red-400 transition-colors disabled:opacity-50"
                >
                  <Loader2 v-if="updating === item.id" class="w-4 h-4 animate-spin" />
                  <Trash2 v-else class="w-4 h-4" />
                </button>
              </div>

              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <button
                    @click="updateQuantity(item.id, (item.in_basket?.quantity ?? 1) - 1)"
                    :disabled="updating === item.id"
                    class="p-1 bg-weenie-dark rounded text-gray-400 hover:text-white transition-colors disabled:opacity-50"
                  >
                    <Minus class="w-4 h-4" />
                  </button>
                  <span class="w-8 text-center text-white">{{ item.in_basket?.quantity ?? 0 }}</span>
                  <button
                    @click="updateQuantity(item.id, (item.in_basket?.quantity ?? 0) + 1)"
                    :disabled="updating === item.id"
                    class="p-1 bg-weenie-dark rounded text-gray-400 hover:text-white transition-colors disabled:opacity-50"
                  >
                    <Plus class="w-4 h-4" />
                  </button>
                </div>
                <span class="font-semibold text-weenie-gold">
                  {{ formatCurrency(tebexStore.basket?.currency) }}{{ (item.in_basket?.price ?? 0).toFixed(2) }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <div v-if="tebexStore.basket?.packages?.length" class="p-6 border-t border-white/10">
          <!-- Discount Codes Section -->
          <div class="mb-4">
            <button
              @click="showDiscountSection = !showDiscountSection"
              class="flex items-center justify-between w-full text-sm text-gray-400 hover:text-white transition-colors"
            >
              <span class="flex items-center gap-2">
                <Tag class="w-4 h-4" />
                Have a discount code?
              </span>
              <ChevronDown v-if="!showDiscountSection" class="w-4 h-4" />
              <ChevronUp v-else class="w-4 h-4" />
            </button>

            <Transition
              enter-active-class="transition duration-200 ease-out"
              enter-from-class="opacity-0 -translate-y-2"
              enter-to-class="opacity-100 translate-y-0"
              leave-active-class="transition duration-150 ease-in"
              leave-from-class="opacity-100 translate-y-0"
              leave-to-class="opacity-0 -translate-y-2"
            >
              <div v-if="showDiscountSection" class="mt-3 space-y-3">
                <!-- Applied Coupons -->
                <div v-if="hasAppliedCoupon" class="space-y-2">
                  <div
                    v-for="coupon in tebexStore.basket?.coupons"
                    :key="coupon.id"
                    class="flex items-center justify-between bg-green-500/10 border border-green-500/30 rounded-lg px-3 py-2"
                  >
                    <span class="text-green-400 text-sm flex items-center gap-2">
                      <Tag class="w-4 h-4" />
                      {{ coupon.code }}
                    </span>
                    <button
                      @click="handleRemoveCoupon(coupon.id)"
                      class="text-green-400 hover:text-red-400 transition-colors"
                    >
                      <X class="w-4 h-4" />
                    </button>
                  </div>
                </div>

                <!-- Coupon Input -->
                <div v-else class="flex gap-2">
                  <input
                    v-model="couponCode"
                    type="text"
                    placeholder="Coupon code"
                    class="flex-1 px-3 py-2 bg-weenie-darker border border-gray-700 rounded-lg text-white text-sm placeholder-gray-500 focus:outline-none focus:border-weenie-gold"
                    @keyup.enter="handleApplyCoupon"
                  />
                  <button
                    @click="handleApplyCoupon"
                    :disabled="applyingCoupon || !couponCode.trim()"
                    class="px-4 py-2 bg-weenie-dark hover:bg-weenie-gradient text-white text-sm font-medium rounded-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <Loader2 v-if="applyingCoupon" class="w-4 h-4 animate-spin" />
                    <span v-else>Apply</span>
                  </button>
                </div>

                <!-- Applied Creator Code -->
                <div v-if="hasAppliedCreatorCode" class="flex items-center justify-between bg-purple-500/10 border border-purple-500/30 rounded-lg px-3 py-2">
                  <span class="text-purple-400 text-sm flex items-center gap-2">
                    <Sparkles class="w-4 h-4" />
                    Creator: {{ tebexStore.basket?.creator_code }}
                  </span>
                  <button
                    @click="handleRemoveCreatorCode"
                    class="text-purple-400 hover:text-red-400 transition-colors"
                  >
                    <X class="w-4 h-4" />
                  </button>
                </div>

                <!-- Creator Code Input -->
                <div v-else class="flex gap-2">
                  <input
                    v-model="creatorCode"
                    type="text"
                    placeholder="Creator code"
                    class="flex-1 px-3 py-2 bg-weenie-darker border border-gray-700 rounded-lg text-white text-sm placeholder-gray-500 focus:outline-none focus:border-weenie-gold"
                    @keyup.enter="handleApplyCreatorCode"
                  />
                  <button
                    @click="handleApplyCreatorCode"
                    :disabled="applyingCreator || !creatorCode.trim()"
                    class="px-4 py-2 bg-weenie-dark hover:bg-weenie-gradient text-white text-sm font-medium rounded-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <Loader2 v-if="applyingCreator" class="w-4 h-4 animate-spin" />
                    <span v-else>Apply</span>
                  </button>
                </div>
              </div>
            </Transition>
          </div>

          <!-- Totals -->
          <div class="space-y-2 mb-4">
            <div v-if="tebexStore.basket && tebexStore.basket.base_price !== tebexStore.basketTotal" class="flex items-center justify-between text-sm">
              <span class="text-gray-500">Subtotal</span>
              <span class="text-gray-400">
                {{ formatCurrency(tebexStore.basket.currency) }}{{ (tebexStore.basket.base_price ?? 0).toFixed(2) }}
              </span>
            </div>
            <div v-if="tebexStore.basket && tebexStore.basket.sales_tax > 0" class="flex items-center justify-between text-sm">
              <span class="text-gray-500">Tax</span>
              <span class="text-gray-400">
                {{ formatCurrency(tebexStore.basket.currency) }}{{ (tebexStore.basket.sales_tax ?? 0).toFixed(2) }}
              </span>
            </div>
            <div class="flex items-center justify-between pt-2 border-t border-gray-700">
              <span class="text-gray-400">Total</span>
              <span class="text-2xl font-bold text-weenie-gold">
                {{ formatCurrency(tebexStore.basket?.currency) }}{{ tebexStore.basketTotal.toFixed(2) }}
              </span>
            </div>
          </div>

          <button
            @click="checkout"
            class="btn-primary w-full flex items-center justify-center gap-2"
          >
            Checkout
            <ExternalLink class="w-4 h-4" />
          </button>

          <p class="text-center text-xs text-gray-500 mt-3">
            Payments processed securely by Tebex
          </p>
        </div>
      </div>
    </Transition>
  </Teleport>

  <!-- Floating cart button -->
  <button
    v-if="tebexStore.basketItemCount > 0"
    @click="uiStore.toggleCart"
    class="fixed bottom-6 right-6 z-40 p-4 bg-weenie-gradient rounded-full shadow-lg shadow-weenie-red/30 text-white hover:scale-110 transition-transform"
  >
    <ShoppingBag class="w-6 h-6" />
    <span class="absolute -top-1 -right-1 w-5 h-5 bg-white text-weenie-red text-xs font-bold rounded-full flex items-center justify-center">
      {{ tebexStore.basketItemCount }}
    </span>
  </button>
</template>
