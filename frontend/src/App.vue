<template>
  <div class="app-shell">
    <router-view v-slot="{ Component }">
      <keep-alive>
        <component :is="Component" v-if="Component" :key="$route.fullPath" />
      </keep-alive>
    </router-view>
    <AppTabbar v-if="showTabbar" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppTabbar from '@/components/AppTabbar.vue'

const route = useRoute()
const hiddenTabbarRoutes = ['login', 'register', 'chat', 'book-detail', 'publish-book', 'publish-wish', 'audit']
const showTabbar = computed(() => !hiddenTabbarRoutes.includes(String(route.name)))
</script>

<style>
html,
body {
  margin: 0;
  padding: 0;
  background: #f7f8fa;
  font-family: -apple-system, BlinkMacSystemFont, 'Helvetica Neue', Helvetica, Segoe UI, Arial, Roboto, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft Yahei', sans-serif;
  -webkit-font-smoothing: antialiased;
}
.app-shell {
  max-width: 640px;
  margin: 0 auto;
  min-height: 100vh;
  position: relative;
  background: #f7f8fa;
}
.page {
  padding: 12px;
  padding-bottom: 60px;
}
.page-title {
  font-size: 18px;
  font-weight: 600;
  margin: 4px 0 12px;
}
.van-nav-bar {
  background: #fff;
}
</style>
