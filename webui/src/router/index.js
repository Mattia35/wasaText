import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import UserInfoView from '../views/UserInfoView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/home', component: HomeView},
		{path: '/user-info', component: UserInfoView},
		{path: '/conversation', component: HomeView},
		{path: '/some/:id/link', component: HomeView},
	]
})

export default router
