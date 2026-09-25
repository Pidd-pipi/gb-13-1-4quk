<template>
  <div class="page book-list-page">
    <van-nav-bar title="找书" fixed placeholder />
    <van-search v-model="query.keyword" placeholder="书名 / 作者 / ISBN" show-action @search="onSearch">
      <template #action><span @click="onSearch">搜索</span></template>
    </van-search>

    <div class="filter-bar">
      <van-dropdown-menu>
        <van-dropdown-item v-model="query.subject_category" :options="subjectOptions" @change="onFilter" />
        <van-dropdown-item v-model="query.condition" :options="conditionOptions" @change="onFilter" />
        <van-dropdown-item v-model="query.sort" :options="sortOptions" @change="onFilter" />
      </van-dropdown-menu>
    </div>

    <div class="price-filter">
      <van-field v-model="minPrice" type="number" placeholder="最低价" />
      <span class="dash">—</span>
      <van-field v-model="maxPrice" type="number" placeholder="最高价" />
      <van-button size="small" type="primary" @click="onFilter">确定</van-button>
    </div>

    <div class="borrow-filter">
      <van-checkbox v-model="borrowableOnly" shape="square" icon-size="16px" @change="onFilter">只看可短借</van-checkbox>
    </div>

    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list v-model:loading="loading" :finished="finished" finished-text="没有更多了" @load="onLoad">
        <div v-if="list.length" class="book-list">
          <BookCard v-for="b in list" :key="b.id" :book="b" />
        </div>
        <EmptyState v-else-if="finished" description="未找到相关书籍" />
      </van-list>
    </van-pull-refresh>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { listBooks } from '@/api/book'
import type { Book } from '@/types'
import { ConditionOptions, SortOptions, SubjectCategoryOptions } from '@/constants/enums'
import BookCard from '@/components/BookCard.vue'
import EmptyState from '@/components/EmptyState.vue'

const route = useRoute()
const query = reactive({
  keyword: (route.query.keyword as string) || '',
  subject_category: (route.query.subject_category as string) || '',
  condition: '',
  sort: 'newest',
})
const minPrice = ref('')
const maxPrice = ref('')
const borrowableOnly = ref(false)
const list = ref<Book[]>([])
const page = ref(1)
const total = ref(0)
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)

const subjectOptions = [{ text: '全部分类', value: '' }, ...SubjectCategoryOptions.map((o) => ({ text: o.label, value: o.value }))]
const conditionOptions = [{ text: '新旧不限', value: '' }, ...ConditionOptions.map((o) => ({ text: o.label, value: o.value }))]
const sortOptions = SortOptions.map((o) => ({ text: o.label, value: o.value }))

const hasMore = computed(() => list.value.length < total.value)

async function fetchPage() {
  loading.value = true
  try {
    const data = await listBooks({
      keyword: query.keyword,
      subject_category: query.subject_category,
      condition: query.condition,
      sort: query.sort,
      min_price: minPrice.value ? Number(minPrice.value) : undefined,
      max_price: maxPrice.value ? Number(maxPrice.value) : undefined,
      borrowable: borrowableOnly.value ? true : undefined,
      page: page.value,
      page_size: 10,
    })
    total.value = data.total
    if (page.value === 1) {
      list.value = data.list
    } else {
      list.value = list.value.concat(data.list)
    }
    finished.value = list.value.length >= total.value
  } catch {
    finished.value = true
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

function onLoad() {
  if (!hasMore.value) return
  page.value += 1
  fetchPage()
}

function onRefresh() {
  page.value = 1
  finished.value = false
  fetchPage()
}

function onFilter() {
  page.value = 1
  finished.value = false
  list.value = []
  fetchPage()
}

function onSearch() {
  onFilter()
}

watch(
  () => route.query,
  (q) => {
    query.keyword = (q.keyword as string) || ''
    query.subject_category = (q.subject_category as string) || ''
    onFilter()
  },
)

fetchPage()
</script>

<style scoped>
.filter-bar {
  margin: 4px 0;
}
.book-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 0 12px;
}
.price-filter {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #fff;
}
.price-filter .van-field {
  flex: 1;
  padding: 6px 8px;
  background: #f7f8fa;
  border-radius: 6px;
}
.dash {
  color: #969799;
}
.borrow-filter {
  padding: 0 12px 8px;
  background: #fff;
}
</style>
