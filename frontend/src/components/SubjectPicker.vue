<template>
  <van-field
    :model-value="display"
    is-link
    readonly
    :label="label"
    :placeholder="placeholder"
    @click="show = true"
  />
  <van-popup v-model:show="show" round position="bottom">
    <van-picker
      :columns="columns"
      @confirm="onConfirm"
      @cancel="show = false"
      :model-value="currentValue ? [currentValue] : []"
    />
  </van-popup>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { SubjectCategoryOptions, SubjectCategoryText } from '@/constants/enums'

const props = defineProps<{ label?: string; placeholder?: string; modelValue?: string }>()
const emit = defineEmits<{ (e: 'update:modelValue', value: string): void }>()

const show = ref(false)
const columns = SubjectCategoryOptions.map((o) => ({ text: o.label, value: o.value }))

const currentValue = computed(() => props.modelValue || '')
const display = computed(() => (props.modelValue ? SubjectCategoryText[props.modelValue] : ''))

function onConfirm({ selectedOptions }: { selectedOptions: Array<{ value: string }> }) {
  if (selectedOptions.length) {
    emit('update:modelValue', selectedOptions[0].value)
  }
  show.value = false
}
</script>
