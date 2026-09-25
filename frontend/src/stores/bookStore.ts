import { ref } from 'vue'
import { defineStore } from 'pinia'
import * as bookApi from '@/api/book'
import type { BookQuery } from '@/api/book'
import type { Book } from '@/types'
import type { PageData } from '@/types'

// bookStore：书籍列表、推荐、收藏、浏览历史的共享状态
export const useBookStore = defineStore('book', () => {
  const books = ref<Book[]>([])
  const total = ref(0)
  const loading = ref(false)
  const recommendations = ref<Book[]>([])
  const favorites = ref<Book[]>([])
  const history = ref<Book[]>([])

  async function loadBooks(query: BookQuery) {
    loading.value = true
    try {
      const data: PageData<Book> = await bookApi.listBooks(query)
      books.value = data.list
      total.value = data.total
      return data
    } finally {
      loading.value = false
    }
  }

  async function loadRecommendations(limit = 10) {
    recommendations.value = await bookApi.getRecommendations({ limit })
    return recommendations.value
  }

  async function loadFavorites() {
    const data = await bookApi.listFavorites({ page: 1, page_size: 50 })
    favorites.value = data.list
    return data
  }

  async function loadHistory(limit = 20) {
    const data = await bookApi.listHistory({ limit })
    history.value = data.list
    return data
  }

  async function toggleFavorite(book: Book) {
    if (book.is_favorite) {
      await bookApi.removeFavorite(book.id)
      book.is_favorite = false
      book.favorite_count = Math.max(0, (book.favorite_count || 1) - 1)
    } else {
      await bookApi.addFavorite(book.id)
      book.is_favorite = true
      book.favorite_count = (book.favorite_count || 0) + 1
    }
  }

  return { books, total, loading, recommendations, favorites, history, loadBooks, loadRecommendations, loadFavorites, loadHistory, toggleFavorite }
})
