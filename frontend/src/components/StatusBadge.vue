<template>
  <van-tag :type="badgeType" size="medium" round>{{ text }}</van-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { BookStatus, BookStatusBadge, BookStatusText, WishStatus, WishStatusText } from '@/constants/enums'

const props = defineProps<{ status: string; kind?: 'book' | 'wish' }>()

const text = computed(() => {
  if (props.kind === 'wish') return WishStatusText[props.status] || props.status
  return BookStatusText[props.status] || props.status
})

const badgeType = computed(() => {
  if (props.kind === 'wish') {
    return props.status === WishStatus.OPEN ? 'success' : 'default'
  }
  return BookStatusBadge[props.status] || 'default'
})
</script>
