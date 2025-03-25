<script>
export default {
    props: {
        show: Boolean,
    },
    data: function() {
        return {
            usernameValidation: new RegExp('^[a-z0-9]{1,15}$'),
            username: "",
        }
    },
    emits: ['to-home', 'login-success', 'update-username'],
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

                // Chiudi il modale
                this.closeModal();

                // Emette l'evento di login avvenuto con successo
                this.$emit('update-username');
            } catch (e) {
                this.errorMsg = e.toString();
                document.getElementsByTagName("input")[0].style.outline = "auto";
                document.getElementsByTagName("input")[0].style.outlineColor = "red";
            };

        }
}
}</script>

<template>
    <div v-if="show" class="modal-mask" >
        <div class="modal-wrapper">
            <div class="modal-container">
                <div class="modal-header">
                    <h3>Set new Username</h3>
                    <button @click="closeModal">X</button>
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
    }
    .modal-container {
    width: 300px;
    margin: 0px auto;
    padding: 20px;
    background-color: white;
    border-radius: 10px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.33);
    transition: all 0.3s ease;
    }
    .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    }
    
</style>