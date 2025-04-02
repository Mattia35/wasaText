<script>
import setUsername from '../components/setUsername.vue'
export default {
    data: function() {
        return {
            username: sessionStorage.username,
            photo: sessionStorage.photo,
            errorMsg: "",
            showSetUsernameModal: false,
            newPhoto: null,
            showBigPhoto: false,
        }
    },
    emits: ['to-home', 
            'login-success', 
            'update-username', 
            'close', 
            'update-photo', 
            'update-groupname', 
            'update-group-photo',
            'update-group-members',
            'update-group-info'],
    components: {
        setUsername
    },
    methods: {
      handleFileChange(event) {
        const file = event.target.files[0]; 
        if (!file) {
          this.errorMsg = "Nessun file selezionato";
          return;
        }
        if (file.type !== "image/jpeg" && file.type !== "image/jpg") {
          this.errorMsg = "File type not supported, only jpg and jpeg are allowed";
          return;
        }
        if (file.size > 5242880) {
          this.errorMsg = "File size is too big. Max size is 5MB";
          return;
        }
        this.newPhoto = file;
        
        this.setNewUserPhoto();
        },

      async setNewUserPhoto() {
        try {
          const formData = new FormData();
          formData.append('image', this.newPhoto);

          // Effettua la richiesta di login al server con la foto inserita
          let response = await this.$axios.put(`/users/${sessionStorage.userID}/photo`, 
          formData, {headers: {Authorization: `${sessionStorage.token}`}});
          
          // Salva i dati dell'utente nella sessionStorage
          sessionStorage.photo = response.data.photo;
          this.photo = response.data.photo;
          
          // Emette l'evento di login avvenuto con successo
          this.$emit('update-photo');
          
        } catch (e) {
            this.errorMsg = e.toString();
        };
      },

      usernameModify() {
          //modale per la modifica dell'username
          this.showSetUsernameModal = !this.showSetUsernameModal;
      },

      goToHome() {
        this.$router.push("/home");
        this.$emit('to-home');
      },

      handleUpdateUsername() {
        this.username = sessionStorage.username;
        this.$emit('update-username');
      },

      fileInput(){
        this.$refs.file.click();
      },

      showPhoto(){
        this.showBigPhoto = !this.showBigPhoto;
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
            <img :src="`data:image/jpg;base64,${photo}`" alt="User photo" @click="showPhoto">
            <p class="username">{{username}}</p>
            <div class="user-info-box-buttons">
                <button @click="usernameModify">Modify your username</button>
                <input type="file" ref="file" accept=".jpg,.jpeg" @change="handleFileChange" style="display: none;"/>
                <button @click="fileInput">Modify your photo</button>
            </div>
            <div class="return-to-home-box">
                <button @click="goToHome">Return to home</button>
            </div>
        </div>
    </div>
    <div v-if="showBigPhoto" @click="showPhoto" class="modal-mask" >
      <img :src="`data:image/jpg;base64,${photo}`" alt="User photo" class="BigImage">
    </div>
    <setUsername :show="showSetUsernameModal" @close="usernameModify" @update-username="handleUpdateUsername" />
  </template>

  <style>
  .user-info-container {
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    height: 90vh;
  }

  .modal-mask {
    position: fixed;
    z-index: 9998;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center; 
}

.BigImage {
    width: 70vh; 
    height: 70vh;
    max-width: 80%;
    max-height: 80vh;
    border-radius: 50%; 
    object-fit: cover; 
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
    border-radius: 5px;
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
    border: 5px solid #007bff;
    border-radius: 50%;

  }

  .user-info-box-buttons {
    display: flex;
    gap: 10px;
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