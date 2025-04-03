<script>
export default {
    props: {
        show: Boolean,
    },
    data: function() {
        return {
            usernameValidation: new RegExp('^[a-z0-9]{1,15}$'),
            username: "",
            errorMsg: "",
        }
    },
    watch: {
		username() {
			this.errorMsg = "";
		},
    },
    emits: ['to-home', 
            'login-success', 
            'update-username', 
            'close', 
            'update-photo', 
            'update-groupname', 
            'update-group-photo',
            'update-group-info'],
    methods: {
        closeModal() {
            this.$emit('close');
        },
        async setNewUsername() {
            try {
                // Controlla che l'username sia valido
                if (this.username.length < 1 || this.username.length > 15) throw "Invalid username, it must contains min 1 character and max 15 characters"
                if (!this.usernameValidation.test(this.username)) throw "Invalid username, it must contain only letters and numbers"

                // Effettua la richiesta di login al server con l'username inserito (se l'username non esiste, verrà creato un nuovo utente)
                let response = await this.$axios.put(`/users/${sessionStorage.userID}/username`, {
                username: this.username,
                }, {headers: {Authorization: `${sessionStorage.token}`}});
                
                // Salva i dati dell'utente nella sessionStorage
                sessionStorage.username = response.data.username;

                // Emette l'evento di login avvenuto con successo
                this.$emit('update-username');

                // Chiudi il modale
                this.closeModal();
            } catch (e) {
                if (e.response && e.response.status === 400) {
                    this.errorMsg = "Username already taken, please choose another one";
                } else {
                    this.errorMsg = "An error occurred, please try again.";
                }
            };
        }
    }
}</script>

<template>
    <div v-if="show" class="modal-mask" >
        <div class="modal-wrapper">
            <div class="modal-container">
                <button class="exit-button" @click="closeModal">X</button>
                <div class="modal-header">
                    <form @submit.prevent="setNewUsername">
                        <h3>Set new username</h3>
                        <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                        <input type="text" v-model="username" placeholder="Enter your new username" />
                        <button type="submit">confirm</button>
                    </form>
                </div>
            </div>
        </div>        
    </div>
</template>

<style>
    .modal-mask {
        position: fixed;
        z-index: 9998;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-color: rgba(0, 0, 0, 0.5);
        display: table;
        transition: opacity 0.3s ease;
    }
    .modal-wrapper {
        display: table-cell;
        vertical-align: middle;
        margin-left: 280px;
    }
    .modal-container {
        position: relative;
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 600px;
        margin: 0px auto;
        padding: 20px;
        background-color: white;
        border-radius: 10px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.33);
        transition: all 0.3s ease;
    }

    .modal-header {
        display: flex;
        justify-content: center;
        align-items: center;
    }

    .modal-header form {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 10px;
    }
    .modal-header input {
        padding: 10px;
        border: 1px solid #ccc;
        border-radius: 5px;
    }

    .modal-header button {
        padding: 10px 20px;
        border: none;
        border-radius: 5px;
        background-color: #007bff;
        color: white;
        cursor: pointer;
    }
    .exit-button {
        position: absolute;
        top: 5%;
        right: 1.7%;
        padding: 5px 10px;
        border: none;
        border-radius: 50%;
        background-color: #ff0000;
        color: white;
        cursor: pointer;
    }
</style>