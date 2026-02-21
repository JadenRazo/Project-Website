<script setup lang="ts">
import { ref } from 'vue'
import { HelpCircle, ChevronDown, ChevronUp, MessageCircle } from 'lucide-vue-next'

interface FAQItem {
  question: string
  answer: string
}

const faqs: FAQItem[] = [
  {
    question: "How long does rank delivery take?",
    answer: "Ranks are delivered automatically within 5 minutes. You must be online on the server to receive your rank. If you don't receive it, try relogging or contact support on Discord."
  },
  {
    question: "Do I need to be online to receive my purchase?",
    answer: "Yes, you need to be online on play.weeniesmp.net with the exact username you entered during checkout. Make sure your username is spelled correctly!"
  },
  {
    question: "Can I play on Bedrock Edition?",
    answer: "Yes! Bedrock players can connect via play.weeniesmp.net:19011. All ranks and purchases work on both Java and Bedrock editions."
  },
  {
    question: "How do I request a refund?",
    answer: "Refunds are handled through Tebex according to their refund policy. Contact us on Discord with your transaction ID and we'll help resolve any issues."
  },
  {
    question: "Are ranks permanent?",
    answer: "MVP and PRO ranks are lifetime purchases that never expire. VIP is a 1-month subscription that can be renewed."
  },
  {
    question: "What payment methods are accepted?",
    answer: "We accept PayPal, credit/debit cards (Visa, Mastercard, American Express), and various local payment methods through our secure payment processor Tebex."
  },
  {
    question: "I made a purchase but didn't receive it. What should I do?",
    answer: "First, make sure you're online on the server with the exact username you used during checkout. Try relogging. If you still don't receive your purchase after 10 minutes, contact us on Discord with your transaction ID."
  },
  {
    question: "Can I gift a rank to another player?",
    answer: "Currently, you need to enter the recipient's username during checkout to purchase a rank for another player. Make sure to enter their exact Minecraft username."
  }
]

const openIndex = ref<number | null>(null)

function toggle(index: number) {
  openIndex.value = openIndex.value === index ? null : index
}
</script>

<template>
  <div class="min-h-screen pt-24 pb-16 bg-weenie-darker">
    <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="text-center mb-12">
        <HelpCircle class="w-12 h-12 text-weenie-gold mx-auto mb-4" />
        <h1 class="text-4xl md:text-5xl font-bold mb-4">
          <span class="gradient-text">FAQ</span>
        </h1>
        <p class="text-gray-400 max-w-xl mx-auto">
          Frequently asked questions about purchases, ranks, and gameplay.
        </p>
      </div>

      <!-- FAQ Accordion -->
      <div class="space-y-4">
        <div
          v-for="(faq, index) in faqs"
          :key="index"
          class="card overflow-hidden"
        >
          <button
            @click="toggle(index)"
            class="w-full flex items-center justify-between p-5 text-left hover:bg-white/5 transition-colors"
          >
            <span class="font-medium text-white pr-4">{{ faq.question }}</span>
            <ChevronUp v-if="openIndex === index" class="w-5 h-5 text-weenie-gold flex-shrink-0" />
            <ChevronDown v-else class="w-5 h-5 text-gray-500 flex-shrink-0" />
          </button>
          <Transition
            enter-active-class="transition-all duration-300 ease-out"
            enter-from-class="max-h-0 opacity-0"
            enter-to-class="max-h-96 opacity-100"
            leave-active-class="transition-all duration-200 ease-in"
            leave-from-class="max-h-96 opacity-100"
            leave-to-class="max-h-0 opacity-0"
          >
            <div v-if="openIndex === index" class="overflow-hidden">
              <p class="px-5 pb-5 text-gray-400 border-t border-white/10 pt-4">
                {{ faq.answer }}
              </p>
            </div>
          </Transition>
        </div>
      </div>

      <!-- Still need help? -->
      <div class="mt-12 text-center">
        <div class="card p-8">
          <MessageCircle class="w-10 h-10 text-[#5865F2] mx-auto mb-4" />
          <h2 class="text-xl font-semibold text-white mb-2">Still have questions?</h2>
          <p class="text-gray-400 mb-6">Our team is happy to help you on Discord!</p>
          <a
            href="https://discord.com/invite/weeniesmp"
            target="_blank"
            rel="noopener noreferrer"
            class="inline-flex items-center gap-2 px-6 py-3 text-sm font-medium text-white bg-[#5865F2] hover:bg-[#4752C4] rounded-lg transition-colors"
          >
            Join Discord
          </a>
        </div>
      </div>
    </div>
  </div>
</template>
