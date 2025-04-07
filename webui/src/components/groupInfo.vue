<script>
export default {
    props: {
        show1: Boolean,
        members: Array,
        membersQuantity: Number,
    },
    data: function() {
        return {
            usernameValidation: new RegExp('^[a-z0-9]{1,15}$'),
            errorMsg: "",
            nameOfGroup: sessionStorage.recipientName,
            photo: sessionStorage.recipientPhoto,
            userPhoto: sessionStorage.photo,
            showInputChangeGroupname: false,
            showInputAddUsersToGroup: false,
            newNameOfGroup: "",
            newPhoto: null,
            newGroupMembers: [],
            searchText: "",
            filteredUsers: [],
            localMembers : [],
            localMembersQuantity : 0,
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
    watch: {

        // watch the search text and call the filterUsers function (that searches for users)
        searchText() {
            this.filterUsers();
        },
        // watch the members and membersQuantity props and update the local variables (when the group info is displayed)
        members() {
            this.localMembers = [...this.members];
        },
        membersQuantity() {
            this.localMembersQuantity = this.membersQuantity;
        },
	},
    methods: {
        // function to add users to the group
        async addUsersToGroup(){
            // try the request to add users to the group
			try {
				this.errorMsg = "";
                // make the request to add users to the group
                let response = await this.$axios.put(`/users/${sessionStorage.userID}/groups/${sessionStorage.groupID}`, {
				users: this.newGroupMembers
                }, {headers: {Authorization: `${sessionStorage.token}`}});
                // send the response (array of members of the group, excluding the user that is has just added) to chat view
                this.$emit('update-group-info', response.data);
                // update the local members and quantity with the response
                this.localMembers = response.data;
                this.localMembersQuantity = response.data.length + 1;
                // reset the new group members and close the input for adding users
                this.newGroupMembers = [];
                this.showInputAddUsersToGroup = false;
                // update the list of members (unless the user) in sessionStorage
                let usernames = "";
				for (let i = 0; i < this.localMembers.length; i++) {
					usernames += this.localMembers[i].username + ", ";
				}
				usernames += "me";
				sessionStorage.members = usernames;
                // send the event to update the group members in the chat view
                this.$emit('update-group-members');
            } catch (e) {
                // get and print the error message
                this.errorMsg = e.toString();
            };
		},

        // function to filter users (search for users)
        async filterUsers() {
			this.errorMsg = "";
			this.filteredUsers = [];
            // check if the search text is empty. If it is, doesn't do anything
			if (this.searchText.length > 0) {
                // check if the search text is too long or if it contains invalid characters
				if (this.searchText.length > 15 || !this.usernameValidation.test(this.searchText)) {
				this.errorMsg = "Invalid username, it can contain only letters and numbers for a maximum of 16 characters.";
				this.filteredUsers = [];
				return;
				}
                // try the request to search for users
				try {
                    // make the request to search for users
					let response = await this.$axios.get(`/users?query=${this.searchText}`, { headers: { 'Authorization': `${sessionStorage.token}` } });
					// reset the list of filtered users if the response is null
                    if (response.data == null) {
					this.filteredUsers = [];
					return;
					}
                    // get list of filtered users by the response
					this.filteredUsers = response.data;
				} catch (e) {
                    // check if the error is 404 (user not found) or other error
					if (e.response && e.response.status === 404) {
                        this.errorMsg = "User not found";
                    } else {
                        this.errorMsg = "An error occurred, please try again.";
                    }
                    // reset the list of filtered users
					this.filteredUsers = [];
				}
			}
		},

        // function to select a user from the list of filtered users
        selectUser(user) {
            // check if the user is already in the group or if it is the user that is logged in
            if (this.localMembers.some(member => member.userId === user.userId) || user.userId === Number(sessionStorage.userID)) {
                this.errorMsg = "User is already in the group";
            }
            // check if the user is already in the new group members list. If it is not, add it
			else if (!this.newGroupMembers.includes(user)) {
				this.newGroupMembers.push(user);
			}
            // reset the search text and the list of filtered users
			this.searchText = "";  
			this.filteredUsers = [];  
			},
        
        // function to remove a user from the new group members list
        removeMember(index) {
			this.newGroupMembers.splice(index, 1);
		},

        // function to show or hide the input for adding users to the group
        addUsers() {
            this.showInputAddUsersToGroup = !this.showInputAddUsersToGroup;
        },

        // function to leave the group
        async leaveGroup() {
            // try the request to leave the group
            try {
                // make the request to leave the group
                let response = await this.$axios.delete(`/users/${sessionStorage.userID}/groups/${sessionStorage.groupID}`, {headers: {Authorization: `${sessionStorage.token}`}});
                // clear the part of sessionStorage about the group
                sessionStorage.removeItem("groupID");
                sessionStorage.removeItem("recipientName");
                sessionStorage.removeItem("recipientPhoto");
                sessionStorage.removeItem("membersOfGroup");
                sessionStorage.removeItem("members");
                // go to the home page
                this.$router.push("/home");
            } catch (e) {
                // get and print the error message
                this.errorMsg = e.toString();
            };
        },

        // function to handle the file change (when the user selects a file)
        handleFileChange(event) {
            // get the file from the input
            const file = event.target.files[0]; 
            // check if the file is null (if the user has not selected a file)
            if (!file) {
            this.errorMsg = "Nessun file selezionato";
            return;
            }
            // check if the file is an image (jpeg) and if it is too big
            if (file.type !== "image/jpeg" && file.type !== "image/jpg") {
            this.errorMsg = "File type not supported, only jpg and jpeg are allowed";
            return;
            }
            // check if the file is too big (5MB)
            if (file.size > 5242880) {
            this.errorMsg = "File size is too big. Max size is 5MB";
            return;
            }
            // assign the file to the newPhoto variable
            this.newPhoto = file;
            // go to the function to set the new group photo
            this.setNewGroupPhoto();
        },
        
        // function to open the file selector when the user clicks on the button 
        fileInput(){
            this.$refs.file.click();
        },

        // function to set the new group photo
        async setNewGroupPhoto() {
            // try the request to set the new group photo
            try {
            // create a new FormData object and append the file to it
            const formData = new FormData();
            // append the file to the FormData object
            formData.append('image', this.newPhoto);
            // make the request to set the new group photo
            let response = await this.$axios.put(`/users/${sessionStorage.userID}/groups/${sessionStorage.groupID}/photo`, 
            formData, {headers: {Authorization: `${sessionStorage.token}`}});
            // save the new group photo in sessionStorage and in the component
            sessionStorage.recipientPhoto = response.data.photo;
            this.photo = response.data.photo;
            // send the event to update the group photo in the chat view
            this.$emit('update-group-photo');
            } catch (e) {
            // get and print the error message
                this.errorMsg = e.toString();
            };
        },

        // function to close the modal
        closeModal() {
            this.$emit('close');
        },

        // function to show or hide the input for changing the group name
        async changeGroupName() {
            this.showInputChangeGroupname = !this.showInputChangeGroupname;
        },

        // function to set the new group name
        async setGroupName() {
            // try the request to set the new group name
            try {
                // make the request to set the new group name
                let response = await this.$axios.put(`/users/${sessionStorage.userID}/groups/${sessionStorage.groupID}/name`, {
                groupname: this.newNameOfGroup,
                }, {headers: {Authorization: `${sessionStorage.token}`}});
                // save the new group name in sessionStorage and in the component
                sessionStorage.recipientName = response.data.username;
                this.nameOfGroup = sessionStorage.recipientName;
                // send the event to update the group name in the chat view
                this.$emit('update-groupname');
                // go to the function to close the input for changing the group name
                this.changeGroupName();
            } catch (e) {
                // get and print the error message
                this.errorMsg = e.toString();
            };
        },
        mounted() {
        }
    }
}</script>

<template>
    <!-- Modal to show the group info -->
    <div v-if="show1" class="modal-mask" >
        <div class="modal-wrapper">
            <div class="modal-container">
                <!-- button to close the modal -->
                <button class="exit-button" @click="closeModal">X</button>
                <div class="info-group">
                    <!-- show the group photo -->
                    <img class="info-group-image" :src="`data:image/jpg;base64,${photo}`" alt="user photo"/>
                    <!-- show the group name -->
                    <p class="name-group">{{nameOfGroup}}</p>
                    <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                    <!-- buttons to: change groupname, change group photo, add users to group or leave group -->
                    <div class="button-group">
                        <!-- button to change the group name -->
                        <button @click="changeGroupName">
                            <svg class="feather"> 
                                <use href="/feather-sprite-v4.29.0.svg#edit-3" />
                            </svg>
                            <p>Change groupname
                            </p>
                        </button>
                        <!-- button to change the group photo -->
                        <input type="file" ref="file" accept=".jpg,.jpeg" @change="handleFileChange" style="display: none;"/>
                        <button @click="fileInput">
                            <svg class="feather"> 
                                <use href="/feather-sprite-v4.29.0.svg#edit" />
                            </svg>
                            <p>Change photo
                            </p>
                        </button>
                        <!-- button to add users to the group -->
                        <button @click="addUsers">
                            <svg class="feather"> 
                                <use href="/feather-sprite-v4.29.0.svg#user-plus" />
                            </svg>
                            <p>Add users to group
                            </p>
                        </button>
                        <!-- button to leave the group -->
                        <button @click="leaveGroup">
                            <svg class="feather"> 
                                <use href="/feather-sprite-v4.29.0.svg#log-out" />
                            </svg>
                            <p>Leave group
                            </p>
                        </button>
                    </div>
                    <!-- show the number of members in the group and the list of members -->
                    <div class="members-wrapper">
                        <p class="members-count">{{ localMembersQuantity }} members</p>
                        <div class="members-box">
                            <div class="members-bands">
                                <div v-for="member in localMembers" :key="member.userId" class="member">
                                    <img :src="`data:image/jpg;base64,${member.userPhoto}`" alt="User photo"/>
                                    <h4>{{ member.username }}</h4>
                                </div>
                                <div class="member">
                                    <img :src="`data:image/jpg;base64,${userPhoto}`" alt="User photo"/>
                                    <h4>me</h4>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>        
    </div>

    <!-- Modal to show the input for changing the group name -->
    <div v-if="showInputChangeGroupname" class="modal-mask2" @click="changeGroupName">
        <div class="modal-wrapper2">
            <div class="modal-container2" @click.stop>
                <div class="modal-header2">
                    <form @submit.prevent="setGroupName">
                        <h3>Change groupname</h3>
                        <!-- input to set the new group name -->
                        <input type="text" v-model="newNameOfGroup" placeholder="Enter the new groupname" />
                        <!-- button to set the new group name -->
                        <button type="submit">confirm</button>
                    </form>
                </div>
                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
            </div>
        </div>
    </div>

    <!-- Modal to show the input for adding users to the group -->
    <div v-if="showInputAddUsersToGroup" class="modal-mask3" @click="addUsers">
        <div class="modal-wrapper3">
            <div class="modal-container3" @click.stop>
                <div class="modal-header3">
                    <h3>Choose users</h3>
                    <form @submit.prevent="addUsersToGroup">
                        <!-- Select the members of the group -->
                        <div class="new-group-members">
                            <label for="newGroupMembers">Select a new group member</label>
                            <div class="input-members">
                                <input type="text" class="form-control" v-model="searchText" placeholder="Enter a new member's name" />
                            </div>
                            <!-- Search and print all users who have username that start with the text in the input space -->
                            <div class="search-results">
                                <div v-for="user in filteredUsers" :key="user.userId" @click="selectUser(user)" class="user">
                                    <p>{{ user.username }}</p>
                                </div>
                            </div>	
                        </div>
                        <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                        <ul>
                            <!-- List of all selected new members of the group. If the list is empty, it doesn't show anything 
                            every member has a button to remove it from the list -->
                            <label v-if="newGroupMembers.length > 0">New group members list</label>
                            <li v-for="(member, index) in newGroupMembers" :key="index">
                                {{ member.username }} <button @click.prevent="removeMember(index)">x</button>
                            </li>
                        </ul>
                        <hr class="separator">
                        <button class="buttonSubmit" type="submit">Add</button>
                    </form>
                </div>
            </div>
        </div>
    </div>
</template>

<style>
.modal-mask {
    position: fixed;
    z-index: 100000;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, .5);
    display: table;
    transition: opacity .3s;
}
.modal-mask2 {
    position: fixed;
    z-index: 100001;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, .5);
    display: table;
    transition: opacity .3s;
}

.modal-mask3 {
    position: fixed;
    z-index: 100002;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, .5);
    display: table;
    transition: opacity .3s;
}

.modal-wrapper {
    display: table-cell;
    vertical-align: middle;
}

.modal-wrapper2 {
    display: table-cell;
    vertical-align: middle;
}

.modal-wrapper3 {
    display: table-cell;
    vertical-align: middle;
}

.modal-container {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 800px;
    margin: 0 auto;
    padding: 20px;
    background-color: white;
}

.modal-container2 {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 800px;
    margin: 0 auto;
    padding: 20px;
    background-color: white;
    border-radius: 5px;
}

.modal-container3 {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 800px;
    margin: 0 auto;
    padding: 20px;
    background-color: white;
    border-radius: 5px;
}

.modal-header2 {
    display: flex;
    justify-content: center;
    align-items: center;
}

.modal-header3 {
    display: flex; 
    flex-direction: column; 
    gap: 10px; 
    width: 100%;
}

.modal-header2 form {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
}
.modal-header2 input {
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 5px;
}

.modal-header2 button {
    padding: 10px 20px;
    border: none;
    border-radius: 5px;
    background-color: #007bff;
    color: white;
    cursor: pointer;
}

.info-group {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
}

.name-group {
    font-size: 30px;
    font-weight: bold;

    
}

.info-group-image {
    width: 160px;
    height: 160px;
    border-radius: 50%;
    object-fit: cover;
    border: 5px solid #007bff;
}

.button-group {
    display: flex;
    gap: 10px;
}

.button-group button {
    padding: 10px 20px;
    width: 183px;
    border: none;
    border-radius: 5px;
    background-color: #007bff;
    color: white;
    cursor: pointer;
    font-size: 13px;
    display: flex;
    flex-direction: column; 
    align-items: center; 
    justify-content: center;
    gap: 5px;
}

.button-group button svg {
    width: 24px;  
    height: 24px;
    display: block;
}

.button-group button p {
    margin: 0;
    padding: 0;
    font-size: 12px;
    line-height: 1;
}

.members-wrapper {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    width: 100%;
}

.members-count {
    font-size: 15px;
    color: #555;
    margin-bottom: 5px;
    padding-left: 5px;
    padding-top: 10px;
}

.members-box {
    background: white;
    border-radius: 10px;
    padding: 15px;
    width: 100%;
    box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
}

.members-bands {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.member {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px;
    border-radius: 8px;
    background: #f0f0f0;
}

.member h4 {
    margin: 0;
    padding: 0;
    font-size: 20px;
}

.member img {
    width: 50px;
    height: 50px;
    border-radius: 50%;
    object-fit: cover;
    border: 3px solid #007bff;
}

.exit-button {
    position: absolute;
    color: white;
    top: 10px;
    right: 10px;
    border: none;
    background-color: red;
    font-size: 15px;
    cursor: pointer;
    border-radius: 50%;
    padding: 5px 11.5px;

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
</style>