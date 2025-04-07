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

      // function to handle the file change (when the user selects a file)
      handleFileChange(event) {
        // get the file from the input
        const file = event.target.files[0]; 
        // check if the file is null (if the user has not selected a file)
        if (!file) {
        this.errorMsg = "Nessun file selezionato";
        return;
        }
        // check if the file is an image (jpeg) and if it is too big
        if (file.type !== "image/jpeg" && file.type !== "image/jpg") {
        this.errorMsg = "File type not supported, only jpg and jpeg are allowed";
        return;
        }
        // check if the file is too big (5MB)
        if (file.size > 5242880) {
        this.errorMsg = "File size is too big. Max size is 5MB";
        return;
        }
        // assign the file to the newPhoto variable
        this.newPhoto = file;
        // go to the function to set the new group photo
        this.setNewGroupPhoto();
      },

      // function to set the new group photo
      async setNewUserPhoto() {
        // try the request to set the new user photo
        try {
          // create a FormData object to send the file
          const formData = new FormData();
          formData.append('image', this.newPhoto);
          // make the request to set the new user photo
          let response = await this.$axios.put(`/users/${sessionStorage.userID}/photo`, 
          formData, {headers: {Authorization: `${sessionStorage.token}`}});
          // save the new user photo in sessionStorage and in this vue
          sessionStorage.photo = response.data.photo;
          this.photo = response.data.photo;
          // send the event to update the user photo in the parent component
          this.$emit('update-photo');
        } catch (e) {
          // save and print the error message
            this.errorMsg = e.toString();
        };
      },

      // function to open or close the modal to set the username
      usernameModify() {
          this.showSetUsernameModal = !this.showSetUsernameModal;
      },

      // function to go to the home page and emit the event to the parent component, who set the showUserInfo to false
      goToHome() {
        this.$router.push("/home");
        this.$emit('to-home');
      },

      // function to update the username in the parent component and in this vue
      handleUpdateUsername() {
        this.username = sessionStorage.username;
        this.$emit('update-username');
      },

      // function to open the file selector when the user clicks on the button
      fileInput(){
        this.$refs.file.click();
      },

      // function to show the photo in a bigger size
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
    <!-- container of the user info -->
    <div class="user-info-container">
        <h1>User Info</h1>
        <div class="user-info-box">
            <!-- user photo -->
            <img :src="`data:image/jpg;base64,${photo}`" alt="User photo" @click="showPhoto">
            <!-- user name -->
            <p class="username">{{username}}</p>
            <!-- buttons to modify the username and the photo -->
            <div class="user-info-box-buttons">
                <button @click="usernameModify">Modify your username</button>
                <input type="file" ref="file" accept=".jpg,.jpeg" @change="handleFileChange" style="display: none;"/>
                <button @click="fileInput">Modify your photo</button>
            </div>
            <!-- button to go to the home page -->
            <div class="return-to-home-box">
                <button @click="goToHome">Return to home</button>
            </div>
        </div>
    </div>

    <!-- modal to show the photo in a bigger size -->
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