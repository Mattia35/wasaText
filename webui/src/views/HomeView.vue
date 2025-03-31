<script>
export default {
	emits: ['to-home', 
			'login-success', 
			'update-username', 
			'close', 
			'update-photo', 
			'update-groupname', 
			'update-group-photo',
			'update-group-members'],
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
		searchText() {
			this.filterUsers();
		}
	},
	methods: {
		goToChat(conversation) {
			
			if (conversation.conversation.group != 0) {
				sessionStorage.groupID = conversation.conversation.group;
				sessionStorage.recipientPhoto = conversation.group.groupPhoto;
				sessionStorage.recipientName = conversation.group.username;
				sessionStorage.membersOfGroup = JSON.stringify(conversation.groupUsers);
				let usernames = "";
				for (let i = 0; i < conversation.groupUsers.length; i++) {
					usernames += conversation.groupUsers[i].username + ", ";
				}
				usernames += "me";
				sessionStorage.members = usernames;
			} else {
				sessionStorage.members = "";
				sessionStorage.groupID = 0;
				sessionStorage.chatUserID = conversation.user.userId;
				sessionStorage.recipientName = conversation.user.username;
				sessionStorage.recipientPhoto = conversation.user.userPhoto;
			}
			sessionStorage.convId = conversation.conversation.convId;
			this.$router.push(`/chat/${conversation.conversation.convId}`);
		},
		async CreatePrivateChat() {
			// check if the user is already in a conversation with the selected user
			for (let i = 0; i < this.conversations.length; i++) {
				if (this.conversations[i].conversation.group == 0 && this.conversations[i].user.userId == this.groupMembers[0].userId) {
					this.errorMsg2 = "A conversation with this user already exists!";
					return;
				}
				sessionStorage.recipientName = this.groupMembers[0].username;
				sessionStorage.recipientPhoto = this.groupMembers[0].userPhoto;
				sessionStorage.members = "";
				sessionStorage.groupID = 0;
				sessionStorage.chatUserID = this.groupMembers[0].userId;
			}
			this.$router.push(`/temporary-chat/${this.groupMembers[0].username}`);
		},
		async getConversations() {
			this.errorMsg = "";
			this.conversations = [];
			try{
				let response = await this.$axios.get(`/users/${sessionStorage.userID}/conversations`, { headers: { 'Authorization': `${sessionStorage.token}` } });
				this.conversations = response.data;
			} catch (e) {
				this.errorMsg = e.toString();
				document.getElementsByTagName("input")[0].style.outline = "auto";
                document.getElementsByTagName("input")[0].style.outlineColor = "red";
			}
		},
		async filterUsers() {
			this.errorMsg2 = "";
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
					document.getElementsByTagName("input")[0].style.outline = "auto";
                	document.getElementsByTagName("input")[0].style.outlineColor = "red";
					this.filteredUsers = [];
				}
			}
		},

		selectUser(user) {
			if (user.userId === Number(sessionStorage.userID)) {
				this.errorMsg2 = "It's not necessary that you select yourself!";
			}
			else if (!this.groupMembers.some(member => member.userId === user.userId)) {
				this.groupMembers.push(user);
			}
			this.searchText = "";  
			this.filteredUsers = [];  
			},

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

		showCreateGroup(){
			this.CreationConv();
			this.showCreationGroup = !this.showCreationGroup;
		},

		showCreatePrivateChat(){
			this.CreationConv();
			this.showCreationPrivateChat = !this.showCreationPrivateChat;
		},

		async CreateGroup(){
			try {
				this.errorMsg2 = "";
                // Controlla che il groupname sia valido
                if (this.groupName.trim() === "" || this.groupMembers.length === 0) throw "Il nome del gruppo e almeno un membro sono richiesti!";
                // Effettua la richiesta di creazione gruppo al server
                let response = await this.$axios.post(`/users/${sessionStorage.userID}/groups`, {
                groupname: this.groupName,
				users: this.groupMembers
                }, {headers: {Authorization: `${sessionStorage.token}`}});
				this.CreationConv();
				this.getConversations();

            } catch (e) {
                this.errorMsg2 = e.toString();
                document.getElementsByTagName("input")[0].style.outline = "auto";
                document.getElementsByTagName("input")[0].style.outlineColor = "red";
            };
		},

		removeMember(index) {
			this.groupMembers.splice(index, 1);
			if (this.errorMsg2) {
				this.errorMsg2 = "";
			}
		},
	},
	mounted() {
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
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Home page</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="exportList">
						Export
					</button>
				</div>
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="CreationConv">New</button>
				</div>
			</div>
		</div>
		<ul class="conversation-list">
			<li v-for="conversation in conversations" :key="conversation.conversation.convId" class="conversation-item" @click="goToChat(conversation)">
				<div class="conversation">
					<img v-if="conversation.conversation.group != 0" :src="`data:image/jpg;base64,${conversation.group.groupPhoto}`" alt="Conversation photo" >
					<img v-if="conversation.conversation.group == 0" :src="`data:image/jpg;base64,${conversation.user.userPhoto}`" alt="Conversation photo" >
					<div class="conversation-name-and-last-message">
						<p class="conversation-name" v-if="conversation.conversation.group != 0">{{ conversation.group.username }}</p>
						<p class="conversation-name" v-if="conversation.conversation.group == 0">{{ conversation.user.username }}</p>
						<p v-if="conversation.message.photo" class="last-message">
							<p v-if="conversation.senderUser.username === username" class="sender-last-message"><strong>me: </strong></p>
							<p v-if="conversation.senderUser.username !== username" class="sender-last-message"><strong>{{ conversation.senderUser.username }}: </strong></p>
						
							<svg class="feather"> 
								<use href="/feather-sprite-v4.29.0.svg#image" />
							</svg>
							image
						</p>
						<p v-else-if="conversation.message.text && conversation.senderUser.username === username" class="last-message"><strong>me: </strong>{{ conversation.message.text }}</p>
						<p v-else-if="conversation.message.text && conversation.senderUser.username !== username" class="last-message"><strong>{{ conversation.senderUser.username }}: </strong>{{ conversation.message.text }}</p>
						<p v-else class="last-message">
							<p v-if="conversation.senderUser.username === username" class="sender-last-message"><strong>me: </strong></p>
							<p v-if="conversation.senderUser.username !== username" class="sender-last-message"><strong>{{ conversation.senderUser.username }}: </strong></p>
							<svg class="feather"> 
								<use href="/feather-sprite-v4.29.0.svg#image" />
							</svg>
							gif
						</p>
					</div>
				</div>
				<p class="conversation-datetime">{{ conversation.dateTime }}</p>
			</li>
		</ul>
	</div>
	
	<div v-if="showCreationConv || showCreationGroup || showCreationPrivateChat" class="fullscreen-container" @click="CreationConv">
		<div v-if="showCreationConv" class="dropdownConv" @click.stop>
			<!-- Button for the creation of a new group -->
			<button type="button" class="group-or-conv" @click="showCreateGroup">New group</button>
			<!-- Button for the creation of a new private conversation -->
			<button type="button" class="group-or-conv" @click="showCreatePrivateChat">New private chat</button>
		</div>
	</div>
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
				<!-- List of all selected new members of the group -->
				<label v-if="groupMembers.length > 0">Group members list</label>
				<li v-for="(member, index) in groupMembers" :key="index">
					{{ member.username }} <button @click.prevent="removeMember(index)">x</button>
				</li>
			</ul>
			<hr class="separator">
			<button class="buttonSubmit" type="submit">Create</button>
		</form>
	</div>
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
				<!-- User who is selected -->
				<label v-if="groupMembers.length > 0">You have selected an user. If you have selected a wrong person, unselect him and search for the right one!</label>
				<li v-for="(member, index) in groupMembers" :key="index">
					{{ member.username }} <button @click.prevent="removeMember(index)">x</button>
				</li>
			</ul>
			<hr class="separator">
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

}

.sender-last-message, .last-message {
  display: inline-block;
  margin-right: 5px;
}

</style>
