<template>
  <div class="page">
    <van-nav-bar title="我发布的书籍" left-arrow fixed placeholder @click-left="$router.back()" />
    <van-tabs v-model:active="activeTab" @change="onFilter">
      <van-tab title="在售" name="on_sale" />
      <van-tab title="已预约" name="reserved" />
      <van-tab title="借出" name="lent_out" />
      <van-tab title="已售出" name="sold" />
    </van-tabs>

    <div v-if="list.length" class="book-list">
      <BookCard v-for="b in list" :key="b.id" :book="b" show-status />
    </div>
    <EmptyState v-else description="暂无相关书籍" />

    <van-button round block type="primary" class="fab" icon="plus" @click="$router.push('/publish-book')">发布新书</van-button>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { listBooks } from '@/api/book'
import { useAuthStore } from '@/stores/authStore'
import type { Book } from '@/types'
import BookCard from '@/components/BookCard.vue'
import EmptyState from '@/components/EmptyState.vue'

const auth = useAuthStore()
const activeTab = ref('on_sale')
const list = ref<Book[]>([])

onMounted(() => {
  if (auth.user) {
    onFilter()
  }
})

async function onFilter() {
  if (!auth.user) return
  try {
    const data = await listBooks({ seller_id: auth.user.id, status: activeTab.value, page: 1, page_size: 50 })
    list.value = data.list
  } catch {
    list.value = []
  }
}
</script>

<style scoped>
.book-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 12px;
}
.fab {
  position: fixed;
  right: 16px;
  bottom: 30px;
  width: 120px;
  z-index: 10;
}
</style>
