<script>
import groupInfo from "../components/groupInfo.vue";
export default {
    data: function() {
        return {
            errorMsg: "",
            groupId: Number(sessionStorage.groupID),
            recipientPhoto: sessionStorage.recipientPhoto,
            recipientname: sessionStorage.recipientName,
            members: sessionStorage.members,
            username: sessionStorage.username,
            messages: [],
            newMessage: "",
            newPhoto: null,
            showGroupInfo: false,
        }
    },
    emits: ['to-home', 
            'login-success', 
            'update-username', 
            'close', 
            'update-photo', 
            'update-groupname', 
            'update-group-photo',
            'update-group-members'],
    components: {
        groupInfo
    },
    methods: {
        updateGroupMembers() {
            this.members = sessionStorage.members;
        },
        updateGroupname() {
            this.recipientname = sessionStorage.recipientName;
        },
        updateGroupPhoto() {
            this.recipientPhoto = sessionStorage.recipientPhoto;
        },
        removeFile() {
            this.newPhoto = null;
        },
        fileInput(){
            this.$refs.file.click();
        },
        handleFileChange(event) {
            const file = event.target.files[0]; 
            if (!file) {
            this.errorMsg = "Nessun file selezionato";
            return;
            }
            if (file.type !== "image/jpeg" && file.type !== "image/jpg" && file.type !== "image/gif") {
            this.errorMsg = "File type not supported, only jpg, jpeg and gif are allowed";
            return;
            }
            if (file.size > 5242880) {
            this.errorMsg = "File size is too big. Max size is 5MB";
            return;
            }
            this.newPhoto = file;
        },

        goToHome() {
            if (sessionStorage.groupID) {
                sessionStorage.removeItem("groupID");
                sessionStorage.removeItem("membersOfGroup");
                sessionStorage.removeItem("members");
            }
            sessionStorage.removeItem("recipientName");
            sessionStorage.removeItem("recipientPhoto");
            this.$router.push("/home");
        },
        async getMessages() {
            this.errorMsg = "";
            this.messages = [];
            try{
				let response = await this.$axios.get(`/users/${sessionStorage.userID}/conversations/${sessionStorage.convId}`, { headers: { 'Authorization': `${sessionStorage.token}` } });
				this.messages = response.data;
			} catch (e) {
				this.errorMsg = e.toString();
				document.getElementsByTagName("input")[0].style.outline = "auto";
                document.getElementsByTagName("input")[0].style.outlineColor = "red";
			}
        },
        async sendMessage() {
            this.errorMsg = "";
            const formData = new FormData();
            if (this.newPhoto && (this.newPhoto.type === "image/jpeg" || this.newPhoto.type === "image/jpg")) {
                formData.append('image', this.newPhoto);
            }
            if (this.newPhoto && (this.newPhoto.type === "image/gif")) {
                formData.append('gif', this.newPhoto);
            }
            if (this.newMessage) {
                formData.append('text', this.newMessage);
            }
            if (this.newMessage && this.newPhoto && (this.newPhoto.type === "image/gif")) {
                this.errorMsg = "You must write a message or select a photo/gif";
                return;
            }
            try{
                let response = await this.$axios.post(`/users/${sessionStorage.userID}/conversations/${sessionStorage.convId}/messages`,
                formData, { headers: { 'Authorization': `${sessionStorage.token}` }});
                this.newMessage = "";
                this.newPhoto = null;
                this.getMessages();
            } catch (e) {
                this.errorMsg = e.toString();
                document.getElementsByTagName("input")[0].style.outline = "auto";
                document.getElementsByTagName("input")[0].style.outlineColor = "red";
            }
        },
        groupInfo() {
            this.showGroupInfo = !this.showGroupInfo;
        },
    },
    mounted() {
        this.getMessages();
        this.intervalId = setInterval(async () => {
			clearInterval(this.intervalId);
			await this.getMessages();
			this.intervalId = setInterval(this.getMessages, 5000);
		}, 5000);
    }

}
</script>

<template>
    <div class="chat-container">
        <div v-if="groupId != 0" class="chat-header" @click="groupInfo">
            <button @click="goToHome" class="back-button"> 
                <svg class="feather"> 
                    <use href="/feather-sprite-v4.29.0.svg#chevron-left" />
                </svg>
            </button>
            <img :src="`data:image/jpg;base64,${recipientPhoto}`" alt="Conversation photo" class="recipient-photo">
            <div class="recipient-info" @click="goToUserInfo">
                <h2 class="recipient-name">{{ recipientname }}</h2>
                <p class="members-list">{{ members }}</p>
            </div>
        </div>
        <div v-else class="chat-header">
            <button @click="goToHome" class="back-button"> 
                <svg class="feather"> 
                    <use href="/feather-sprite-v4.29.0.svg#chevron-left" />
                </svg>
            </button>
            <img :src="`data:image/jpg;base64,${recipientPhoto}`" alt="Conversation photo" class="recipient-photo">
            <div class="recipient-info" @click="goToUserInfo">
                <h2 class="recipient-name">{{ recipientname }}</h2>
            </div>
        </div>
        <div class="chat-body">
            <div class="chat-messages">
                <div v-for="object in messages" :key="object.message.messageId" :class="{'message': true, 'user-message': object.sender.username === username, 'other-message': object.sender.username !== username}">
                    <div v-if="groupId != 0 && object.sender.username !== username" class="message-header">
                        <h4 >{{ object.sender.username }}</h4>
                        <p>{{ object.dateTime }}</p>
                    </div>
                    <div v-else class="message-header" >
                        <p>{{ object.dateTime }}</p>
                    </div>
                    <div v-if="!object.message.text && object.message.photo" class="message-body">
                        <img class="img" :src="`data:image/jpg;base64,${object.message.photo}`" alt="Message image" />
                    </div>
                    <div v-if="object.message.text && object.message.photo" class="message-body">
                        <img class="img" :src="`data:image/jpg;base64,${object.message.photo}`" alt="Message image" />
                        <p>{{ object.message.text }}</p>
                    </div>
                    <div v-if="object.message.gif" class="message-body">
                        <img class="gif":src="`data:image/gif;base64,${object.message.gif}`" alt="Message gif" />
                    </div>
                    <div v-if="object.message.text && !object.message.photo" class="message-body">
                        <p>{{ object.message.text }}</p>
                    </div>
                </div>
            </div>
        </div>
        <div class="chat-footer">
            <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
            <form @submit.prevent="sendMessage">
                <input @keydown.enter.prevent="sendMessage" type="text" v-model="newMessage" placeholder="Type a message" />
                <input type="file" ref="file" accept=".jpg,.jpeg,.gif" @change="handleFileChange" style="display: none;"/>
                <button type="button" @click="fileInput" class="file-button" v-if="!newPhoto">
                    <svg class="feather"> 
                        <use href="/feather-sprite-v4.29.0.svg#paperclip" />
                </svg>
                </button>
                <button class="remove-file-button" @click="removeFile" v-if="newPhoto">
                    <svg class="feather"> 
                        <use href="/feather-sprite-v4.29.0.svg#x" />
                    </svg>
                </button>
                <button class="send-button" type="submit">Send</button>
            </form>
        </div>
    </div>
    <groupInfo :show1="showGroupInfo" @close="groupInfo" @update-groupname="updateGroupname" @update-group-photo="updateGroupPhoto" @update-group-members="updateGroupMembers"></groupInfo>
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

.members-list {
    font-size: 14px;
    color: #777;
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
.message {
    display: flex;
    flex-direction: column;
    gap: 5px;
    padding: 8px; 
    border-radius: 5px;
    background-color: #f9f9f9;
    min-width: 300px;
    max-width: 80%;
    margin-bottom: 10px;
}
.message.user-message {
    align-self: flex-end; 
    background-color: #007bff; 
    color: white; 
}
.message.other-message {
    align-self: flex-start;
    background-color: #f0f0f0;  
    color: black; 
}
.message-header {
    display: flex;
    justify-content: space-between;
}
.message-header h4 {
    font-size: 20px; 
    font-weight: normal;
    color:#007bff;
}

.message-header p {
    font-size: 14px; 
    text-align: right;
    margin-left: auto;
}
.message-body {
    display: flex;
    justify-content: space-between;
    gap: 10px;
    flex-direction: column;
}

.message-body p {
    font-size: 16px; 
    color: inherit;
}

.img {
    width: 500px;
    height: auto;
    border-radius: 5px;
}

.gif {
    width: 200px;
    height: auto;
    border-radius: 5px;
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
.remove-file-button {
    padding: 10px 20px;
    border: none;
    border-radius: 5px;
    background-color: #ff0000;
    color: white;
    cursor: pointer;
}
.file-button {
    padding: 10px 20px;
    border: 2px solid #007bff(0, 0, 0, 0.2);
    border-radius: 5px;
    background-color: white;
    color:#333;
    cursor: pointer;
}

</style>