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
        // watch the search text and call the filterUsers function
		searchText() {
			this.filterUsers();
		},
	},
    methods: {
        // function to update data about group members
        updateData(data) {
            this.infoGroupMembers = data;
            this.infoGroupMembersQuantity = data.length + 1;
        },

        // function to comment a message
        async commentMessage(input) {
            // check if the user is the sender of the message, so he can't comment his own message
            for (let i = 0; i < this.messages.length; i++) {
                if (this.messages[i].message.messageId === this.messId) {
                    if (this.messages[i].sender.userId === this.userId) {
                        this.errorMsg3 = "You can't comment your own message!";
                        return;
                    }
                }
            }
            this.errorMsg3 = "";
            // try the request to comment the message
            try {
                // make the request to comment the message
                let response = await this.$axios.put(`/users/${sessionStorage.userID}/conversations/${sessionStorage.convId}/messages/${this.messId}/comments`, {
                    emoji: input,
                }, { headers: { 'Authorization': `${sessionStorage.token}` }});
                // update the comments array with the function getMessages
                this.getMessages();
                // check if the user has already commented the message. If yes, update the comment, otherwise add the new comment
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
                // save and print the error message
                this.errorMsg3 = e.toString();
            }
        },

        // function to uncomment a message
        async uncommentMessage() {
            // check if the user is the sender of the message, so he can't uncomment his own message
            let commentId = 0;
            for (let i = 0; i < this.messages.length; i++) {
                if (this.messages[i].message.messageId === this.messId) {
                    if (this.messages[i].sender.userId === this.userId) {
                        this.errorMsg3 = "You haven't commented this message, and can't comment your own message!";
                        return;
                    }
                }
            }
            // check if the user has already commented the message. If not show an error message
            if (!this.HaveICommented) {
                this.errorMsg3 = "You haven't commented this message!";
                return;
            }
            // else, get the commentId of the comment to delete
            else {
                for (let i = 0; i < this.comments.length; i++) {
                    if (this.comments[i].sender.userId === this.userId) {
                        commentId = this.comments[i].commentId;
                        break;
                    }
                }
            }
            this.errorMsg3 = "";
            // try the request to uncomment the message
            try {
                // make the request to uncomment the message
                await this.$axios.delete(`/users/${sessionStorage.userID}/conversations/${sessionStorage.convId}/messages/${this.messId}/comments/${commentId}`, { headers: { 'Authorization': `${sessionStorage.token}` }});
                this.getMessages();
                // update the comments array, removing the comment
                for (let i = 0; i < this.comments.length; i++) {
                    if (this.comments[i].sender.userId === this.userId) {
                        this.comments.splice(i, 1);
                        this.HaveICommented = false;
                        break;
                    }
                }
            } catch (e) {
                // save and print the error message
                this.errorMsg3 = e.toString();
            }
        },

        // function to show the comments menu, where the user can also comment or uncomment a message
        showComments(object) {
            // get message and comment data
            this.messId = object.message.messageId;
            this.HaveICommented = false;
            this.showCom = !this.showCom;
            // if the message has no comments, comments is an empty array
            if (object.comments == null) {
                this.comments = [];
                return;
            }
            // else, get the comments of the message
            this.comments = object.comments;
            // check if the user has already commented the message. If yes, set it as the last comment of the array
            for (let i = 0; i < this.comments.length; i++) {
                if (this.comments[i].sender.userId === this.userId) {
                    this.HaveICommented = true;
                    const comment = this.comments.splice(i, 1)[0];
                    this.comments.push(comment);
                    break;
                }
            }
        },

        // function to hide the comments menu
        showComments2() {
            this.showCom = !this.showCom;
            this.errorMsg3 = "";
        },

        // function to select the message to reply to. This function is called when the user clicks on the reply button
        selectMessage() {
            this.showReplyTo = true;
            this.selectedMessageId = this.messId;
            this.option = false;
        },

        // function to unselect the message to reply to. This function is called when the user clicks on the remove button
        unselectMessage() {
            this.showReplyTo = false;
            this.selectedMessageId = 0;
        },

        // function to delete the message. This function is called when the user clicks on the delete button
        async deleteMessage() {
            this.errorMsg = "";
            // try the request to delete the message
            try {
                // make the request to delete the message
                this.$axios.delete(`/users/${sessionStorage.userID}/conversations/${sessionStorage.convId}/messages/${this.messId}`, { headers: { 'Authorization': `${sessionStorage.token}`} });
                this.option = false;
                // call the getMessages function to update the messages and their data
                this.getMessages();
            } catch (e) {
                // save and print the error message
                this.errorMsg = e.toString();
            }
        },

        // function to forward the message. This function is called when the user clicks on the forward button
        async ForwardMessage() {
            this.errorMsg2 = "";
            // check if the user has selected at least one conversation or user to forward the message
            if (this.conversationsSelected.length === 0 && this.selectedUsers.length === 0) {
                this.errorMsg2 = "You must select at least one conversation or user!";
                return;
            }
            // prepare the destination array, with objects containing the user or group ids
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
            // clear all variables used to store the selected conversations and users
            object = {
                user: 0,
                group: 0,
            }
            this.conversationsSelected = [];
            this.selectedUsers = [];
            this.conversationsFiltered = [];
            this.filteredUsers = [];
            this.searchText = "";
            // create the input object to send to the server
            const input = {
                destination: destination,
            }
            // try the request to forward the message
            try { 
                // make the request to forward the message
                this.$axios.post(`/users/${sessionStorage.userID}/conversations/${sessionStorage.convId}/messages/${this.messId}`, input, { headers: { 'Authorization': `${sessionStorage.token}` } });
                // close the forward and option menus
                this.showForwardBool = false;
                this.option = false;
            } catch (e) {
                // save and print the error message
                this.errorMsg = e.toString();
            };
        },

        // function to select a group conversation to forward the message. This function is called when the user clicks on the conversation
        selectConversation(conversation) {
            // check if the user has already selected the conversation 
            if (this.conversationsSelected.some(conv => conv.conversation.convId === conversation.conversation.convId)) {
                this.errorMsg2 = "You have already selected this conversation!";
                return;
            }
            // else, add the conversation to the selected conversations array
            this.conversationsSelected.push(conversation);
        },

        //  function to remove a group conversation from the selected conversations array
        removeGroup(index) {
            this.conversationsSelected.splice(index, 1);
            if (this.errorMsg2) {
                this.errorMsg2 = "";
            }
        },

        // function to remove a user from the selected conversations array
        removeMember(index) {
			this.selectedUsers.splice(index, 1);
			if (this.errorMsg2) {
				this.errorMsg2 = "";
			}
		},

        // function to search for users with a specific string in their username 
        async filterUsers() {
            if (this.errorMsg2 !== "It's not necessary that you select yourself!" || this.searchText.length !== 0) {
                this.errorMsg2 = "";
            }
			this.filteredUsers = [];
            // check if the search text is empty
			if (this.searchText.length > 0) {
                // check if the search text is too long or if it contains invalid characters
				if (this.searchText.length > 15 || !this.usernameValidation.test(this.searchText)) {
				this.errorMsg2 = "Invalid username, it can contain only letters and numbers for a maximum of 16 characters.";
				this.filteredUsers = [];
				return;
				}
                // try the request to get the users
				try {
                    // make the request to get the users
					let response = await this.$axios.get(`/users?query=${this.searchText}`, { headers: { 'Authorization': `${sessionStorage.token}` } });
					// check if the response is empty. If yes, set the filteredUsers array to an empty array
                    if (response.data == null) {
					this.filteredUsers = [];
					return;
					}
                    // else, set the filteredUsers array to the response data
					this.filteredUsers = response.data;
				} catch (e) {
                    // save and print the error message
					this.errorMsg2 = e.toString();
					this.filteredUsers = [];
				}
			}
		},

        // function to select a user to forward the message. This function is called when the user clicks on the user
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

        // function to update the group members that are changed in the group info component
        updateGroupMembers() {
            this.members = sessionStorage.members;
        },

        // function to update the group name that is changed in the group info component
        updateGroupname() {
            this.recipientname = sessionStorage.recipientName;
        },

        // function to update the group photo that is changed in the group info component
        updateGroupPhoto() {
            this.recipientPhoto = sessionStorage.recipientPhoto;
        },

        // function to remove the selected photo to send as new message
        removeFile() {
            this.newPhoto = null;
        },

        // function to open the file input to select a new photo to send as new message
        fileInput(){
            this.$refs.file.click();
        },

        // function to handle the file input change event
        handleFileChange(event) {
            // get the file from the input
            const file = event.target.files[0]; 
            // check if the file is null (if the user has not selected a file)
            if (!file) {
            this.errorMsg = "Nessun file selezionato";
            return;
            }
            // check if the file is an image (jpeg or gif)
            if (file.type !== "image/jpeg" && file.type !== "image/jpg" && file.type !== "image/gif") {
            this.errorMsg = "File type not supported, only jpg, jpeg and gif are allowed";
            return;
            }
            // check if the file is too big (5MB)
            if (file.size > 5242880) {
            this.errorMsg = "File size is too big. Max size is 5MB";
            return;
            }
            // assign the file to the newPhoto variable
            this.newPhoto = file;
        },

        // function to go to home page
        goToHome() {
            // if the user is in a group conversation, remove the group data (id and members) from sessionStorage
            if (sessionStorage.groupID) {
                sessionStorage.removeItem("groupID");
                sessionStorage.removeItem("membersOfGroup");
                sessionStorage.removeItem("members");
            }
            // remove the recipient name and photo from sessionStorage, and set the control variable to true
            sessionStorage.removeItem("recipientName");
            sessionStorage.removeItem("recipientPhoto");
            this.control = true;
            this.$router.push("/home");
        },

        // function to get the messages of the conversation
        async getMessages() {
            // check if the user is in the home page, if yes, return
            if (this.control) {
                return;
            }
            this.errorMsg = "";
            this.messages = [];
            // try the request to get the messages
            try{
                // make the request to get the messages
				let response = await this.$axios.get(`/users/${sessionStorage.userID}/conversations/${sessionStorage.convId}`, { headers: { 'Authorization': `${sessionStorage.token}` } });
				// get the messages data from the response
                this.messages = response.data;
			} catch (e) {
                // save and print the error message
				this.errorMsg = e.toString();
			}
        },

        // function to send a new message
        async sendMessage() {
            this.errorMsg = "";
            // create a new FormData object to send a message
            const formData = new FormData();
            // check if the message has a photo. If yes, append it to the formData object
            if (this.newPhoto && (this.newPhoto.type === "image/jpeg" || this.newPhoto.type === "image/jpg")) {
                formData.append('image', this.newPhoto);
            }
            // check if the message has a gif. If yes, append it to the formData object
            if (this.newPhoto && (this.newPhoto.type === "image/gif")) {
                formData.append('gif', this.newPhoto);
            }
            // check if the message has a text. If yes, append it to the formData object
            if (this.newMessage) {
                formData.append('text', this.newMessage);
            }
            // check if the message is a response to other. If yes, append his id to the formData object
            if (this.selectedMessageId != 0) {
                formData.append('messToReplyTo', this.selectedMessageId);
            }
            // check if the message is empty (no text, no photo or no gif). If yes, show an error message
            if (this.newMessage && this.newPhoto && (this.newPhoto.type === "image/gif")) {
                this.errorMsg = "You must write a message or select a photo/gif";
                return;
            }
            // try the request to send the message
            try{
                // make the request to send the message
                let response = await this.$axios.post(`/users/${sessionStorage.userID}/conversations/${sessionStorage.convId}/messages`,
                formData, { headers: { 'Authorization': `${sessionStorage.token}` }});
                // update the messages array with the function getMessages and reset the input fields
                this.newMessage = "";
                this.newPhoto = null;
                this.getMessages();
                this.showReplyTo = false;
                this.selectedMessageId = 0;
            } catch (e) {
                // save and print the error message
                this.errorMsg = e.toString();
            }
        },

        // function to open or close the group info component
        groupInfo() {
            this.showGroupInfo = !this.showGroupInfo;
            // get the group members from sessionStorage. It is important in the case of you are closing the group info component 
            if (sessionStorage.membersOfGroup) {
                this.infoGroupMembers = JSON.parse(sessionStorage.membersOfGroup);
            }
            this.infoGroupMembersQuantity = this.infoGroupMembers.length+1;
        },

        // function to open the message options menu and get the message id
        showOption(object) {
            this.messId = object.message.messageId;
            this.option = !this.option;
        },

        // function to close the message options menu
        showOption2() {
            this.option = !this.option;
        },

        // function to open or close the forward menu, and update the list of group conversations
        showForward() {
            this.showForwardBool = !this.showForwardBool;
            this.getConversations();
            this.option = !this.option;
        },

        // function to get the conversations of the user (important in the case of you want to select a group conversation to forward the message)
        async getConversations() {
			this.errorMsg = "";
			this.conversations = [];
            // try the request to get the conversations
			try{
                // make the request to get the conversations
				let response = await this.$axios.get(`/users/${sessionStorage.userID}/conversations`, { headers: { 'Authorization': `${sessionStorage.token}` } });
				// get the conversations data from the response
                this.conversations = response.data;
                // update the conversationsFiltered array with the conversations that are groups
                for (let i = 0; i < this.conversations.length; i++) {
                    if (this.conversations[i].conversation.group != 0) {
                        this.conversationsFiltered.push(this.conversations[i]);
                    }
                }

			} catch (e) {
                // save and print the error message
				this.errorMsg = e.toString();
			}
		}
    },
    mounted() {
        // every time the component is mounted, get the messages and set the interval to update the messages every 5 seconds
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
    <!-- Chat container -->
    <div class="chat-container">

        
        <!-- Chat header of group conversation -->
        <div v-if="groupId != 0" class="chat-header" @click="groupInfo">
            <!-- back button to go to the home page -->
            <button @click="goToHome" class="back-button"> 
                <svg class="feather"> 
                    <use href="/feather-sprite-v4.29.0.svg#chevron-left" />
                </svg>
            </button>
            <!-- group photo, groupname and members list -->
            <img :src="`data:image/jpg;base64,${recipientPhoto}`" alt="Conversation photo" class="recipient-photo">
            <div class="recipient-info">
                <h2 class="recipient-name">{{ recipientname }}</h2>
                <p class="members-list">{{ members }}</p>
            </div>
        </div>

        
        <!-- Chat header of private conversation -->
        <div v-else class="chat-header">
            <!-- back button to go to the home page -->
            <button @click="goToHome" class="back-button"> 
                <svg class="feather"> 
                    <use href="/feather-sprite-v4.29.0.svg#chevron-left" />
                </svg>
            </button>
            <!-- userphoto and username -->
            <img :src="`data:image/jpg;base64,${recipientPhoto}`" alt="Conversation photo" class="recipient-photo">
            <div class="recipient-info">
                <h2 class="recipient-name">{{ recipientname }}</h2>
            </div>
        </div>
        

        <!-- Chat body that contains all messages -->
        <div class="chat-body">
            <div class="chat-messages">

                <!-- list of messages, that contains all information for single message. It's possible, clicking on a message,
                  to open the option menu of the relative message -->
                <div v-for="object in messages" @click="showOption(object)" :key="object.message.messageId" :class="{'message': true, 'user-message': object.sender.username === username, 'other-message': object.sender.username !== username}">
                    
                    <!-- if the message is a reply to another message, show the reply to message -->
                    <div v-if="object.message.replyId" class="replyToMessage">
                        <!-- check if the reply to message that was sended by the user or by another user (has a different layout) -->
                        <span class="replyToUser" v-if="messages.find(mes => mes.message.messageId === object.message.replyId).sender.userId === userId">Me</span>
                        <span class="replyToUser" v-else>{{ messages.find(mes => mes.message.messageId === object.message.replyId).sender.username }}</span>
                        <!-- check if the reply to message has a photo, text or gif -->
                        <img v-if="messages.find(mes => mes.message.messageId === object.message.replyId).message.photo" :src="`data:image/jpg;base64,${messages.find(mes => mes.message.messageId === object.message.replyId).message.photo}`" alt="Reply to message image" />
                        <span v-if="messages.find(mes => mes.message.messageId === object.message.replyId).message.text">{{ messages.find(mes => mes.message.messageId === object.message.replyId).message.text }}</span>
                        <img v-if="messages.find(mes => mes.message.messageId === object.message.replyId).message.gif" :src="`data:image/gif;base64,${messages.find(mes => mes.message.messageId === object.message.replyId).message.gif}`" alt="Reply to message gif" />
                    </div>
                    
                    <!-- Message header, that contains username, if is a forwarded message and datetime -->
                    <!-- CASE 1: message is sended by other user, in a group conversation (show the name of user) -->
                    <div v-if="groupId != 0 && object.sender.username !== username" class="message-header">
                        <h4>{{ object.sender.username }}</h4>
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
                    <!-- CASE 2: message is sended by other user or the user, in a private conversation (not show the name of user) -->
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

                    <!-- Message body, that contains the message text, photo or gif -->
                    <div v-if="!object.message.text && object.message.photo" class="message-body">
                        <img class="img" :src="`data:image/jpg;base64,${object.message.photo}`" alt="Message image" />
                    </div>
                    <div v-if="object.message.text && object.message.photo" class="message-body">
                        <img class="img" :src="`data:image/jpg;base64,${object.message.photo}`" alt="Message image" />
                        <p>{{ object.message.text }}</p>
                    </div>
                    <div v-if="object.message.gif" class="message-body">
                        <img class="gif" :src="`data:image/gif;base64,${object.message.gif}`" alt="Message gif" />
                    </div>
                    <div v-if="object.message.text && !object.message.photo" class="message-body">
                        <p>{{ object.message.text }}</p>
                    </div>

                    <!-- Message footer, that contains the checkmark and the comments -->
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


        <!-- Chat footer, that contains the input field to send a new message -->
        <div class="chat-footer">
            <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
            <form @submit.prevent="sendMessage">
                <!-- input field to insert a string -->
                <input @keydown.enter.prevent="sendMessage" type="text" v-model="newMessage" placeholder="Type a message" />
                <!-- input field (not showed) to insert a photo or gif -->
                <input type="file" ref="file" accept=".jpg,.jpeg,.gif" @change="handleFileChange" style="display: none;"/>
                <!-- buttons remove the selected photo or gif, unselect message to reply and send finally new message -->
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
    

    <!-- Message options menu, that contains the reply to message, forward message and delete message buttons -->
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


    <!-- Message comments menu, that contains the comment input field and the list of comments -->
    <div v-if="showCom" @click="showComments2" class="fullscreen-container">
        <div @click.stop class="comments-option-container">
            <div class="comments">
                <p>Message comments</p>
                <!-- check if the message has comments -->
                <div v-if="comments.length === 0" class="no-comments">
                    <p>No comments</p>
                </div>
                <!-- if the message has comments, show the list of comments -->
                <div v-for="object in comments" :key="object.commentId" class="comment">
                    <h4 v-if="object.sender.userId !== userId">{{ object.sender.username }}:</h4>
                    <h4 v-if="object.sender.userId === userId">Me:</h4>
                    <p >{{ object.content }}</p>
                </div>
            </div>
            <!-- input field to choose a new comment (an emojis between) -->
            <div class="comment-input">
                <h3>Comment message!</h3>
                <div v-for="emoji in emojis" :key="emoji" class="emoji" @click="commentMessage(emoji)">
                    {{ emoji }}
                </div>
                <!-- trash emoji to delete the comment -->
                <div class="emoji" @click="uncommentMessage()">
                    <svg class="feather"> 
                        <use href="/feather-sprite-v4.29.0.svg#trash-2" />
                    </svg>
                </div>
            </div>
            <ErrorMsg v-if="errorMsg3" :msg="errorMsg3"></ErrorMsg>
        </div>
    </div>


    <!-- Message forward menu, that contains the list of conversations to select to forward the message -->
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
                            <!-- Selected groups, that have the unselect button -->
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
                        <!-- Selected users, that have the unselect button -->
                        <li class="userChoosed" v-for="(conv, index) in selectedUsers" :key="index">
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