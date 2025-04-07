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

		// function to go do the logout 
		logOut(){
			// remove all data from sessionStorage
			sessionStorage.clear();
			// set the variable that controls the login status to false
			this.isLogged = false;
			// set the variable that controls if the user info icon is shown to false
			this.showUserInfo = false;
			// set the username and photo to empty
			this.username = "";
			this.photo = "";
			// go to the login page
			this.$router.push("/");
		},

		// function to show the user info
		handleHomeView() {
			this.showUserInfo = true;
		},

		// function to handle the login success
		handleLoginSuccess(){
			// set the variable that controls the login status to true
			this.isLogged = true;
			// set the variable that controls if the user info icon is shown to true
			this.showUserInfo = true;
			// set the username and photo to the values stored in sessionStorage
			this.username = sessionStorage.username;
			this.photo = sessionStorage.photo;
		},

		// function to go to the user info page
		goToUserInfo(){
			// set the variable that controls if the user info icon is shown to false
			this.showUserInfo = false;
			// go to the user info page
			this.$router.push("/user-info");
		},
		
		// function to handle the update of the username
		handleUpdateUsername(){
			// set the username to the value stored in sessionStorage
			this.username = sessionStorage.username;
		},

		// function to handle the update of the user photo
		handleUpdateUserPhoto(){
			// set the photo to the value stored in sessionStorage
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

	<!-- Main container -->
	<div class="container-fluid">
		<div class="row">
			<!-- Sidebar, that is showed when the user is logged -->
			<nav v-if="isLogged" id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<!-- menu where user can do logout -->
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
		
		<!-- User info icon that is showed when the user is logged. When the user clicks on it, the user info page is opened -->
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
