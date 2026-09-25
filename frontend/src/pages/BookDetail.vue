<template>
  <div class="page detail-page">
    <van-nav-bar title="书籍详情" left-arrow fixed placeholder @click-left="$router.back()" />

    <template v-if="book">
      <van-swipe v-if="images.length" class="detail-swipe" :autoplay="0" indicator-color="white">
        <van-swipe-item v-for="img in images" :key="img">
          <img :src="img" class="detail-img" alt="book" />
        </van-swipe-item>
      </van-swipe>
      <div v-else class="no-img"><van-icon name="book-o" size="48" color="#c8c9cc" /></div>

      <div class="detail-card">
        <div class="detail-title">{{ book.title }}</div>
        <div class="detail-price">
          <span class="price">{{ formatPrice(book.price) }}</span>
          <span v-if="book.original_price" class="origin">原价 {{ formatPrice(book.original_price) }}</span>
        </div>
        <div class="detail-tags">
          <ConditionTag :condition="book.condition" />
          <StatusBadge :status="book.status" />
          <van-tag plain type="primary">{{ book.subject_text }}</van-tag>
          <van-tag plain>{{ book.trade_type_text }}</van-tag>
        </div>
        <van-cell-group inset class="meta-group">
          <van-cell title="作者" :value="book.author || '佚名'" />
          <van-cell title="ISBN" :value="book.isbn || '-'" />
          <van-cell title="课程" :value="book.course_name || '-'" />
          <van-cell title="校区" :value="book.campus || '-'" />
          <van-cell title="浏览量" :value="String(book.view_count)" />
        </van-cell-group>
        <div class="desc">
          <div class="desc-title">书籍描述</div>
          <div class="desc-body">{{ book.description || '暂无描述' }}</div>
        </div>
        <van-cell-group inset>
          <van-cell :title="book.seller?.name || '卖家'" :label="`${book.seller?.department || ''} · ${book.seller?.campus || ''}`" @click="goSeller">
            <template #icon><van-image round width="36" height="36" :src="book.seller?.avatar_url || fallbackAvatar" /></template>
            <template #value><van-icon name="arrow" /></template>
          </van-cell>
        </van-cell-group>
      </div>

      <div class="action-bar">
        <van-button plain type="warning" icon="star-o" @click="onFavorite">
          {{ book.is_favorite ? '已收藏' : '收藏' }}
        </van-button>
        <template v-if="isOwner">
          <van-button v-if="book.status === 'on_sale'" plain type="primary" icon="edit" @click="$router.push({ name: 'publish-book', query: { id: String(book.id) } })">编辑</van-button>
          <van-button v-if="book.status === 'on_sale'" type="danger" icon="delete-o" @click="onDelete">下架</van-button>
          <van-button v-if="book.status === 'reserved'" type="success" icon="passed" @click="onSold">确认售出</van-button>
          <van-button v-if="book.status === 'reserved'" plain type="default" @click="onCancelReserve">取消预约</van-button>
        </template>
        <template v-else>
          <van-button v-if="book.status === 'on_sale'" type="primary" icon="chat-o" @click="onContact">联系卖家</van-button>
          <van-button v-if="book.status === 'on_sale'" type="warning" plain icon="clock-o" @click="onReserve">预约</van-button>
          <van-button v-if="book.status === 'reserved'" plain type="default" icon="clock-o" @click="onCancelReserve">取消预约</van-button>
        </template>
      </div>
    </template>
    <EmptyState v-else description="书籍不存在或已删除" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showConfirmDialog, showToast } from 'vant'
import { getBook, addFavorite, removeFavorite, reserveBook, cancelReserveBook, markSold, deleteBook } from '@/api/book'
import { createConversationFromBook } from '@/api/conversation'
import type { Book } from '@/types'
import { formatPrice, formatImages } from '@/utils/format'
import { useAuthStore } from '@/stores/authStore'
import ConditionTag from '@/components/ConditionTag.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import EmptyState from '@/components/EmptyState.vue'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const book = ref<Book | null>(null)
const fallbackAvatar = ''

const images = computed(() => (book.value ? formatImages(book.value.images) : []))
const isOwner = computed(() => !!book.value && !!auth.user && book.value.seller_id === auth.user.id)

onMounted(async () => {
  const id = Number(route.params.id)
  try {
    book.value = await getBook(id)
  } catch {
    book.value = null
  }
})

async function onFavorite() {
  if (!auth.isLoggedIn) {
    router.push('/login')
    return
  }
  if (!book.value) return
  try {
    if (book.value.is_favorite) {
      await removeFavorite(book.value.id)
      showToast('已取消收藏')
    } else {
      await addFavorite(book.value.id)
      showToast('收藏成功')
    }
    book.value.is_favorite = !book.value.is_favorite
    book.value.favorite_count = Math.max(0, (book.value.favorite_count || 0) + (book.value.is_favorite ? 1 : -1))
  } catch {
    /* toast handled by interceptor */
  }
}

async function onReserve() {
  if (!auth.isLoggedIn) {
    router.push('/login')
    return
  }
  if (!book.value) return
  try {
    await showConfirmDialog({ title: '预约', message: '确认预约这本书吗？' })
    book.value = await reserveBook(book.value.id)
    showToast('预约成功')
  } catch {
    /* cancelled or error */
  }
}

async function onCancelReserve() {
  if (!book.value) return
  try {
    await showConfirmDialog({ title: '取消预约', message: '确认取消预约吗？' })
    book.value = await cancelReserveBook(book.value.id)
    showToast('已取消预约')
  } catch {
    /* cancelled */
  }
}

async function onSold() {
  if (!book.value) return
  try {
    await showConfirmDialog({ title: '确认售出', message: '确认该书籍已完成交易并售出？' })
    book.value = await markSold(book.value.id)
    showToast('已标记售出')
  } catch {
    /* cancelled */
  }
}

async function onDelete() {
  if (!book.value) return
  try {
    await showConfirmDialog({ title: '下架', message: '确认下架该书籍？' })
    await deleteBook(book.value.id)
    showToast('已下架')
    router.back()
  } catch {
    /* cancelled */
  }
}

async function onContact() {
  if (!auth.isLoggedIn) {
    router.push('/login')
    return
  }
  if (!book.value) return
  try {
    const conv = await createConversationFromBook(book.value.id, '你好，我对这本书感兴趣，还在吗？')
    router.push(`/chat/${conv.id}`)
  } catch {
    /* toast handled */
  }
}

function goSeller() {
  if (book.value?.seller_id) {
    router.push(`/profile?userId=${book.value.seller_id}`)
  }
}
</script>

<style scoped>
.detail-swipe {
  height: 260px;
}
.detail-img {
  width: 100%;
  height: 260px;
  object-fit: cover;
}
.no-img {
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f2f3f5;
}
.detail-card {
  background: #fff;
  border-radius: 8px;
  padding: 14px;
  margin: 12px 0;
}
.detail-title {
  font-size: 18px;
  font-weight: 700;
}
.detail-price {
  margin: 8px 0;
}
.price {
  color: #ee0a24;
  font-size: 22px;
  font-weight: 700;
}
.origin {
  margin-left: 8px;
  color: #969799;
  font-size: 13px;
  text-decoration: line-through;
}
.detail-tags {
  display: flex;
  gap: 6px;
  align-items: center;
  flex-wrap: wrap;
}
.meta-group {
  margin: 12px 0;
}
.desc {
  margin-top: 12px;
}
.desc-title {
  font-weight: 600;
  margin-bottom: 6px;
}
.desc-body {
  color: #646566;
  font-size: 14px;
  line-height: 1.6;
  white-space: pre-wrap;
}
.action-bar {
  position: sticky;
  bottom: 0;
  display: flex;
  gap: 10px;
  padding: 10px 12px;
  background: #fff;
  border-top: 1px solid #ebedf0;
}
.action-bar .van-button {
  flex: 1;
}
</style>
