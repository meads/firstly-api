## Notes App Authentication Architecture

This document describes the stateless-but-revocable JWT authentication system used in this Notes application. The architecture prioritizes simplicity, clean database state, and seamless user experience across browser sessions.

### 🏗️ Architecture Overview

The system uses a dual-token architecture (Access Token + Refresh Token) stored entirely on the client side via sessionStorage. To maintain control over active sessions without checking the database on every single API request, the system introduces a lightweight Sessions table.

```
+-----------------------------------+
|          Client Browser           |
+-----------------------------------+
                  |
    1. /login or /register
                  v
+-----------------------------------+
|          Backend Server           |
+-----------------------------------+
                  |
   2. Generates Tokens & Sessions Record
                  v
+-----------------------------------+
|         Database (PostgreSQL)     |
+-----------------------------------+
```

### 🗝️ Token Specifications

  * Access Token: Short-lived (15 minutes). Sent in the Authorization: Bearer [token] header. Verified entirely statelessly by backend middleware.
  * Refresh Token: Long-lived (24 hours). Sent in the JSON body to the /refresh endpoint to request a new access token.

## 🛠️ Tech Stack

* **Frontend:** [Vite/React]
* **Styling:** [CSS]
* **HTTP Client:** [Fetch API]
* **HTTP Server:** [Gin]
* **Backend:** [Go]
* **Tokens:** [JWT]
* **Database:** [sqlc/Postgresql]

## 📦 Getting Started

### Prerequisites
* Go (version 1.26)
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


### 🗄️ Database Schema

The sessions table tracks active, valid login instances per browser tab/device.

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

  2. The backend verifies credentials (or hashes a new password using bcrypt).

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
 * Isolation Benefit: Because sessionStorage is strictly sandboxed to a single browser tab, multiple open tabs maintain separate tokens and session contexts, eliminating cross-tab token conflicts.


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
To prevent this, a scheduled routine runs once every 24 hours to prune expired sessions from the database.


SQL Cleanup Query

Any session older than 24 hours is mathematically guaranteed to have an expired Refresh Token and is safe to delete.

```sql
DELETE FROM sessions 
WHERE created_at < NOW() - INTERVAL '24 hours';
```

Production Implementations

Option A: Node.js Scheduled Task (node-cron)

```javascript
const cron = require('node-cron');
const db = require('./db');

// Runs every night at midnight (00:00)
cron.schedule('0 0 * * *', async () => {
  try {
    const result = await db.query(
      "DELETE FROM sessions WHERE created_at < NOW() - INTERVAL '24 hours'"
    );
    console.log(`[Cleanup] Purged ${result.affectedRows} expired sessions.`);
  } catch (err) {
    console.error('[Cleanup Error]', err);
  }
});
```

Option B: Using pg_cron

The pg_cron extension turns PostgreSQL into a cron-based scheduler, allowing you to run SQL commands directly from inside the database.

1. Enable the extension

You must add pg_cron to your shared_preload_libraries in postgresql.conf and restart the database. Once done, run:

```sql
CREATE EXTENSION pg_cron;
```

2. Schedule the DELETE query

Run the cron.schedule function. The first argument is a standard cron expression (0 2 * * * means everyday at 2:00 AM), and the second argument is your SQL query:

```sql
SELECT cron.schedule(
    'daily-cleanup-job',      -- Job name
    '0 2 * * *',              -- Cron schedule (2:00 AM daily)
    $$DELETE FROM sessions WHERE created_at < NOW() - INTERVAL '24 hours'$$
);
```

3. Managing the Job

```sql
-- View active jobs
SELECT * FROM cron.job;

-- Unschedule/delete the job
SELECT cron.unschedule('daily-cleanup-job');
```

### ⚠️ Security Notes

1. XSS Protection: Storing tokens in sessionStorage makes them accessible via JavaScript. Ensure strict Content Security Policies (CSP) are active and input sanitization is strictly applied on all user-submitted notes to mitigate Cross-Site Scripting risks.
2. HTTPS Only: All token exchanges must occur over TLS/HTTPS to protect credentials and JWTs from interception.

### 📄 License

This project is open-source and available under the **MIT License**.
