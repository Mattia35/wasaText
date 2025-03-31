import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import UserInfoView from '../views/UserInfoView.vue'
import ChatView from '../views/ChatView.vue'
import TemporaryChatView from '../views/TemporaryChatView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/home', component: HomeView},
		{path: '/user-info', component: UserInfoView},
		{path: '/chat/:convId', component: ChatView},
		{path: '/temporary-chat/:usernameDestination', component: TemporaryChatView},
	]
})

export default router
