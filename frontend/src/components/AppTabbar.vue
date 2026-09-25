<template>
  <van-tabbar v-model="active" fixed placeholder route>
    <van-tabbar-item replace to="/" icon="wap-home-o">首页</van-tabbar-item>
    <van-tabbar-item replace to="/books" icon="search">找书</van-tabbar-item>
    <van-tabbar-item replace to="/publish-book" icon="add-o">发布</van-tabbar-item>
    <van-tabbar-item replace to="/wishes" icon="todo-list-o">求购</van-tabbar-item>
    <van-tabbar-item replace to="/profile" icon="user-o">我的</van-tabbar-item>
  </van-tabbar>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const active = ref(0)

const map: Record<string, number> = {
  home: 0,
  'book-list': 1,
  'publish-book': 2,
  'wish-list': 3,
  profile: 4,
}

watch(
  () => route.name,
  (name) => {
    if (name && typeof name === 'string' && map[name] !== undefined) {
      active.value = map[name]
    }
  },
  { immediate: true },
)
</script>
