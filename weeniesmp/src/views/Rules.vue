<script setup lang="ts">
import { Shield, AlertTriangle, Users, MessageSquare, Swords, Ban } from 'lucide-vue-next'

const rules = [
  {
    icon: Users,
    title: 'Respect Everyone',
    description: 'Treat all players and staff with respect. Harassment, bullying, or discrimination will not be tolerated.',
    severity: 'high'
  },
  {
    icon: Ban,
    title: 'No Cheating',
    description: 'Hacked clients, exploits, duping, or any form of cheating is strictly prohibited and will result in a permanent ban.',
    severity: 'high'
  },
  {
    icon: MessageSquare,
    title: 'Appropriate Chat',
    description: 'Keep chat family-friendly. No excessive swearing, spam, advertising, or inappropriate content.',
    severity: 'medium'
  },
  {
    icon: Swords,
    title: 'Fair PvP',
    description: 'Combat logging, TP killing, and spawn killing are not allowed. Fight fair and honor duels.',
    severity: 'medium'
  },
  {
    icon: Shield,
    title: 'No Griefing',
    description: 'Do not grief, steal from, or destroy other players builds even if unclaimed. Respect others property.',
    severity: 'high'
  },
  {
    icon: AlertTriangle,
    title: 'No Lag Machines',
    description: 'Intentionally causing server lag through redstone, entities, or other means is prohibited.',
    severity: 'high'
  }
]

const severityColors: Record<string, string> = {
  high: 'border-red-500/50 bg-red-500/5',
  medium: 'border-yellow-500/50 bg-yellow-500/5',
  low: 'border-green-500/50 bg-green-500/5'
}
</script>

<template>
  <div class="min-h-screen pt-24 pb-16 bg-weenie-darker">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="text-center mb-12">
        <Shield class="w-12 h-12 text-weenie-red mx-auto mb-4" />
        <h1 class="text-4xl md:text-5xl font-bold mb-4">
          Server <span class="gradient-text">Rules</span>
        </h1>
        <p class="text-gray-400 max-w-xl mx-auto">
          Please read and follow these rules to ensure a great experience for everyone.
        </p>
      </div>

      <div class="space-y-4">
        <div
          v-for="(rule, index) in rules"
          :key="rule.title"
          class="card border-l-4"
          :class="severityColors[rule.severity]"
        >
          <div class="flex items-start gap-4">
            <div class="flex-shrink-0 w-10 h-10 rounded-lg bg-weenie-dark flex items-center justify-center">
              <span class="text-weenie-gold font-bold">{{ index + 1 }}</span>
            </div>
            <div class="flex-1">
              <div class="flex items-center gap-3 mb-2">
                <component :is="rule.icon" class="w-5 h-5 text-weenie-red" />
                <h3 class="text-lg font-semibold text-white">{{ rule.title }}</h3>
              </div>
              <p class="text-gray-400">{{ rule.description }}</p>
            </div>
          </div>
        </div>
      </div>

      <div class="mt-12 card bg-weenie-red/10 border-weenie-red/30">
        <div class="flex items-start gap-4">
          <AlertTriangle class="w-6 h-6 text-weenie-red flex-shrink-0 mt-1" />
          <div>
            <h3 class="font-semibold text-white mb-2">Important Notice</h3>
            <p class="text-gray-300 text-sm">
              Breaking these rules may result in warnings, temporary mutes/bans, or permanent bans
              depending on the severity. Staff decisions are final. If you believe a punishment was
              unfair, appeal on our Discord server.
            </p>
          </div>
        </div>
      </div>

      <div class="mt-8 text-center">
        <a
          href="https://discord.gg/weeniesmp"
          target="_blank"
          rel="noopener noreferrer"
          class="btn-secondary"
        >
          Appeal on Discord
        </a>
      </div>
    </div>
  </div>
</template>
