<script>
export default {
	emits: ['to-home', 
			'login-success', 
			'update-username', 
			'close', 
			'update-photo', 
			'update-groupname', 
			'update-group-photo',
			'update-group-members',
			'update-group-info'],
	data: function() {
		return {
			username: sessionStorage.username,
			errorMsg: "",
			errorMsg2: "",
			showCreationConv: false,
			showCreationGroup: false,
			showCreationPrivateChat: false,
			groupName: "",
			newMember: "",
      		groupMembers: [],
			filteredUsers: [],
			searchText: "",
			usernameValidation: new RegExp('^[a-z0-9]{1,15}$'),
			conversations: [],
			intervalId: null,
		}
	},
	watch: {
		// watch the search text and call the filterUsers function when it changes
		searchText() {
			this.filterUsers();
		}
	},
	methods: {

		// function to go to the chat page
		goToChat(conversation) {
			// check if the conversation is a group or a private chat
			if (conversation.conversation.group != 0) {
				// if it is a group, save the group id, the group name, the group photo and the members of the group in sessionStorage
				sessionStorage.groupID = conversation.conversation.group;
				sessionStorage.recipientPhoto = conversation.group.groupPhoto;
				sessionStorage.recipientName = conversation.group.username;
				sessionStorage.membersOfGroup = JSON.stringify(conversation.groupUsers);
				// create and save the string with the members of the group
				let usernames = "";
				for (let i = 0; i < conversation.groupUsers.length; i++) {
					usernames += conversation.groupUsers[i].username + ", ";
				}
				usernames += "me";
				sessionStorage.members = usernames;
			} else {
				// if it is a private chat, save the user id, the user name and the user photo in sessionStorage
				sessionStorage.members = "";
				sessionStorage.groupID = 0;
				sessionStorage.chatUserID = conversation.user.userId;
				sessionStorage.recipientName = conversation.user.username;
				sessionStorage.recipientPhoto = conversation.user.userPhoto;
			}
			// save the conversation id in sessionStorage and go to the chat page
			sessionStorage.convId = conversation.conversation.convId;
			this.$router.push(`/chat/${conversation.conversation.convId}`);
		},

		// function to go to the temporary chat page, for the creation of a new private chat
		async CreatePrivateChat() {
			// check if the user is already in a conversation with the selected user. If yes, show an error message
			for (let i = 0; i < this.conversations.length; i++) {
				if (this.conversations[i].conversation.group == 0 && this.conversations[i].user.userId == this.groupMembers[0].userId) {
					this.errorMsg2 = "A conversation with this user already exists!";
					return;
				}
			}
			// save the user id, the user name and the user photo in sessionStorage
			sessionStorage.recipientName = this.groupMembers[0].username;
			sessionStorage.recipientPhoto = this.groupMembers[0].userPhoto;
			sessionStorage.members = "";
			sessionStorage.groupID = 0;
			sessionStorage.chatUserID = this.groupMembers[0].userId;
			// go to the temporary chat page
			this.$router.push(`/temporary-chat/${this.groupMembers[0].username}`);
		},

		// function to get the conversations from the server
		async getConversations() {
			this.errorMsg = "";
			this.conversations = [];
			// try the request
			try{
				// make the request to get the conversations
				let response = await this.$axios.get(`/users/${sessionStorage.userID}/conversations`, { headers: { 'Authorization': `${sessionStorage.token}` } });
				// save the conversations in the data variable
				this.conversations = response.data;
			} catch (e) {
				// save and print the error message
				this.errorMsg = e.toString();
			}
		},

		// function to filter the users based on the search text
		async filterUsers() {
			this.errorMsg2 = "";
			this.filteredUsers = [];
			// check if the text in input is valid. If not, show an error message
			if (this.searchText.length > 0) {
				if (this.searchText.length > 15 || !this.usernameValidation.test(this.searchText)) {
				this.errorMsg2 = "Invalid username, it can contain only letters and numbers for a maximum of 16 characters.";
				this.filteredUsers = [];
				return;
				}
				// try the request to get the users
				try {
					// make the request to get the users
					let response = await this.$axios.get(`/users?query=${this.searchText}`, { headers: { 'Authorization': `${sessionStorage.token}` } });
					// check if the response is empty. If yes, get in output an empty array
					if (response.data == null) {
					this.filteredUsers = [];
					return;
					}
					// else, save the users in the data variable
					this.filteredUsers = response.data;
				} catch (e) {
					// save and print the error message
					if (e.response && e.response.status === 404) {
						this.errorMsg2 = "User not found";
					} else {
						this.errorMsg2 = "An error occurred, please try again.";
					}
				}
			}
		},

		// function to select a user from the search results
		selectUser(user) {
			// check if the user is the same as the logged user. If yes, show an error message
			if (user.userId === Number(sessionStorage.userID)) {
				this.errorMsg2 = "It's not necessary that you select yourself!";
				this.filteredUsers = [];
				return;
			}
			// check if the user is already in the list of members of new group. If not, add him to the group members
			else if (!this.groupMembers.some(member => member.userId === user.userId)) {
				this.groupMembers.push(user);
			}
			// set searchText to empty string and filteredUsers to empty array
			this.searchText = "";  
			this.filteredUsers = [];
		},

		// function to show or hide the menu for the creation of a new group or a new private chat
		CreationConv(){
			if (this.showCreationGroup) {
				this.showCreationGroup = false;
				this.groupMembers = [];
				this.groupName = "";
				this.searchText = "";
				return;
			} 
			if (this.showCreationPrivateChat) {
				this.showCreationPrivateChat = false;
				this.groupMembers = [];
				this.searchText = "";
				return;
			}
			this.showCreationConv = !this.showCreationConv;
		},

		// function to show or hide the menu for the creation of a new group
		showCreateGroup(){
			this.CreationConv();
			this.showCreationGroup = !this.showCreationGroup;
		},

		// function to show or hide the menu for the creation of a new private chat
		showCreatePrivateChat(){
			this.CreationConv();
			this.showCreationPrivateChat = !this.showCreationPrivateChat;
		},

		// function to create a new group
		async CreateGroup(){
			// try the request to creation of a new group
			try {
				this.errorMsg2 = "";
                // check if the group name is valid
                if (this.groupName.trim() === "" || this.groupMembers.length === 0) throw "Il nome del gruppo e almeno un membro sono richiesti!";
                // make the request to create a new group
                let response = await this.$axios.post(`/users/${sessionStorage.userID}/groups`, {
                groupname: this.groupName,
				users: this.groupMembers
                }, {headers: {Authorization: `${sessionStorage.token}`}});
				// close the menu for the creation of a new group, and update the list of user conversations
				this.CreationConv();
				this.getConversations();

            } catch (e) {
				// save and print the error message
                this.errorMsg2 = e.toString();
            };
		},

		// function to remove a member from the new group members list
		removeMember(index) {
			this.groupMembers.splice(index, 1);
			if (this.errorMsg2) {
				this.errorMsg2 = "";
			}
		},
	},
	mounted() {
		// get the conversations from the server when the component is mounted, and set the interval to update the conversations every 5 seconds
		this.getConversations();
		this.intervalId = setInterval(async () => {
			clearInterval(this.intervalId);
			await this.getConversations();
			this.intervalId = setInterval(this.getConversations, 5000);
		}, 5000);
	}
}
</script>

<template>
	<div>
		<!-- Home page vue -->
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Home page</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<!-- Button to create a new conversation. It open a menu where the user choose which conv he want to create (group or priv) -->
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="CreationConv">New</button>
				</div>
			</div>
		</div>
		<!-- List of all conversations -->
		<ul class="conversation-list">
			<!-- conversation object -->
			<li v-for="conversation in conversations" :key="conversation.conversation.convId" class="conversation-item" @click="goToChat(conversation)">
				<div class="conversation">
					<!-- conversation photo. If it is a group, show the group photo, otherwise show the user photo -->
					<img v-if="conversation.conversation.group != 0" :src="`data:image/jpg;base64,${conversation.group.groupPhoto}`" alt="Conversation photo" >
					<img v-if="conversation.conversation.group == 0" :src="`data:image/jpg;base64,${conversation.user.userPhoto}`" alt="Conversation photo" >
					<!-- conversation name and last message -->
					<div class="conversation-name-and-last-message">
						<!-- conversation name. If it is a group, show the group name, otherwise show the user name -->
						<p class="conversation-name" v-if="conversation.conversation.group != 0">{{ conversation.group.username }}</p>
						<p class="conversation-name" v-if="conversation.conversation.group == 0">{{ conversation.user.username }}</p>
						<!-- if message is a photo, show the sender username, the photo icon and the text "image" --> 
						<div v-if="conversation.message.photo && conversation.message.messageId != 0" class="last-message">
							<p v-if="conversation.senderUser.username === username" class="sender-last-message"><strong>me: </strong></p>
							<p v-if="conversation.senderUser.username !== username" class="sender-last-message"><strong>{{ conversation.senderUser.username }}: </strong></p>
						
							<svg class="feather"> 
								<use href="/feather-sprite-v4.29.0.svg#image" />
							</svg>
							image
						</div>
						<!-- if message is a text and the sender is the logged user, show "Me" and the text of the message -->
						<div v-else-if="conversation.message.text && conversation.senderUser.username === username && conversation.message.messageId != 0" class="last-message"><strong>me: </strong>{{ conversation.message.text }}</div>
						<!-- if message is a text and the sender is not the logged user, show the sender username and the text of the message -->
						<div v-else-if="conversation.message.text && conversation.senderUser.username !== username && conversation.message.messageId != 0" class="last-message"><strong>{{ conversation.senderUser.username }}: </strong>{{ conversation.message.text }}</div>
						<!-- if message is a gif, show the sender username, the gif icon and the text "gif" -->
						<div v-else-if="conversation.message.messageId != 0" class="last-message">
							<p v-if="conversation.senderUser.username === username" class="sender-last-message"><strong>me: </strong></p>
							<p v-if="conversation.senderUser.username !== username" class="sender-last-message"><strong>{{ conversation.senderUser.username }}: </strong></p>
							<svg class="feather"> 
								<use href="/feather-sprite-v4.29.0.svg#image" />
							</svg>
							gif
						</div>
					</div>
				</div>
				<!-- date and time of the last message -->
				<p v-if="conversation.message.messageId != 0" class="conversation-datetime">{{ conversation.dateTime }}</p>
			</li>
		</ul>
	</div>
	

	<!-- creation of a new conversation menu -->
	<div v-if="showCreationConv || showCreationGroup || showCreationPrivateChat" class="fullscreen-container" @click="CreationConv">
		<div v-if="showCreationConv" class="dropdownConv" @click.stop>
			<!-- Button for the creation of a new group -->
			<button type="button" class="group-or-conv" @click="showCreateGroup">New group</button>
			<!-- Button for the creation of a new private conversation -->
			<button type="button" class="group-or-conv" @click="showCreatePrivateChat">New private chat</button>
		</div>
	</div>


	<!-- menu for the creation of a new group -->
	<div v-if="showCreationGroup" class="menu-Creation">
		<h3>New group info</h3>
		<form @submit.prevent="CreateGroup">
			<!-- Select the name of the group -->
			<div class="group-name">
				<label for="groupName">Group name</label>
				<input type="text" class="form-control" id="groupName" v-model="groupName" placeholder="Enter the group name" required>
			</div>
			<!-- Select the members of the group -->
			<div class="group-members">
				<label for="groupMembers">Select a new group member</label>
				<div class="input-members">
					<input type="text" class="form-control" v-model="searchText" placeholder="Enter a member's name" />
				</div>
				 <!-- Search and print all users who have username that start with the text in the input space -->
				 <div class="search-results">
					<div v-for="user in filteredUsers" :key="user.userId" @click="selectUser(user)" class="user">
						<p>{{ user.username }}</p>
					</div>
				</div>	
			</div>
			<ErrorMsg v-if="errorMsg2" :msg="errorMsg2"></ErrorMsg>
			<ul>
				<!-- List of all selected new members of the group. If the user is selected, show the name of the user and a button to remove him from the list -->
				<label v-if="groupMembers.length > 0">Group members list</label>
				<li v-for="(member, index) in groupMembers" :key="index">
					{{ member.username }} <button @click.prevent="removeMember(index)">x</button>
				</li>
			</ul>
			<hr class="separator">
			<!-- Button to create the new group -->
			<button class="buttonSubmit" type="submit">Create</button>
		</form>
	</div>


	<!-- menu for the creation of a new private chat -->
	<div v-if="showCreationPrivateChat" class="menu-Creation">
		<h3>Create a new private chat</h3>
		<form @submit.prevent="CreatePrivateChat">
			<!-- Select the user -->
			<div v-if="groupMembers.length != 1" class="select-user">
				<label for="groupMembers">Select the user</label>
				<div class="input-members">
					<input type="text" class="form-control" v-model="searchText" placeholder="Enter a username" required/>
				</div>
				 <!-- Search and print all users who have username that start with the text in the input space -->
				 <div class="search-results">
					<div v-for="user in filteredUsers" :key="user.userId" @click="selectUser(user)" class="user">
						<p>{{ user.username }}</p>
					</div>
				</div>	
			</div>
			<ErrorMsg v-if="errorMsg2" :msg="errorMsg2"></ErrorMsg>
			<ul>
				<!-- User who is selected. If the user is selected, show the name of the user and a button to remove him from the list -->
				<label v-if="groupMembers.length > 0">You have selected an user. If you have selected a wrong person, unselect him and search for the right one!</label>
				<li v-for="(member, index) in groupMembers" :key="index">
					{{ member.username }} <button @click.prevent="removeMember(index)">x</button>
				</li>
			</ul>
			<hr class="separator">
			<!-- Button to create the new private chat -->
			<button class="buttonSubmit" type="submit">Create</button>
		</form>
	</div>
</template>

<style>
.d-flex {
  position: sticky;
  top: 0;
  background-color: white;
  z-index: 999;
  padding: 10px 20px;
}

.fullscreen-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  display: flex; 
  justify-content: center; 
  align-items: center;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 9999;
}
.dropdownConv {
  position: absolute;
  top: 6.7%; 
  right: 1.8%;
  background-color: #fff; 
  box-shadow: 0px 4px 15px rgba(0, 0, 0, 0.1); 
  border-radius: 8px;
  padding: 0; 
  z-index: 9999; 
  width: 250px; 
  display: flex;
  flex-direction: column; 
  gap: 0px; 
}

.group-or-conv {
  background-color: #ffffff;
  color: black; 
  padding: 15px 20px; 
  font-size: 16px;
  cursor: pointer; 
  text-align: center;
  border: none;
  border-radius: 5px; 
  transition: background-color 0.3s ease, color 0.3s ease;
}

.group-or-conv:hover {
  background-color: #007bff;
  color: white; 
}

.menu-Creation {
  position: absolute; 
  top: 10%; 
  left: 10%; 
  background-color: #fff; 
  box-shadow: 0px 4px 15px rgba(0, 0, 0, 0.1); 
  border-radius: 8px; 
  padding: 20px; 
  z-index: 9999; 
  width: 80%; 

}

.menu-Creation form{
  display: flex; 
  flex-direction: column; 
  gap: 10px; 
}

.input-members {
    display: flex;
    align-items: center;  
    gap: 10px;            
}

.input-members input {
    flex-grow: 1;          
}

ul {
  list-style-type: none; 
  padding: 0; 
  margin: 0;
}

li {
  display: flex;
  align-items: center; 
  justify-content: space-between; 
  background-color: #f8f9fa; 
  padding: 4px 4px; 
  margin-bottom: 5px; 
  border-radius: 5px; 
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1); 
  font-size: 16px; 
  transition: background-color 0.3s ease; 
}


li:hover {
  background-color: #e2e6ea; 
}


li button {
  padding: 5px 10px; 
  font-size: 16px; 
  border: none;
  background-color: #ff0000;
  color: white;
  cursor: pointer;
  border-radius: 5px; 
  flex-shrink: 0; 
}

.buttonSubmit {
  padding: 10px 20px; 
  border: none; 
  border-radius: 5px; 
  background-color: #007bff; 
  color: white; 
  cursor: pointer; 
  font-size: 16px; 
}

.conversation-list {
  list-style: none;
  padding: 0;
  margin: 0;
  height: calc(100vh - 100px);
  overflow-y: auto;
  overflow-x: hidden;
}

.conversation-item {
  padding: 10px;
  border-bottom: 1px solid #ddd;
  cursor: pointer;
  transition: background 0.3s;
  position: relative;
}

.conversation-item:hover {
  background: #f5f5f5;
}

.conversation {
  display: flex;
  align-items: center;
  gap: 10px;
}

.conversation-datetime {
  position: absolute;
  top: 0; 
  right: 0; 
  font-size: 12px;
  color: #666;
  margin: 5px;
  padding: 3px 5px;
  border-radius: 3px;
}

.conversation img {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #007bff;
}

.conversation-name-and-last-message {
  display: flex;
  flex-direction: column;
}

.conversation-name {
  font-weight: bold;
  font-size: 16px;
  margin: 0;
}

.last-message {
  display: flex;
  align-items: center;
  font-size: 14px;
  color: #666;
  margin: 0;
  max-height: 20px;

}

.sender-last-message, .last-message {
  display: inline-block;
  margin-right: 5px;
}

</style>
