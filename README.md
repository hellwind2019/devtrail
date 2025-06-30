# 🛠️ DevTrail — Gamified Coding Progress Tracker

A lightweight web application written in **Go** that allows developers to register, log in, and track their daily coding progress.  
Inspired by Boot.dev and habit trackers, but focused on personal projects and developer productivity.

## 📌 Features

- ✅ User registration via a simple HTML form
- ✅ User login with session-based authentication using `gorilla/sessions`
- ✅ Protected dashboard accessible only to authenticated users
- ✅ Logout functionality with session cleanup
- ✅ Clean, modular project structure
- ✅ Data stored in a simple local text file (`users.txt`) for quick prototyping

## 🖥️ Tech Stack

- **Go**
- **HTML (templates)**
- **Gorilla Sessions**


## Setup .env

Variables :
```
SESSION_KEY= any value for session enctyprion
GH_BASIC_CLIENT_ID = client id form github oauth
GH_BASIC_CLIENT_SECRET = ***********************
```