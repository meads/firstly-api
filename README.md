# Notes App

This application achieves two primary objectives. One is to serve as a template for horizontally layered architecture written in Go and two is to provide an example of a backend that authenticates web clients as well as potientially mobile clients using a token based architecture which is described in the subsequent sections of this document.

This document describes the stateless-but-revocable JWT authentication system used in this Notes application. The architecture prioritizes simplicity, clean database state, and seamless user experience across browser sessions.

The application allows a user to register using a username and password. Following the initial signup the user can create, read, update and delete notes. The functionality was intentionally made to be simple so as just to provide the proof of concept.

## 🛠️ Tech Stack

* **Frontend:** [Vite/React]
* **Styling:** [CSS]
* **HTTP Client:** [Fetch API]
* **HTTP Server:** [Gin]
* **Backend:** [Go] (version 1.26)
* **Tokens:** [JWT]
* **Database:** [sqlc/Postgresql]

## 📦 Getting Started

### Prerequisites
* Docker Desktop [typical docker install steps](https://www.docker.com/get-started/)
* A running instance of the [Notes API](https://github.com/meads/firstly-api)


### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/meads/firstly-api.git
    cd firstly-api
    ```

2. Configure environment variables:
   Create a `.env` file in the root directory using the example found in CONTRIBUTING.md

3. Build and start Docker containers:
   ```bash
   docker compose build 
   docker compose up
   ```


### Running the App

Open `http://localhost:3000` in your browser to view the UI.



<hr>


# 🏗️ Architecture Overview

The system uses a cookie-less dual-token architecture (Access Token + Refresh Token) stored on the client side via sessionStorage. 

To maintain control over active sessions without checking the database on every single API request, the system uses a Sessions table. The sessions table PRIMARY KEY is a uuid that mirrors the refresh token RegisteredClaims.ID. The user_id is the user the session is for. The refresh_token is the actual JWT refresh token string. The is_revoked column is a boolean that allows revocation of a session/refresh_token. The expires_at is another field that mirrors the JWT refresh token but for the expires timestamp.

### 🗄️ Database Schema

```sql
CREATE TABLE sessions (
    id VARCHAR(36) PRIMARY KEY NOT NULL,  -- Unique UUID or Session ID
    user_id BIGINT NOT NULL,              -- Foreign key to Users table
    refresh_token VARCHAR(512) NOT NULL,       -- The issued 24-hour Refresh JWT
    is_revoked BOOLEAN DEFAULT FALSE NOT NULL, -- Flag to invalidate the session preventing token rotation
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### 🔄 Lifecycle Flows

#### Registration & Authentication (/register & /login)

  1. The client submits credentials.

  2. The backend verifies credentials.

  3. The backend generates a 15-minute Access Token and a 24-hour Refresh Token.

  4. The backend inserts a new record into the sessions table.

  5. The backend returns the following JSON payload to the client:

```json
{
  "sessionId": "uuid-string",
  "userId": 123,
  "accessToken": "eyJhbGci...",
  "refreshToken": "eyJhbGci..."
}
```

#### Client-Side Storage

The client stores all returned properties directly in sessionStorage.

#### Transparent Token Refresh (The Interceptor Queue)

To keep the user logged in seamlessly, the client wraps the native browser fetch function (fetchClient).

```
Expired Access Token Triggered -> 401 Unauthorized
                                      |
                         Queue all pending requests
                                      |
                     POST /refresh (Send Refresh Token)
                                      |
              [Backend verifies Session exists in DB]
                                      |
                     Returns NEW 15-min Access Token
                                      |
              Replay queued requests with updated header 
```

 * Static Refresh Approach: The Refresh Token remains static for its 24-hour lifespan to prevent race conditions during concurrent frontend fetches.
 * Failure Catch: If the /refresh call fails (e.g., token expired or deleted from DB), the queue is rejected, sessionStorage is wiped, and the user is redirected to the login screen.

#### Explicit Logout (/logout)

  1. The user clicks "Log Out".
  2. The client fires a POST /logout passing the sessionId.
  3. The backend explicitly deletes the corresponding row from the sessions table.
  4. The client purges sessionStorage.

### 🧹 Automated Database Cleanup

Because sessionStorage is destroyed when a user closes a browser tab, the server is never notified of "abandoned" sessions. Left unchecked, the sessions table would grow indefinitely.
To prevent this, create a scheduled routine that runs once every 24 hours to prune expired sessions from the database.


SQL Cleanup Query

Any session older than 24 hours is guaranteed to have an expired Refresh Token and is safe to delete.

```sql
DELETE FROM sessions 
WHERE created_at < NOW() - INTERVAL '24 hours';
```


### ⚠️ Security Notes

1. XSS Protection: Storing tokens in sessionStorage makes them accessible via JavaScript. Ensure strict Content Security Policies (CSP) are active and input sanitization is strictly applied on all user-submitted notes to mitigate Cross-Site Scripting risks.
2. HTTPS Only: All token exchanges must occur over TLS/HTTPS to protect credentials and JWTs from interception.



### 📄 License

This project is open-source and available under the **MIT License**.

