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
            control: false,
            userId: Number(sessionStorage.userID),
            infoGroupMembers: [],
            infoGroupMembersQuantity: 0,
            option: false,
            showForwardBool: false,
            conversations: [],
            conversationsFiltered: [],
            conversationsSelected: [],
            searchText: "",
            filteredUsers: [],
            errorMsg2: "",
            errorMsg3: "",
            usernameValidation: new RegExp('^[a-z0-9]{1,15}$'),
            selectedUsers: [],
            messId: 0,
            selectedMessageId: 0,
            showReplyTo: false,
            showCom: false,
            comments: [],
            HaveICommented: false,
            emojis: ["😀", "😂", "😍", "😎", "😭", "😡", "🎉", "❤️", "👍", "🔥"],
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
        groupInfo
    },
    watch: {
		searchText() {
			this.filterUsers();
		},
	},
    methods: {
        updateData(data) {
            this.infoGroupMembers = data;
            this.infoGroupMembersQuantity = data.length + 1;
        },
        async commentMessage(input) {
            for (let i = 0; i < this.messages.length; i++) {
                if (this.messages[i].message.messageId === this.messId) {
                    if (this.messages[i].sender.userId === this.userId) {
                        this.errorMsg3 = "You can't comment your own message!";
                        return;
                    }
                }
            }
            this.errorMsg3 = "";
            try {
                let response = await this.$axios.put(`/users/${sessionStorage.userID}/conversations/${sessionStorage.convId}/messages/${this.messId}/comments`, {
                    emoji: input,
                }, { headers: { 'Authorization': `${sessionStorage.token}` }});
                this.getMessages();
                const newComment = response.data;
                for (let i = 0; i < this.comments.length; i++) {
                    if (this.comments[i].sender.userId === this.userId) {
                        this.HaveICommented = true;
                        this.comments[i].content = newComment.content;
                        break;
                    }
                }
                if (!this.HaveICommented) {
                    this.comments.push(newComment);
                    this.HaveICommented = true;
                }
            } catch (e) {
                this.errorMsg3 = e.toString();
            }
        },
        async uncommentMessage() {
            let commentId = 0;
            for (let i = 0; i < this.messages.length; i++) {
                if (this.messages[i].message.messageId === this.messId) {
                    if (this.messages[i].sender.userId === this.userId) {
                        this.errorMsg3 = "You haven't commented this message, and can't comment your own message!";
                        return;
                    }
                }
            }
            if (!this.HaveICommented) {
                this.errorMsg3 = "You haven't commented this message!";
                return;
            }
            else {
                for (let i = 0; i < this.comments.length; i++) {
                    if (this.comments[i].sender.userId === this.userId) {
                        commentId = this.comments[i].commentId;
                        break;
                    }
                }
            }
            this.errorMsg3 = "";
            try {
                await this.$axios.delete(`/users/${sessionStorage.userID}/conversations/${sessionStorage.convId}/messages/${this.messId}/comments/${commentId}`, { headers: { 'Authorization': `${sessionStorage.token}` }});
                this.getMessages();
                for (let i = 0; i < this.comments.length; i++) {
                    if (this.comments[i].sender.userId === this.userId) {
                        this.comments.splice(i, 1);
                        this.HaveICommented = false;
                        break;
                    }
                }
            } catch (e) {
                this.errorMsg3 = e.toString();
            }
        },
        showComments(object) {
            this.messId = object.message.messageId;
            this.HaveICommented = false;
            this.showCom = !this.showCom;
            if (object.comments == null) {
                this.comments = [];
                return;
            }
            this.comments = object.comments;
            for (let i = 0; i < this.comments.length; i++) {
                if (this.comments[i].sender.userId === this.userId) {
                    this.HaveICommented = true;
                    const comment = this.comments.splice(i, 1)[0];
                    this.comments.push(comment);
                    break;
                }
            }
        },
        showComments2() {
            this.showCom = !this.showCom;
            this.errorMsg3 = "";
        },
        selectMessage() {
            this.showReplyTo = true;
            this.selectedMessageId = this.messId;
            this.option = false;
        },
        unselectMessage() {
            this.showReplyTo = false;
            this.selectedMessageId = 0;
        },
        async deleteMessage() {
            this.errorMsg = "";
            try {
                this.$axios.delete(`/users/${sessionStorage.userID}/conversations/${sessionStorage.convId}/messages/${this.messId}`, { headers: { 'Authorization': `${sessionStorage.token}`} });
                this.option = false;
                this.getMessages();
            } catch (e) {
                this.errorMsg = e.toString();
            }
        },

        async ForwardMessage() {
            this.errorMsg2 = "";
            if (this.conversationsSelected.length === 0 && this.selectedUsers.length === 0) {
                this.errorMsg2 = "You must select at least one conversation or user!";
                return;
            }
            let object = {
                user: 0,
                group: 0,
            }
            let destination = [];
            for (let i = 0; i < this.conversationsSelected.length; i++) {
                object = {
                    user: 0,
                    group: this.conversationsSelected[i].conversation.group,
                }
                destination.push(object);
            }
            object = {
                user: 0,
                group: 0,
            }
            for (let i = 0; i < this.selectedUsers.length; i++) {
                object = {
                    user: this.selectedUsers[i].userId,
                    group: 0,
                }
                destination.push(object);
            }
            object = {
                user: 0,
                group: 0,
            }
            this.conversationsSelected = [];
            this.selectedUsers = [];
            this.conversationsFiltered = [];
            this.filteredUsers = [];
            this.searchText = "";
            const input = {
                destination: destination,
            }
            try { 
                this.$axios.post(`/users/${sessionStorage.userID}/conversations/${sessionStorage.convId}/messages/${this.messId}`, input, { headers: { 'Authorization': `${sessionStorage.token}` } });
                this.showForwardBool = false;
                this.option = false;
            } catch (e) {
                this.errorMsg = e.toString();
            };
        },
        selectConversation(conversation) {
            if (this.conversationsSelected.some(conv => conv.conversation.convId === conversation.conversation.convId)) {
                this.errorMsg2 = "You have already selected this conversation!";
                return;
            }
            this.conversationsSelected.push(conversation);
        },

        removeGroup(index) {
            this.conversationsSelected.splice(index, 1);
            if (this.errorMsg2) {
                this.errorMsg2 = "";
            }
        },
        removeMember(index) {
			this.selectedUsers.splice(index, 1);
			if (this.errorMsg2) {
				this.errorMsg2 = "";
			}
		},
        async filterUsers() {
            if (this.errorMsg2 !== "It's not necessary that you select yourself!" || this.searchText.length !== 0) {
                this.errorMsg2 = "";
            }
			this.filteredUsers = [];

			if (this.searchText.length > 0) {
				if (this.searchText.length > 15 || !this.usernameValidation.test(this.searchText)) {
				this.errorMsg2 = "Invalid username, it can contain only letters and numbers for a maximum of 16 characters.";
				this.filteredUsers = [];
				return;
				}
				try {
					let response = await this.$axios.get(`/users?query=${this.searchText}`, { headers: { 'Authorization': `${sessionStorage.token}` } });
					if (response.data == null) {
					this.filteredUsers = [];
					return;
					}
					this.filteredUsers = response.data;
				} catch (e) {
					this.errorMsg2 = e.toString();
					this.filteredUsers = [];
				}
			}
		},

        selectUser(user) {
			if (user.userId === Number(sessionStorage.userID)) {
				this.errorMsg2 = "It's not necessary that you select yourself!";
			}
			else if (!this.selectedUsers.some(member => member.userId === user.userId)) {
				this.selectedUsers.push(user);
			}
			this.searchText = "";  
			this.filteredUsers = [];  
			},

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
            this.control = true;
            this.$router.push("/home");
        },
        async getMessages() {
            if (this.control) {
                return;
            }
            this.errorMsg = "";
            this.messages = [];
            try{
				let response = await this.$axios.get(`/users/${sessionStorage.userID}/conversations/${sessionStorage.convId}`, { headers: { 'Authorization': `${sessionStorage.token}` } });
				this.messages = response.data;
			} catch (e) {
				this.errorMsg = e.toString();
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
            if (this.selectedMessageId != 0) {
                formData.append('messToReplyTo', this.selectedMessageId);
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
                this.showReplyTo = false;
                this.selectedMessageId = 0;
            } catch (e) {
                this.errorMsg = e.toString();
            }
        },
        groupInfo() {
            this.showGroupInfo = !this.showGroupInfo;
            if (sessionStorage.membersOfGroup) {
                this.infoGroupMembers = JSON.parse(sessionStorage.membersOfGroup);
            }
            this.infoGroupMembersQuantity = this.infoGroupMembers.length+1;
        },
        showOption(object) {
            this.messId = object.message.messageId;
            this.option = !this.option;
        },

        showOption2() {
            this.option = !this.option;
        },
        showForward() {
            this.showForwardBool = !this.showForwardBool;
            this.getConversations();
            this.option = !this.option;
        },

        async getConversations() {
			this.errorMsg = "";
			this.conversations = [];
			try{
				let response = await this.$axios.get(`/users/${sessionStorage.userID}/conversations`, { headers: { 'Authorization': `${sessionStorage.token}` } });
				this.conversations = response.data;
                for (let i = 0; i < this.conversations.length; i++) {
                    if (this.conversations[i].conversation.group != 0) {
                        this.conversationsFiltered.push(this.conversations[i]);
                    }
                }

			} catch (e) {
				this.errorMsg = e.toString();
			}
		}
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
            <div class="recipient-info">
                <h2 class="recipient-name">{{ recipientname }}</h2>
            </div>
        </div>
        <div class="chat-body">
            <div class="chat-messages">
                <div v-for="object in messages" @click="showOption(object)" :key="object.message.messageId" :class="{'message': true, 'user-message': object.sender.username === username, 'other-message': object.sender.username !== username}">
                    <div v-if="object.message.replyId" class="replyToMessage">
                        <span class="replyToUser" v-if="messages.find(mes => mes.message.messageId === object.message.replyId).sender.userId === userId">Me</span>
                        <span class="replyToUser" v-else>{{ messages.find(mes => mes.message.messageId === object.message.replyId).sender.username }}</span>
                        <img v-if="messages.find(mes => mes.message.messageId === object.message.replyId).message.photo" :src="`data:image/jpg;base64,${messages.find(mes => mes.message.messageId === object.message.replyId).message.photo}`" alt="Reply to message image" />
                        <span v-if="messages.find(mes => mes.message.messageId === object.message.replyId).message.text">{{ messages.find(mes => mes.message.messageId === object.message.replyId).message.text }}</span>
                        <img v-if="messages.find(mes => mes.message.messageId === object.message.replyId).message.gif" :src="`data:image/gif;base64,${messages.find(mes => mes.message.messageId === object.message.replyId).message.gif}`" alt="Reply to message gif" />
                    </div>
                    <div v-if="groupId != 0 && object.sender.username !== username" class="message-header">
                        <h4 >{{ object.sender.username }}</h4>
                        <div class="forwarded-date">
                            <p class="only-date">{{ object.dateTime }}</p>
                            <div v-if="object.message.forward" class="forwarded-message">
                                <svg class="feather"> 
                                    <use href="/feather-sprite-v4.29.0.svg#message-square" />
                                </svg>
                                <p>Forwarded message</p>
                            </div>
                        </div>
                    </div>
                    <div v-else class="message-header" >
                        <div class="forwarded-date">
                            <p class="only-date">{{ object.dateTime }}</p>
                            <div v-if="object.message.forward" class="forwarded-message">
                                <svg class="feather"> 
                                    <use href="/feather-sprite-v4.29.0.svg#message-square" />
                                </svg>
                                <p>Forwarded message</p>
                            </div>
                        </div>    
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
                    <div class="checkmark-and-comments">
                        <p v-if="object.message.senderId === userId" class="comments-user" @click.stop="showComments(object)">See all comments</p>
                        <p v-else class="comments-dest" @click.stop="showComments(object)">See all comments</p>
                        <div v-if="!object.message.status && object.sender.userId === userId" class="checkmark">
                            <svg class="feather"> 
                                <use href="/feather-sprite-v4.29.0.svg#check" />
                            </svg>
                        </div>
                        <div v-if="object.message.status && object.sender.userId === userId" class="checkmark">
                            <svg class="feather"> 
                                <use href="/feather-sprite-v4.29.0.svg#check" />
                            </svg>
                            <svg class="feather"> 
                                <use href="/feather-sprite-v4.29.0.svg#check" />
                            </svg>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="chat-footer">
            <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
            <form @submit.prevent="sendMessage">
                <input @keydown.enter.prevent="sendMessage" type="text" v-model="newMessage" placeholder="Type a message" />
                <input type="file" ref="file" accept=".jpg,.jpeg,.gif" @change="handleFileChange" style="display: none;"/>
                <button class="remove-file-button" @click="unselectMessage" v-if="showReplyTo">
                    <svg class="feather"> 
                        <use href="/feather-sprite-v4.29.0.svg#trash-2" />
                    </svg>
                </button>
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
    

    <div v-if="option" @click="showOption2" class="fullscreen-container">
		<div v-if="option" @click.stop class="message-option-container">
			<!-- Button to the select the message as cause of reply -->
			<button type="button" class="option-button" @click="selectMessage">Reply to</button>
            <!-- Button to forward the message -->
			<button type="button" class="option-button" @click="showForward">Forward</button>
            <!-- Button to delete message -->
			<button type="button" class="option-button" @click="deleteMessage">Delete</button>
		</div>
	</div>


    <div v-if="showCom" @click="showComments2" class="fullscreen-container">
        <div @click.stop class="comments-option-container">
            <div class="comments">
                <p>Message comments</p>
                <div v-if="comments.length === 0" class="no-comments">
                    <p>No comments</p>
                </div>
                <div v-for="object in comments" :key="object.commentId" class="comment">
                    <h4 v-if="object.sender.userId !== userId">{{ object.sender.username }}:</h4>
                    <h4 v-if="object.sender.userId === userId">Me:</h4>
                    <p >{{ object.content }}</p>
                </div>
            </div>
            <div class="comment-input">
                <h3>Comment message!</h3>
                <div v-for="emoji in emojis" :key="emoji" class="emoji" @click="commentMessage(emoji)">
                    {{ emoji }}
                </div>
                <div class="emoji" @click="uncommentMessage()">
                    <svg class="feather"> 
                        <use href="/feather-sprite-v4.29.0.svg#trash-2" />
                    </svg>
                </div>
            </div>
            <ErrorMsg v-if="errorMsg3" :msg="errorMsg3"></ErrorMsg>
        </div>
    </div>


    <div v-if="showForwardBool" @click="showForward" class="fullscreen-container">
        <div @click.stop class="menu-forward">
            <h3 class="menu-forward-title">Select conversations to forward the message</h3>
            <form @submit.prevent="ForwardMessage">
                <!-- Select conversations -->
                <div class="select-conv">
                    <!-- Select groups -->
                    <div class="all-conversations">
                        <p class="instruction">Your groups</p>
                        <button type="button" v-for="conversation in conversationsFiltered" :key="conversation.conversation.userId" @click="selectConversation(conversation)" class="conversation">
                            <img v-if="conversation.conversation.group != 0" :src="`data:image/jpg;base64,${conversation.group.groupPhoto}`" alt="Conversation photo" >
					        <p class="conversation-name" v-if="conversation.conversation.group != 0">{{ conversation.group.username }}</p>
						</button>
                        <ul>
                            <!-- Group who is selected -->
                            <li v-for="(conv, index) in conversationsSelected" :key="index">
                                <p>{{ conv.group.username }}</p>
                                <button @click.prevent="removeGroup(index)">x</button>
                            </li>
                        </ul>
                    </div>
                    <p class="instruction">Other users</p>
                    <!-- Select users -->
                    <div class="group-members">
                        <div class="input-members">
                            <input type="text" class="form-control" v-model="searchText" placeholder="Enter a username" />
                        </div>
                            <!-- Search and print all users who have username that start with the text in the input space -->
                            <div class="search-results">
                            <div v-for="user in filteredUsers" :key="user.userId" @click="selectUser(user)" class="user">
                                <p>{{ user.username }}</p>
                            </div>
                        </div>	
                    </div>
                    <ErrorMsg v-if="errorMsg2" :msg="errorMsg2"></ErrorMsg>	
                </div>
                <div class="usersChoosed">
                    <ul>
                        <!-- User who is selected -->
                        <li class= "userChoosed "v-for="(conv, index) in selectedUsers" :key="index">
                            <p>{{ conv.username }}</p>
                            <button @click.prevent="removeMember(index)">x</button>
                        </li>
                    </ul>
                </div>
                <hr class="separator">
                <button class="buttonSubmit" type="submit">Forward</button>
            </form>
        </div>
    </div>

    <groupInfo :members="infoGroupMembers" :membersQuantity="infoGroupMembersQuantity" @update-group-info="updateData" :show1="showGroupInfo" @close="groupInfo" @update-groupname="updateGroupname" @update-group-photo="updateGroupPhoto" @update-group-members="updateGroupMembers" />
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
    cursor: pointer;
}

.replyToMessage {
    display: flex;
    flex-direction: column;
    margin: 5px 0;
    padding: 5px 10px;
    border-radius: 5px;
    background: white;
    border-left: 3px solid #a6a6a6;
    font-size: 12px;
    color: #6c757d;
    
}

.replyToUser {
    color: #007bff;
}

.replyToMessage img {
    width: 50px;
    height: 50px;
    border-radius: 5px;
    object-fit: cover;
} 

.message.user-message {
    align-self: flex-end; 
    background-color: #007bff; 
    color: white; 
}
.message.other-message {
    align-self: flex-start;
    background-color: #e0e0e0;  
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

.checkmark-and-comments {
    display: flex;
    width: 100%;
    height: 20px;
}

.checkmark {
    margin-left: auto;
}

.comments-user {
    margin-right: auto;
    font-size: 14px;

}

.comments-dest {
    margin-left: auto;
    font-size: 14px;
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

.fullscreen-container {
    position: fixed;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
}

.message-option-container {
    position: fixed;
    width: 300px;
    height: 120px;
    background-color: #f9f9f9;
    display: flex;
    flex-direction: column;
    border-radius: 5px;
    z-index: 10001;
}

.option-button {
    padding: 10px 20px;
    border: none;
    border-radius: 5px;
    background-color: white;
    color: black;
    cursor: pointer;
    width: 100%;
    text-align: center;
    transition: background 0.3s;
}

.option-button:hover {
    background-color: #007bff;
    color: white;
}

.menu-forward {
    position: fixed;
    width: 700px;
    height: 800px;
    background-color: #f9f9f9;
    display: flex;
    flex-direction: column;
    border-radius: 5px;
    z-index: 10001;
    text-align: center;
}

.menu-forward-title {
    font-size: 25px;
    color: #333;
    margin: 20px;
}

.select-conv {
    display: flex;
    flex-direction: column;
    gap: 10px;
    padding: 10px;
}

.all-conversations {
    display: flex;
    flex-direction: column;
    gap: 5px;
    padding: 5px;

}

.instruction {
    font-size: 15px;
    color: #333;
    align-self: flex-start;
    padding-left: 30px;
    margin: 5px;
}

.conversation {
    display: flex;
    cursor: pointer;
    border-radius: 5px;
    border: none;
    padding: 5px;
}


.conversation img {
    width: 50px;
    height: 50px;
    border-radius: 50%;
    object-fit: cover;
}

.conversation-name {
    font-size: 20px;
    color: black;
    font-weight: normal;
}

.group-members {
    display: flex;
    flex-direction: column;
    padding: 5px;
}

.input-members {
    display: flex;
    gap: 5px;
}

.form-control {
    border: 1px solid #ccc;
    border-radius: 5px;
}

.search-results {
    display: flex;
    flex-direction: column;
    gap: 5px;
}

.user {
    display: flex;
    background-color: #f9f9f9;
    cursor: pointer;
    border-radius: 5px;
}

.user p {
    font-size: 16px;
    color: black;
    font-weight: normal;
    padding: 5px;
}


.separator {
    margin: 20px;
}

.usersChoosed {
    display: flex;
    flex-direction: column;
    gap: 15px;
    padding: 15px;
}

.userChoosed {
    display: flex;
    gap: 5px;
    padding: 5px;
    border-radius: 5px;
    background-color: #f9f9f9;
    
}

.forwarded-message {
    display: flex;
    flex-direction: row;

}

.forwarded-message p {
    font-size: 10px;
    font-style: oblique;
}

.forwarded-message svg {
    width: 15px;
    height: 15px;
}

.forwarded-date {
    display: flex;
    flex-direction: column;
    width: 100%;
    align-items: flex-end;
}

.only-date {
    margin-bottom: 2px;
}

.comments-option-container {
    background: white;
    width: 50%;
    max-width: 600px;
    border-radius: 10px;
    padding: 20px;
    display: flex;
    flex-direction: column;
}

.comments {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.comments p {
    margin-bottom: 10px;
    font-size: 20px;
    font-weight: bold;
}

.no-comments {
    text-align: center;
    font-size: 14px;
}

.comment {
    display: flex;
    align-items: center;
    width: 100%;
    gap: 10px;
    margin-bottom: 10px;
    padding: 10px;
    border-radius: 5px;
    background: #f1f1f1;
    justify-content: center;
}

.comment h4 {
    font-size: 14px;
    font-weight: bold;
    margin: 0;
}

.comment p {
    font-size: 14px;
    margin: 0;
}

.comment-input {
    margin-top: 20px;
    text-align: center;
}

.comment-input h3 {
    margin-bottom: 10px;
    font-size: 15px;
}

.emoji {
    display: inline-block;
    padding: 5px;
    font-size: 1.5em;
    cursor: pointer;
    transition: transform 0.2s ease-in-out;
}

.emoji:hover {
    transform: scale(1.2);
}

.emoji svg {
    width: 20px;
    height: 20px;
    fill: #555;
}

.emoji:hover svg {
    fill: red;
}
</style>