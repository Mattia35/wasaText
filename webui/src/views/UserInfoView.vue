<script>
import setUsername from '../components/setUsername.vue'
export default {
    data: function() {
        return {
            username: sessionStorage.username,
            photo: sessionStorage.photo,
            errorMsg: "",
            showSetUsernameModal: false,
        }
    },
    emits: ['to-home', 'login-success'],
    components: {
        setUsername
    },
    methods: {
        usernameModify() {
            //modale per la modifica dell'username
            this.showSetUsernameModal = !this.showSetUsernameModal
        },

        photoModify() {
            //modale per la modifica della foto
        },

        goToHome() {
            this.$router.push("/home");
            this.$emit('to-home');
        }
        
    },
    mounted() {
        
    }
}
</script>

<template>
    <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
    <div class="user-info-container">
        <h1>User Info</h1>
        <div class="user-info-box">
            <img :src="`data:image/jpg;base64,${photo}`" alt="User photo" >
            <p class="username">{{username}}</p>
            <div class="user-info-box-buttons">
                <button @click="usernameModify">Modify your username</button>
                <button @click="photoModify">Modify your photo</button>
            </div>
            <div class="return-to-home-box">
                <button @click="goToHome">Return to home</button>
            </div>
        </div>
    </div>

    <setUsername :show="showSetUsernameModal" @close="usernameModify"></setUsername>
  </template>

  <style>
  .user-info-container {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    height: 50vh;
  }

  .user-info-box-buttons button {
    width: 250px;       
    cursor: pointer;
    height: 50px;    
    font-size: 16px;    
    padding: 10px;    
    text-align: center; 
    border: none;   
    background-color: #007bff; 
    color: white;   
  }
  .username {
    font-size: 20px;
  }
  .user-info-box {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
  }
  .user-info-box img {
    width: 150px;
    height: 150px;
    border-radius: 50%;
  }
  .user-info-box-buttons {
    display: flex;
    gap: 10px;
  }
  .user-info-box-buttons button {
    padding: 10px 20px;
    border: none;
    border-radius: 5px;
    background-color: #007bff;
    color: white;
    cursor: pointer;
  }

  .return-to-home-box button {
    padding: 10px 20px;
    border: none;
    border-radius: 5px;
    background-color: #ff0000; 
    color: white;
    cursor: pointer;
  }

  </style>