<template>
  <div class="page">
    <van-nav-bar title="我的" fixed placeholder />

    <div class="profile-header" @click="onHeaderClick">
      <van-image round width="64" height="64" :src="user?.avatar_url || fallback" />
      <div class="profile-info">
        <div class="name">{{ user?.name || '未登录' }}</div>
        <div class="sub">{{ user ? `${user.department} · ${user.campus}` : '登录后体验完整功能' }}</div>
        <div v-if="user" class="sub">{{ user.student_no }} · {{ user.email }}</div>
      </div>
    </div>

    <template v-if="user">
      <van-cell-group inset class="stats-card">
        <van-cell title="好评率" :value="stats?.good_rate || '--'" />
        <van-cell title="累计评价" :value="`${stats?.total_evaluations || 0} 条（好 ${stats?.good_count || 0} / 中 ${stats?.neutral_count || 0} / 差 ${stats?.bad_count || 0}）`" />
        <van-cell title="在售 / 已售书籍" :value="`${stats?.on_sale_books || 0} / ${stats?.sold_books || 0}`" />
        <van-cell v-if="stats?.risk_flagged" title="风险提示" value="差评率较高，交易请谨慎" />
      </van-cell-group>

      <van-cell-group inset class="menu">
        <van-cell title="我发布的书籍" icon="records-o" is-link @click="$router.push('/my-books')" />
        <van-cell title="短借管理" icon="exchange" is-link @click="$router.push('/borrows')" />
        <van-cell title="我的收藏" icon="star-o" is-link @click="$router.push('/favorites')" />
        <van-cell title="我的消息" icon="chat-o" is-link @click="$router.push('/messages')" />
        <van-cell title="发布求购" icon="todo-list-o" is-link @click="$router.push('/publish-wish')" />
        <van-cell v-if="isAdmin" title="审计日志" icon="shield-o" is-link @click="$router.push('/audit')" />
      </van-cell-group>

      <div class="section-title">收到的评价</div>
      <van-cell-group inset>
        <van-cell v-for="e in evaluations" :key="e.id">
          <template #title>
            <span class="eval-type" :class="e.type">{{ e.type_text }}</span>
            <span class="eval-from">{{ e.from_user?.name || '用户' }} · {{ e.book?.title || '' }}</span>
          </template>
          <template #label>{{ e.content || '（无文字评价）' }}</template>
          <template #value><span class="eval-time">{{ e.created_at }}</span></template>
        </van-cell>
        <van-cell v-if="!evaluations.length" title="暂无评价" />
      </van-cell-group>

      <van-cell-group inset class="menu">
        <van-cell title="完善资料" icon="edit" is-link @click="showEdit = true" />
        <van-cell title="退出登录" icon="power" is-link @click="onLogout" />
      </van-cell-group>
    </template>
    <template v-else>
      <div class="login-actions">
        <van-button round block type="primary" @click="$router.push('/login')">登录 / 注册</van-button>
      </div>
    </template>

    <van-popup v-model:show="showEdit" round position="bottom" style="height: 60%">
      <div class="edit-panel">
        <div class="edit-title">完善个人信息</div>
        <van-cell-group inset>
          <van-field v-model="editForm.name" label="姓名" placeholder="姓名" />
          <van-field v-model="editForm.department" label="院系" placeholder="院系" />
          <van-field v-model="editForm.campus" label="校区" placeholder="校区" />
          <van-field v-model="editForm.contact" label="联系方式" placeholder="手机 / 微信" />
        </van-cell-group>
        <div class="edit-actions">
          <van-button round block type="primary" @click="saveProfile">保存</van-button>
        </div>
      </div>
    </van-popup>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showSuccessToast } from 'vant'
import { useAuth } from '@/hooks/useAuth'
import { getStats, updateProfile } from '@/api/user'
import { useEvaluationStore } from '@/stores/evaluationStore'
import { useAuthStore } from '@/stores/authStore'
import type { UserStats } from '@/types'

const route = useRoute()
const router = useRouter()
const { user, isLoggedIn, isAdmin, logout } = useAuth()
const evalStore = useEvaluationStore()
const stats = ref<UserStats | null>(null)
const showEdit = ref(false)
const fallback = ''
const editForm = ref({ name: '', department: '', campus: '', contact: '' })

const evaluations = ref(evalStore.evaluations)

onMounted(async () => {
  const userId = route.query.userId ? Number(route.query.userId) : user.value?.id
  if (!userId) return
  try {
    stats.value = await getStats(userId)
    const data = await evalStore.loadUserEvaluations(userId)
    evaluations.value = data.list
  } catch {
    /* toast */
  }
})

function onHeaderClick() {
  if (!isLoggedIn.value) {
    router.push('/login')
  }
}

async function saveProfile() {
  if (!user.value) return
  try {
    await updateProfile(editForm.value)
    await useAuthStoreRefresh()
    showEdit.value = false
    showSuccessToast('已保存')
  } catch {
    /* toast */
  }
}

async function useAuthStoreRefresh() {
  await useAuthStore().fetchMe()
}

function onLogout() {
  logout()
  router.push('/login')
}
</script>

<style scoped>
.profile-header {
  display: flex;
  align-items: center;
  gap: 12px;
  background: linear-gradient(135deg, #1989fa, #5b8ff9);
  border-radius: 8px;
  padding: 20px 16px;
  color: #fff;
  margin-bottom: 12px;
}
.name {
  font-size: 18px;
  font-weight: 700;
}
.sub {
  font-size: 12px;
  opacity: 0.9;
  margin-top: 2px;
}
.stats-card {
  margin-bottom: 12px;
}
.menu {
  margin-bottom: 12px;
}
.section-title {
  font-size: 16px;
  font-weight: 600;
  padding: 4px 4px 10px;
}
.eval-type {
  font-weight: 600;
  margin-right: 6px;
}
.eval-type.good {
  color: #07c160;
}
.eval-type.neutral {
  color: #ff976a;
}
.eval-type.bad {
  color: #ee0a24;
}
.eval-from {
  color: #646566;
  font-size: 13px;
}
.eval-time {
  color: #c8c9cc;
  font-size: 11px;
}
.login-actions {
  padding: 24px 16px;
}
.edit-panel {
  padding: 16px 0 24px;
}
.edit-title {
  text-align: center;
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 12px;
}
.edit-actions {
  padding: 16px;
}
</style>
