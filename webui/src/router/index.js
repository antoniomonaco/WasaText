import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import HomeView from '../views/HomeView.vue'
import ConversationsView from '../views/ConversationsView.vue'
import ChatView from '../views/ChatView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/',component: HomeView},
		{path: '/session', name: 'login', component: LoginView},
		{path: '/conversations/', name: 'conversations', component: ConversationsView},
		{path: '/conversations/:conversationID', name: 'chat', component: ChatView},
	]
})

export default router
