<script>
export default {
  data() {
    return {
      username: "",
      errorMsg: "",
      usernameValidation: new RegExp('^[a-z0-9]{1,15}$'),
    }
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
  methods: {
    // function to do the login
    async doLogin() {
      // try the request
      try {
        // control if the username is valid
        if (this.username.length < 1 || this.username.length > 15) throw "Invalid username, it must contains min 1 character and max 15 characters"
        if (!this.usernameValidation.test(this.username)) throw "Invalid username, it must contain only letters and numbers"
        // make the request to the server
        let response = await this.$axios.post('/session', {
          username: this.username,
        });
        // save the user data in sessionStorage
        sessionStorage.userID = response.data.user.userId;
        sessionStorage.username = response.data.user.username;
        sessionStorage.token = response.data.user.userId;
        sessionStorage.photo = response.data.user.userPhoto;
        // go to the home page
        this.$router.push("/home");
        // send the event to declare the login success
        this.$emit('login-success');
      } catch (e) {
        // save and print the error message
        this.errorMsg = e.toString();
      };
    }
  },
  mounted() {
    // check if the user is already logged in. If yes, go to the home page, otherwise clear the sessionStorage
    if (sessionStorage.token) {
      this.$router.push("/home");
      return;
    }
    sessionStorage.clear();
  },
}

</script>

<template>
  <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
  <!-- Login containter -->
  <div class="login-container">
    <!-- Form to login -->
    <form @submit.prevent="doLogin">
      <h1>WasaText</h1>
      <!-- Input for the username -->
      <input type="text" v-model="username" placeholder="Enter your username" required/>
      <!-- Submit button -->
      <button type="submit">Login</button>
    </form>
  </div>
</template>


<style>
.login-container {
  position: fixed; 
  top: 0;
  left: 0;
  width: 100vw; 
  height: 100vh; 
  display: flex;
  justify-content: center;
  align-items: center;
  background: rgba(0, 0, 0, 0.5); 
}

.login-container form {
  background: white;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.login-container input {
  margin: 15px;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
}

.login-container button {
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  background-color: #007bff;
  color: white;
  cursor: pointer;
}

.login-container button:hover {
  background-color: #0056b3;
}
</style>