import { computed, reactive, ref } from 'vue'
import type { PageData } from '@/types'

// usePagination：跨页面复用的分页加载 hook
export function usePagination<T>(fetcher: (page: number, pageSize: number) => Promise<PageData<T>>, pageSize = 10) {
  const list = reactive([] as T[])
  const total = ref(0)
  const page = ref(1)
  const loading = ref(false)
  const finished = ref(false)
  const refreshing = ref(false)

  const hasMore = computed(() => list.length < total.value)

  async function load() {
    if (loading.value) return
    loading.value = true
    try {
      const data = await fetcher(page.value, pageSize)
      if (page.value === 1) {
        list.splice(0, list.length, ...(data.list as never[]))
      } else {
        list.push(...(data.list as never[]))
      }
      total.value = data.total
      finished.value = list.length >= total.value
    } catch (e) {
      finished.value = true
      console.error(e)
    } finally {
      loading.value = false
      refreshing.value = false
    }
  }

  function onLoad() {
    if (finished.value) return
    page.value += 1
    load()
  }

  function onRefresh() {
    page.value = 1
    finished.value = false
    load()
  }

  function reset() {
    page.value = 1
    list.splice(0, list.length)
    total.value = 0
    finished.value = false
    load()
  }

  return { list, total, loading, finished, refreshing, hasMore, onLoad, onRefresh, reset, load }
}
