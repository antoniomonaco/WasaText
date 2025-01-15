import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import HomeView from '../views/HomeView.vue'
import ConversationsView from '../views/ConversationsView.vue'
import ChatView from '../views/ChatView.vue'
import MainView from '../views/MainView.vue'


const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/',name: "Main" ,component: MainView},
		{path: '/session', name: 'login', component: LoginView},
		{path: '/conversations/', name: 'conversations', component: ConversationsView},
		{path: '/conversations/:conversationID', name: 'chat', component: ChatView},
	]
})
router.beforeEach((to, from, next) => {
	if (to.matched.some(record => record.meta.requiresAuth)) {
	  // Verifica se l'utente è autenticato
	  const token = localStorage.getItem('authToken');
	  if (token) {
		next();
	  } else {
		next({ name: 'Login' }); // Reindirizza alla pagina di login se non autenticato
	  }
	} else {
	  next();
	}
  });
export default router
