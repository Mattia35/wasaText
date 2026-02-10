# WASAText

## Academic Project Information

**WASAText** is a university project developed for the **Web and Software Architecture (WASA)** exam.

**Project:** WASA Project – *“WASAText”*  
**Version:** 1  

As part of the exam, the project focuses on:
- defining APIs using the **OpenAPI** standard  
- designing and developing the **backend** in **Go**  
- designing and developing the **frontend** in **JavaScript**  
- creating **Docker container images** for deployment  

This repository contains a complete implementation of the WASAText system.

---

## Introduction

**WASAText** allows users to connect with friends effortlessly through a web-based messaging application.  
Users can send and receive messages, both in private conversations and group chats, directly from their PC.

The application supports:
- text and GIF messages  
- private and group conversations  
- message reactions and delivery/read indicators  

---

## Functional Overview

The user is presented with a list of conversations (private or group), sorted in reverse chronological order.

Each conversation shows:
- username or group name  
- user or group profile photo  
- date and time of the latest message  
- message preview or an icon for photo messages  

Users can:
- start new conversations with any existing WASAText user  
- search users by username  
- create group chats and add other users  
- leave groups at any time  

Inside a conversation, users can:
- view messages in reverse chronological order  
- see timestamps and sender information  
- send, reply to, forward, and delete messages  
- react to messages with emoticons and remove reactions  

### Message Status Indicators

- **One checkmark (✓)**: the message has been received  
- **Two checkmarks (✓✓)**: the message has been read  

---

## Simplified Login

WASAText uses a **simplified authentication mechanism** designed specifically for this academic project.

- Users log in by entering a **username only**
- No passwords are required
- If the username already exists, the user is logged in
- If the username is new, the user is automatically registered and logged in

The backend returns a **user identifier** that must be included as a Bearer token in the `Authorization` header for all API requests.

⚠️ **Security Notice**  
This authentication system is intentionally simplified and **not suitable for real-world applications**.  
It is used only to focus on software architecture and API design concepts.

---

## Requirements

To run WASAText locally, the following are required:

- **Docker**
- A modern web browser (Chrome, Firefox, Edge, etc.)

No additional dependencies are required for end users.

---

## Running the Application

### Build Docker Images

From the project root directory, build the backend and frontend images:

```bash
docker build -t wasa-text-backend:latest -f Dockerfile.backend .
docker build -t wasa-text-frontend:latest -f Dockerfile.frontend .
```

Run the Containers

Run the backend and frontend in two separate terminal windows:

```bash
docker run -it --rm -p 3000:3000 wasa-text-backend:latest
docker run -it --rm -p 8080:80 wasa-text-frontend:latest
```
Both containers must be running for the application to work correctly.

Accessing WASAText

Once the application is running, open your browser and go to:
```code
http://localhost:8080
```
Enter a username to log in and start using WASAText.
Notes and Limitations
	•	This project is intended for educational purposes only
	•	The login system is intentionally simplified
	•	No password recovery or identity verification is implemented
	•	Data persistence depends on the local backend database


