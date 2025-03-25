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
	emits: ['login-success', 'to-home'],
	methods:{

		handleHomeView() {
			this.showUserInfo = true;
		},

		handleLoginSuccess(){
			this.isLogged = true;
			this.showUserInfo = true;
			this.username = sessionStorage.username;
			this.photo = sessionStorage.photo;
		},

		logout(){
			sessionStorage.clear();
			this.isLogged = false;
			this.$router.push("/");
		},

		goToUserInfo(){
			this.showUserInfo = false;
			this.$router.push("/user-info");
		}
	},
	mounted(){
	}
}
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">Example App</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
	</header>

	<div class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>General</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink to="/" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
								Home
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/link1" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
								Menu item 1
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/link2" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
								Menu item 2
							</RouterLink>
						</li>
					</ul>

					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>Secondary menu</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink :to="'/some/' + 'variable_here' + '/path'" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#file-text"/></svg>
								Item 1
							</RouterLink>
						</li>
					</ul>
				</div>
			</nav>

			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView @login-success="handleLoginSuccess" @to-home="handleHomeView"/>
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
}

.user-info button {
	padding: 5px 10px;
	border: none;
	border-radius: 5px;
	background-color: #007bff;
	color: white;
}
</style>
