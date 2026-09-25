<template>
  <div class="image-uploader">
    <van-uploader
      v-model="fileList"
      :max-count="maxCount"
      :after-read="onAfterRead"
      @delete="onDelete"
      accept="image/*"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { showLoadingToast, closeToast, showToast } from 'vant'
import { uploadFile } from '@/utils/upload'
import type { UploaderFileListItem } from 'vant'

const props = withDefaults(defineProps<{ modelValue: string[]; maxCount?: number }>(), {
  modelValue: () => [],
  maxCount: 5,
})

const emit = defineEmits<{ (e: 'update:modelValue', value: string[]): void }>()

const fileList = ref<UploaderFileListItem[]>(props.modelValue.map((url) => ({ url, isImage: true })))

watch(
  () => props.modelValue,
  (val) => {
    const urls = fileList.value.map((f) => f.url || '').filter(Boolean)
    if (JSON.stringify(urls) !== JSON.stringify(val)) {
      fileList.value = val.map((url) => ({ url, isImage: true }))
    }
  },
)

async function onAfterRead(item: UploaderFileListItem | UploaderFileListItem[]) {
  const file = (Array.isArray(item) ? item[0] : item)?.file
  if (!file) {
    showToast('文件读取失败')
    return
  }
  showLoadingToast({ message: '上传中...', forbidClick: true })
  try {
    const res = await uploadFile(file)
    closeToast()
    emit('update:modelValue', [...props.modelValue, res.url])
  } catch {
    closeToast()
    showToast('上传失败')
  }
}

function onDelete(file: UploaderFileListItem) {
  const url = file.url || ''
  emit(
    'update:modelValue',
    props.modelValue.filter((u) => u !== url),
  )
}
</script>
