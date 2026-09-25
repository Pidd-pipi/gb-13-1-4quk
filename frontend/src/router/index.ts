import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'

const routes: RouteRecordRaw[] = [
  { path: '/', name: 'home', component: () => import('@/pages/Home.vue') },
  { path: '/books', name: 'book-list', component: () => import('@/pages/BookList.vue') },
  { path: '/books/:id', name: 'book-detail', component: () => import('@/pages/BookDetail.vue') },
  { path: '/publish-book', name: 'publish-book', component: () => import('@/pages/PublishBook.vue'), meta: { requiresAuth: true } },
  { path: '/wishes', name: 'wish-list', component: () => import('@/pages/WishList.vue') },
  { path: '/publish-wish', name: 'publish-wish', component: () => import('@/pages/PublishWish.vue'), meta: { requiresAuth: true } },
  { path: '/messages', name: 'messages', component: () => import('@/pages/Messages.vue'), meta: { requiresAuth: true } },
  { path: '/chat/:id', name: 'chat', component: () => import('@/pages/Chat.vue'), meta: { requiresAuth: true } },
  { path: '/my-books', name: 'my-books', component: () => import('@/pages/MyBooks.vue'), meta: { requiresAuth: true } },
  { path: '/borrows', name: 'borrows', component: () => import('@/pages/Borrows.vue'), meta: { requiresAuth: true } },
  { path: '/favorites', name: 'favorites', component: () => import('@/pages/Favorites.vue'), meta: { requiresAuth: true } },
  { path: '/profile', name: 'profile', component: () => import('@/pages/Profile.vue'), meta: { requiresAuth: true } },
  { path: '/login', name: 'login', component: () => import('@/pages/Login.vue') },
  { path: '/register', name: 'register', component: () => import('@/pages/Register.vue') },
  { path: '/audit', name: 'audit', component: () => import('@/pages/Audit.vue'), meta: { requiresAuth: true, requiresAdmin: true } },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.isLoggedIn) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.meta.requiresAdmin && auth.user?.role !== 'admin') {
    return { name: 'home' }
  }
  return true
})

export default router
