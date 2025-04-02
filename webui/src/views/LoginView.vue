<!-- In questa pagina l'utente effettua il login -->

<script>
export default {
  data() {
    return {
      // Username input dell'utente che si sta loggando
      username: "",
      errorMsg: "",

      // Verifica per il campo username
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
    // Funzione per effettuare il login
    async doLogin() {
      try {
        // Controlla che l'username sia valido
        if (this.username.length < 1 || this.username.length > 15) throw "Invalid username, it must contains min 1 character and max 15 characters"
        if (!this.usernameValidation.test(this.username)) throw "Invalid username, it must contain only letters and numbers"

        // Effettua la richiesta di login al server con l'username inserito (se l'username non esiste, verrà creato un nuovo utente)
        let response = await this.$axios.post('/session', {
          username: this.username,
        });

        // Salva i dati dell'utente nella sessionStorage
        sessionStorage.userID = response.data.user.userId;
        sessionStorage.username = response.data.user.username;
        sessionStorage.token = response.data.user.userId;
        sessionStorage.photo = response.data.user.userPhoto;
        // Reindirizza l'utente alla home
        this.$router.push("/home");
        // Emette l'evento di login avvenuto con successo
        this.$emit('login-success');
      } catch (e) {
        this.errorMsg = e.toString();
      };
    }
  },
  mounted() {
    // Se l'utente è già loggato, reindirizza alla home
    if (sessionStorage.token) {
      this.$router.push("/home");
      return;
    }
    // Altrimewnti cancella i dati dell'utente dalla sessionStorage
    sessionStorage.clear();
  },
}

</script>

<template>
  <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
  <div class="login-container">
    <form @submit.prevent="doLogin">
      <h1>WasaText</h1>
      <input type="text" v-model="username" placeholder="Enter your username" required/>
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