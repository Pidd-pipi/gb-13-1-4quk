<template>
  <div class="chat-page">
    <van-nav-bar :title="title" left-arrow fixed placeholder @click-left="$router.back()" />
    <div class="chat-body" ref="bodyRef">
      <div v-for="m in messages" :key="m.id" :class="['msg-row', m.sender_id === auth.user?.id ? 'mine' : 'theirs']">
        <div class="bubble">
          <span v-if="m.image_url"><img :src="m.image_url" class="msg-img" alt="img" /></span>
          <span v-if="m.content">{{ m.content }}</span>
        </div>
      </div>
      <EmptyState v-if="!messages.length" description="开始你们的交易沟通吧" />
    </div>
    <div class="chat-input">
      <van-field v-model="draft" placeholder="输入消息..." @keyup.enter="send" />
      <van-button type="primary" size="small" @click="send">发送</van-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useConversationStore } from '@/stores/conversationStore'
import { useAuthStore } from '@/stores/authStore'
import EmptyState from '@/components/EmptyState.vue'

const route = useRoute()
const convStore = useConversationStore()
const auth = useAuthStore()
const draft = ref('')
const bodyRef = ref<HTMLElement>()
const messages = ref(convStore.messages)

const title = computed(() => {
  const c = convStore.conversations.find((x) => x.id === Number(route.params.id))
  if (!c) return '会话'
  const me = auth.user?.id
  return c.buyer_id === me ? c.seller?.name || '卖家' : c.buyer?.name || '买家'
})

onMounted(async () => {
  const id = Number(route.params.id)
  messages.value = await convStore.openConversation(id)
  scrollBottom()
})

watch(messages, () => scrollBottom(), { deep: true })

function scrollBottom() {
  nextTick(() => {
    if (bodyRef.value) {
      bodyRef.value.scrollTop = bodyRef.value.scrollHeight
    }
  })
}

async function send() {
  const content = draft.value.trim()
  if (!content) return
  draft.value = ''
  await convStore.send(content)
}
</script>

<style scoped>
.chat-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
}
.chat-body {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
  background: #f7f8fa;
}
.msg-row {
  display: flex;
  margin-bottom: 10px;
}
.msg-row.mine {
  justify-content: flex-end;
}
.bubble {
  max-width: 75%;
  background: #fff;
  border-radius: 8px;
  padding: 8px 10px;
  font-size: 14px;
  line-height: 1.5;
  word-break: break-all;
}
.mine .bubble {
  background: #1989fa;
  color: #fff;
}
.msg-img {
  max-width: 180px;
  border-radius: 6px;
  display: block;
}
.chat-input {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  background: #fff;
  border-top: 1px solid #ebedf0;
}
.chat-input .van-field {
  flex: 1;
  background: #f7f8fa;
  border-radius: 18px;
}
</style>
