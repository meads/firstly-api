## 🗺️ Roadmap / TODO

### Backend
- [x] Register route to create new users, create tokens and session.
- [x] Users routes to retrieve and modify database users.
- [x] Login route to authenticate existing users, create tokens and session.
- [x] Sessions route to retrieve and modify sessions. 
- [ ] Notes route to retrieve and modify notes db records.
- [ ] Notes table for persisting notes records.
- [x] Users table for persisting uesrs records.
- [x] Sessions table for persisting sessions records.
- [x] Reliable means of hashing and comparing passwords during Registration and Login.
- [x] JWT middleware to protect certain routes by token validation. 
- [x] Database first design to generate data layer calls from sql statements.
- [x] Generate mocks for interfaces allowing mocks for unit testing code that modifies database state.
- [x] Generation of access & refresh tokens for custom claims using JWT library v5.
- [ ] Config struct to store environment variables in main package, validate and provide ease of use.
- [ ] Tests for all code located in main package.
- [ ] Near 100% code coverage for unit testing of all backend code.

### Frontend
- [x] UI to register users.
- [x] UI to login existing users.
- [ ] UI to view sessions and tokens.
- [ ] UI to create and list notes.
- [ ] UI to view realtime token expiration.
 

