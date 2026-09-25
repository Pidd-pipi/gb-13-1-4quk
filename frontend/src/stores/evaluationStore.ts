import { ref } from 'vue'
import { defineStore } from 'pinia'
import * as evalApi from '@/api/evaluation'
import type { Evaluation } from '@/types'

export const useEvaluationStore = defineStore('evaluation', () => {
  const evaluations = ref<Evaluation[]>([])
  const total = ref(0)

  async function loadUserEvaluations(userId: number) {
    const data = await evalApi.listUserEvaluations(userId, { page: 1, page_size: 50 })
    evaluations.value = data.list
    total.value = data.total
    return data
  }

  return { evaluations, total, loadUserEvaluations }
})
