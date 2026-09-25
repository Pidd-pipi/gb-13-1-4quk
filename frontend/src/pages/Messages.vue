<template>
  <div class="page">
    <van-nav-bar title="消息" fixed placeholder />
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-cell-group inset v-if="list.length">
        <van-cell v-for="c in list" :key="c.id" is-link @click="$router.push(`/chat/${c.id}`)">
          <template #icon>
            <van-badge :content="c.unread_count > 99 ? '99+' : c.unread_count || ''" :show-zero="false">
              <div class="avatar">{{ counterpart(c).slice(0, 1) }}</div>
            </van-badge>
          </template>
          <template #title>
            <div class="conv-title">
              <span>{{ counterpart(c) }}</span>
              <span v-if="c.book" class="conv-book">{{ c.book.title }}</span>
            </div>
            <div class="conv-preview van-ellipsis">{{ c.last_message || '暂无消息' }}</div>
          </template>
          <template #value>
            <span class="conv-time">{{ c.last_message_at || c.created_at }}</span>
          </template>
        </van-cell>
      </van-cell-group>
      <EmptyState v-else description="暂无会话，去书籍详情页联系卖家吧" icon="chat-o" />
    </van-pull-refresh>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useConversationStore } from '@/stores/conversationStore'
import { useAuthStore } from '@/stores/authStore'
import EmptyState from '@/components/EmptyState.vue'

const convStore = useConversationStore()
const auth = useAuthStore()
const list = ref(convStore.conversations)
const refreshing = ref(false)

function counterpart(c: { buyer_id: number; buyer?: { name?: string }; seller?: { name?: string } }) {
  if (c.buyer_id === auth.user?.id) {
    return c.seller?.name || '卖家'
  }
  return c.buyer?.name || '买家'
}

onMounted(load)

async function load() {
  try {
    const data = await convStore.loadConversations()
    list.value = data.list
  } catch {
    /* toast */
  } finally {
    refreshing.value = false
  }
}

async function onRefresh() {
  await load()
}
</script>

<style scoped>
.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #1989fa;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  margin-right: 10px;
}
.conv-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
}
.conv-book {
  font-size: 11px;
  color: #1989fa;
  background: #ecf5ff;
  border-radius: 4px;
  padding: 0 4px;
}
.conv-preview {
  color: #969799;
  font-size: 12px;
  margin-top: 2px;
}
.conv-time {
  color: #c8c9cc;
  font-size: 11px;
}
</style>
