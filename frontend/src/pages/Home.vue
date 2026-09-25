<template>
  <div class="page home-page">
    <van-nav-bar title="校园二手书交易平台" fixed placeholder />
    <van-search v-model="keyword" placeholder="搜索书名 / 作者 / ISBN" @search="onSearch" />

    <div class="banner">
      <van-swipe :autoplay="3000" indicator-color="white">
        <van-swipe-item v-for="b in banners" :key="b.title">
          <div class="banner-item" :style="{ background: b.bg }">
            <div class="banner-text">
              <div class="banner-title">{{ b.title }}</div>
              <div class="banner-sub">{{ b.sub }}</div>
            </div>
          </div>
        </van-swipe-item>
      </van-swipe>
    </div>

    <div class="section">
      <div class="section-title">
        <span>学科分类</span>
        <span class="more" @click="$router.push('/books')">全部 <van-icon name="arrow" /></span>
      </div>
      <div class="category-grid">
        <div v-for="c in categories" :key="c.value" class="category-item" @click="goCategory(c.value)">
          <div class="category-icon" :style="{ background: c.bg }">
            <van-icon :name="c.icon" size="22" color="#fff" />
          </div>
          <div class="category-label">{{ c.label }}</div>
        </div>
      </div>
    </div>

    <div v-if="isLoggedIn && recommendations.length" class="section">
      <div class="section-title"><span>同院系推荐</span></div>
      <van-cell-group inset>
        <van-cell v-for="b in recommendations" :key="b.id" :title="b.title" :label="`${b.seller?.name} · ${b.campus}`" is-link @click="$router.push(`/books/${b.id}`)">
          <template #value><span class="price">{{ formatPrice(b.price) }}</span></template>
        </van-cell>
      </van-cell-group>
    </div>

    <div class="section">
      <div class="section-title">
        <span>最新上架</span>
        <span class="more" @click="$router.push('/books')">更多 <van-icon name="arrow" /></span>
      </div>
      <div v-if="latest.length" class="book-list">
        <BookCard v-for="b in latest" :key="b.id" :book="b" />
      </div>
      <EmptyState v-else description="暂无在售书籍" />
    </div>

    <div class="quick-actions">
      <van-button round block type="primary" @click="$router.push('/publish-book')">发布闲置书籍</van-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { listBooks, getRecommendations } from '@/api/book'
import type { Book } from '@/types'
import { formatPrice } from '@/utils/format'
import { SubjectCategoryOptions } from '@/constants/enums'
import { useAuthStore } from '@/stores/authStore'
import BookCard from '@/components/BookCard.vue'
import EmptyState from '@/components/EmptyState.vue'

const router = useRouter()
const auth = useAuthStore()
const isLoggedIn = computed(() => auth.isLoggedIn)
const keyword = ref('')
const latest = ref<Book[]>([])
const recommendations = ref<Book[]>([])

const banners = [
  { title: '教材循环，省钱环保', sub: '让闲置教材找到新主人', bg: 'linear-gradient(135deg,#1989fa,#5b8ff9)' },
  { title: '按课程名搜索', sub: '快速找到本学期教材', bg: 'linear-gradient(135deg,#07c160,#66bb6a)' },
  { title: '面交更安心', sub: '同校面对面交易', bg: 'linear-gradient(135deg,#ff976a,#ff6b6b)' },
]

const categories = SubjectCategoryOptions.map((c, i) => ({
  ...c,
  icon: ['bulb-o', 'notes-o', 'chart-trending-o', 'flower-o', 'more-o'][i],
  bg: ['#1989fa', '#07c160', '#ff976a', '#7232dd', '#969799'][i],
}))

onMounted(async () => {
  loadLatest()
  if (isLoggedIn.value) {
    await loadRecommendations()
  }
})

// 登录态变化后刷新同院系推荐（keep-alive 缓存页面不会自动重挂载）
watch(isLoggedIn, async (v) => {
  if (v) {
    await loadRecommendations()
  }
})

async function loadRecommendations() {
  try {
    recommendations.value = await getRecommendations({ limit: 5 })
  } catch {
    recommendations.value = []
  }
}

async function loadLatest() {
  try {
    const data = await listBooks({ sort: 'newest', page: 1, page_size: 6 })
    latest.value = data.list
  } catch {
    latest.value = []
  }
}

function onSearch() {
  router.push({ name: 'book-list', query: { keyword: keyword.value } })
}

function goCategory(value: string) {
  router.push({ name: 'book-list', query: { subject_category: value } })
}
</script>

<style scoped>
.banner {
  margin: 4px 0 12px;
  border-radius: 8px;
  overflow: hidden;
}
.banner-item {
  height: 120px;
  display: flex;
  align-items: center;
  padding: 0 20px;
}
.banner-text {
  color: #fff;
}
.banner-title {
  font-size: 20px;
  font-weight: 700;
}
.banner-sub {
  margin-top: 6px;
  font-size: 13px;
  opacity: 0.9;
}
.section {
  margin-bottom: 16px;
}
.section-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
  padding: 4px 4px 10px;
}
.more {
  font-size: 13px;
  color: #969799;
  font-weight: 400;
}
.category-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 10px;
  background: #fff;
  border-radius: 8px;
  padding: 14px 8px;
}
.category-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}
.category-icon {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.category-label {
  font-size: 12px;
  color: #646566;
}
.book-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.price {
  color: #ee0a24;
  font-weight: 600;
}
.quick-actions {
  padding: 8px 0 20px;
}
</style>
