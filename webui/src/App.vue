<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {
	data: function(){
		return {
			showUserInfo : false,
			isLogged : sessionStorage.token ? true : false,
			username : "",
			photo : "",
		};
	},
	emits: ['login-success', 
			'to-home', 
			'update-username', 
			'close', 
			'update-photo', 
			'update-groupname', 
			'update-group-photo',
			'update-group-members',
			'update-group-info'],
	methods:{
		logOut(){
			sessionStorage.clear();
			this.isLogged = false;
			this.showUserInfo = false;
			this.username = "";
			this.photo = "";
			this.$router.push("/");
		},

		handleHomeView() {
			this.showUserInfo = true;
		},

		handleLoginSuccess(){
			this.isLogged = true;
			this.showUserInfo = true;
			this.username = sessionStorage.username;
			this.photo = sessionStorage.photo;
		},

		goToUserInfo(){
			this.showUserInfo = false;
			this.$router.push("/user-info");
		},

		handleUpdateUsername(){
			this.username = sessionStorage.username;
		},

		handleUpdateUserPhoto(){
			this.photo = sessionStorage.photo;
		}

	},
	mounted(){
	}
}
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
	</header>

	<div class="container-fluid">
		<div class="row">
			<nav v-if="isLogged" id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>General</span>
					</h6>
					<ul class="nav flex-column">
						<li @click="logOut" class="nav-item">
							<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#log-out"/></svg>
							Logout
						</li>
					</ul>
				</div>
			</nav>
			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView @login-success="handleLoginSuccess" @to-home="handleHomeView" @update-username="handleUpdateUsername" @update-photo="handleUpdateUserPhoto"/>
			</main>
		</div>

		<div v-if="isLogged && showUserInfo" class="user-info">
			<img :src="`data:image/jpg;base64,${photo}`" alt="User photo" >
			<span>{{ username }}</span>
			<button @click="goToUserInfo">User info</button>
		</div>
	</div>
</template>

<style>
.nav-item {
	display: flex;
	justify-content: flex-start;
}

.nav-item svg {
	margin-right: 10px;
}

.user-info {
	position: fixed;
	bottom: 10px;
	left: 10px;
	padding: 10px;
	background-color: #333;
	box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
	border-radius: 8px;
	flex-direction: column;
	color: white;
	display: flex;
	align-items: center;
	gap: 10px;
	z-index: 1000;
}

.user-info img {
	width: 80px;
	height: 80px;
	border-radius: 50%;
	object-fit: cover;
	border: 5px solid #007bff;
}

.user-info button {
	padding: 5px 10px;
	border: none;
	border-radius: 5px;
	background-color: #007bff;
	color: white;
}
</style>
