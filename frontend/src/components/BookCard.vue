<template>
  <div class="book-card" @click="$router.push(`/books/${book.id}`)">
    <div class="cover">
      <img v-if="firstImage" :src="firstImage" alt="cover" />
      <van-icon v-else name="book-o" size="40" color="#c8c9cc" />
      <StatusBadge v-if="showStatus" :status="book.status" class="status-badge" />
    </div>
    <div class="info">
      <div class="title van-ellipsis">{{ book.title }}</div>
      <div class="meta van-ellipsis">{{ book.author || '佚名' }}</div>
      <div class="tags">
        <ConditionTag :condition="book.condition" />
        <van-tag plain type="primary">{{ book.subject_text }}</van-tag>
        <van-tag v-if="book.trade_type_text" plain>{{ book.trade_type_text }}</van-tag>
      </div>
      <div class="bottom">
        <span class="price">{{ formatPrice(book.price) }}</span>
        <span class="views"><van-icon name="eye-o" /> {{ book.view_count }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Book } from '@/types'
import { formatPrice, formatImages } from '@/utils/format'
import StatusBadge from '@/components/StatusBadge.vue'
import ConditionTag from '@/components/ConditionTag.vue'

const props = withDefaults(defineProps<{ book: Book; showStatus?: boolean }>(), {
  showStatus: false,
})

const firstImage = computed(() => formatImages(props.book.images)[0] || '')
</script>

<style scoped>
.book-card {
  display: flex;
  background: #fff;
  border-radius: 8px;
  padding: 10px;
  gap: 10px;
}
.cover {
  position: relative;
  width: 84px;
  height: 110px;
  background: #f2f3f5;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  overflow: hidden;
}
.cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.status-badge {
  position: absolute;
  top: 4px;
  left: 4px;
}
.info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.title {
  font-size: 15px;
  font-weight: 600;
  color: #323233;
}
.meta {
  font-size: 12px;
  color: #969799;
}
.tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
  align-items: center;
}
.bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
}
.price {
  color: #ee0a24;
  font-size: 16px;
  font-weight: 600;
}
.views {
  font-size: 12px;
  color: #969799;
}
</style>
