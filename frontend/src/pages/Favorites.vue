<template>
  <div class="page">
    <van-nav-bar title="我的收藏" left-arrow fixed placeholder @click-left="$router.back()" />
    <div v-if="list.length" class="book-list">
      <BookCard v-for="b in list" :key="b.id" :book="b" />
    </div>
    <EmptyState v-else description="还没有收藏书籍，去逛逛吧" icon="star-o" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { listFavorites } from '@/api/book'
import type { Book } from '@/types'
import BookCard from '@/components/BookCard.vue'
import EmptyState from '@/components/EmptyState.vue'

const list = ref<Book[]>([])

onMounted(async () => {
  try {
    const data = await listFavorites({ page: 1, page_size: 50 })
    list.value = data.list
  } catch {
    list.value = []
  }
})
</script>

<style scoped>
.book-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 12px;
}
</style>
