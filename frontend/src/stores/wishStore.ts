import { ref } from 'vue'
import { defineStore } from 'pinia'
import * as wishApi from '@/api/wish'
import type { WishQuery } from '@/api/wish'
import type { Wish } from '@/types'
import type { PageData } from '@/types'

export const useWishStore = defineStore('wish', () => {
  const wishes = ref<Wish[]>([])
  const total = ref(0)
  const loading = ref(false)

  async function loadWishes(query: WishQuery) {
    loading.value = true
    try {
      const data: PageData<Wish> = await wishApi.listWishes(query)
      wishes.value = data.list
      total.value = data.total
      return data
    } finally {
      loading.value = false
    }
  }

  return { wishes, total, loading, loadWishes }
})
