<template>
  <span class="condition-tag" :style="style">{{ text }}</span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Condition, ConditionText } from '@/constants/enums'

const props = defineProps<{ condition: string }>()

const text = computed(() => ConditionText[props.condition] || props.condition || '未知')

const colors: Record<string, string> = {
  [Condition.BRAND_NEW]: '#07c160',
  [Condition.NINE_NEW]: '#1989fa',
  [Condition.SEVEN_NEW]: '#ff976a',
  [Condition.FIVE_NEW]: '#969799',
}

const style = computed(() => ({
  color: colors[props.condition] || '#969799',
  borderColor: colors[props.condition] || '#969799',
}))
</script>

<style scoped>
.condition-tag {
  display: inline-block;
  font-size: 12px;
  line-height: 18px;
  padding: 0 6px;
  border: 1px solid;
  border-radius: 4px;
  white-space: nowrap;
}
</style>
