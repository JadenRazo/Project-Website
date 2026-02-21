<script setup lang="ts">
interface Props {
  width?: string
  height?: string
  rounded?: 'none' | 'sm' | 'md' | 'lg' | 'full'
  animate?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  width: '100%',
  height: '1rem',
  rounded: 'md',
  animate: true
})

const roundedClasses = {
  none: 'rounded-none',
  sm: 'rounded-sm',
  md: 'rounded-md',
  lg: 'rounded-lg',
  full: 'rounded-full'
}
</script>

<template>
  <div
    aria-hidden="true"
    class="bg-white/5"
    :class="[
      roundedClasses[props.rounded],
      props.animate ? 'animate-pulse' : ''
    ]"
    :style="{
      width: props.width,
      height: props.height
    }"
  />
</template>
