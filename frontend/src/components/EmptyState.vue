<template>
  <div class="empty-state">
    <van-empty :image-size="90" :description="description">
      <template #image>
        <van-icon :name="icon" size="56" color="#c8c9cc" />
      </template>
      <template v-if="$slots.default" #default>
        <slot />
      </template>
    </van-empty>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{ description?: string; icon?: string }>(), {
  description: '暂无数据',
  icon: 'records-o',
})
</script>

<style scoped>
.empty-state {
  padding: 24px 0;
}
</style>
