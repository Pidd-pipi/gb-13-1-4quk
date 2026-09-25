import { ref } from 'vue'
import { defineStore } from 'pinia'
import * as convApi from '@/api/conversation'
import type { Conversation, Message } from '@/types'

export const useConversationStore = defineStore('conversation', () => {
  const conversations = ref<Conversation[]>([])
  const messages = ref<Message[]>([])
  const currentId = ref<number>(0)
  const loading = ref(false)

  async function loadConversations() {
    loading.value = true
    try {
      const data = await convApi.listConversations({ page: 1, page_size: 50 })
      conversations.value = data.list
      return data
    } finally {
      loading.value = false
    }
  }

  async function openConversation(id: number) {
    currentId.value = id
    messages.value = await convApi.listMessages(id)
    await convApi.markRead(id)
    return messages.value
  }

  async function send(content: string, imageUrl?: string) {
    const msg = await convApi.sendMessage(currentId.value, content, imageUrl)
    messages.value.push(msg)
    await loadConversations()
    return msg
  }

  return { conversations, messages, currentId, loading, loadConversations, openConversation, send }
})
