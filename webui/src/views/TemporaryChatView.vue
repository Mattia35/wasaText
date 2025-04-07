<script>
export default {
    data: function() {
        return {
            errorMsg: "",
            recipientPhoto: sessionStorage.recipientPhoto,
            recipientname: sessionStorage.recipientName,
            newMessage: "",
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
    methods: {
        // function to go to the home page
        goToHome() {
            // remove the recipient name and photo from sessionStorage
            sessionStorage.removeItem("recipientName");
            sessionStorage.removeItem("recipientPhoto");
            // go to the home page
            this.$router.push("/home");
        },
        // function to create a new private conversation
        async createPrivateConversation() {
            this.errorMsg = "";
            // try the request
            try{
                // check if the mesage is empty
                if (this.newMessage.trim() === "") throw "message is empty!";
                // make the request
                let response = await this.$axios.put(`/users/${sessionStorage.userID}/conversations`,
                {
                    user: this.recipientname,
                    text: this.newMessage,
                }, { headers: { 'Authorization': `${sessionStorage.token}` }});
                this.newMessage = "";
                // get the response (array of conversations)
                let conversations = response.data;
                // get the conversation that is just created
                for (let i = 0; i < conversations.length; i++) {
                    if (conversations[i].user.username === this.recipientname) {
                        // save the conversation id in sessionStorage
                        sessionStorage.convId = conversations[i].conversation.convId;
                        break;
                    }
                }
                // go to the conversation page
                this.$router.push(`/chat/${sessionStorage.convId}`);
            } catch (e) {
                // save and print the error message
                this.errorMsg = e.toString();
            }
        },
    },
    mounted() {
    }

}
</script>

<template>
    <!-- container of the chat -->
    <div class="chat-container">
        <!-- header of the chat -->
        <div class="chat-header">
            <!-- back button to go to the home page -->
            <button @click="goToHome" class="back-button"> 
                <svg class="feather"> 
                    <use href="/feather-sprite-v4.29.0.svg#chevron-left" />
                </svg>
            </button>
            <!-- recipient photo and name -->
            <img :src="`data:image/jpg;base64,${recipientPhoto}`" alt="Conversation photo" class="recipient-photo">
            <div class="recipient-info">
                <h2 class="recipient-name">{{ recipientname }}</h2>
            </div>
        </div>
        <!-- body of the chat (empty)-->
        <div class="chat-body">
            <div class="chat-messages">
            </div>
        </div>
        <!-- footer of the chat (input and send button) -->
        <div class="chat-footer">
            <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
            <form @submit.prevent="createPrivateConversation">
                <input @keydown.enter.prevent="createPrivateConversation" type="text" v-model="newMessage" placeholder="Type a message" required/>
                <button class="send-button" type="submit">Send</button>
            </form>
        </div>
    </div>
</template>

<style>

.recipient-info {
    display: flex;
    flex-direction: column;
}

.chat-container {
    position: fixed; 
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background-color: white; 
    display: flex;
    flex-direction: column;
    z-index: 9999; 
}

.chat-header {
    position: fixed; 
    top: 0;
    left: 0;
    width: 100%;
    height: 100px;
    background-color: #ffffff;
    display: flex;
    align-items: center;
    padding: 10px 15px;
    border-bottom: 1px solid #ccc;
    z-index: 10000;
    cursor: pointer;
}

.back-button {
    background: none;
    border: none;
    cursor: pointer;
    padding: 5px;
    display: flex;
    align-items: center;
}

.back-button svg {
    width: 24px;
    height: 24px;
    color: #333;
}

.recipient-photo {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    object-fit: cover;
    border: 3px solid #007bff;
    margin: 0 10px;
}

.recipient-name {
    margin-top: 20px;
    font-size: 30px;
    font-weight: bold;
    color: #333;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.chat-body {
    flex: 1;
    margin-top: 100px;
    overflow-y: auto;
    padding: 10px;
    display: flex;
    flex-direction: column-reverse; 
}
.chat-messages {
    display: flex;
    flex-direction: column-reverse;
    gap: 10px;
    padding: 10px;
}

.chat-footer {
    padding: 10px;
    border-top: 1px solid #ccc;
}
.chat-footer form {
    display: flex;
    gap: 10px;
}
.chat-footer input {
    flex: 1;
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 5px;
}
.send-button {
    padding: 10px 20px;
    border: none;
    border-radius: 5px;
    background-color: #007bff;
    color: white;
    cursor: pointer;
}

</style>