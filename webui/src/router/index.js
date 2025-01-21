import { createRouter, createWebHashHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import HomeView from '../views/HomeView.vue'
import ConversationsView from '../views/ConversationsView.vue'
import ChatView from '../views/ChatView.vue'
import MainView from '../views/MainView.vue'
import ProfileView from '../views/ProfileView.vue'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/session'  // Reindirizza / a /session
    },
    {
      path: '/main',
      name: "Main",
      component: MainView,
      meta: { requiresAuth: true }
    },
    {
      path: '/session',
      name: 'Login',
      component: LoginView
    },
    {
      path: '/conversations/',
      name: 'conversations',
      component: ConversationsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/conversations/:conversationID',
      name: 'chat',
      component: ChatView,
      meta: { requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
      meta: { requiresAuth: true }
    },
  ]
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('authToken')
  
  // Se l'utente è autenticato e prova ad accedere alla pagina di login
  if (to.path === '/session' && token) {
    next('/main') // Reindirizza a main se già loggato
    return
  }
  
  // Se la rotta richiede autenticazione
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!token) {
      next('/session') // Reindirizza al login se non autenticato
    } else {
      next() // Procedi normalmente
    }
  } else {
    next() // Procedi per le rotte che non richiedono auth
  }
})

export default router