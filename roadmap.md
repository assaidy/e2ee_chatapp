# E2EE Chat App Roadmap

---

## ✅ MVP (minimum viable product)
These features give you a working secure 1-to-1 chat system.

- [X] **Auth & user profile**  
  Register, log in, and manage basic user information.

- [x] **Email verification**  
  After registration, send a token via email. User clicks to verify their account. Helps prevent fake accounts.

- [ ] **Create conversations (1-1 chat)**  
  Store metadata for chats between two users (conversation ID, participants).

- [ ] **Offline message queue (store & forward)**  
  If a user is offline, the server holds their encrypted messages and delivers them once they reconnect.

---

## 🚀 Nice-to-have
Makes the app more usable and reliable but not strictly needed for a demo.

- [ ] **Password reset (forgot/change password)**  
  Generate a one-time token for password resets. Logged-in users can change their password directly.

- [ ] **Session management**  
  link/unlink sessions and list active sessions.

- [ ] **Message ordering & delivery receipts**  
  Ensure messages appear in the correct order and support “delivered/read” acknowledgements.

- [ ] **Group chats (basic)**  
  Store group metadata and allow multiple participants to exchange messages.

- [ ] **Group membership management**  
  APIs for inviting, removing, or leaving a group.

- [ ] **Message deletion & retention policies**  
  Support deleting messages (local or global) and optional auto-expiry after a set time.

- [ ] **Media messages**  
  Encrypt files (images, audio, etc.) client-side before upload. Store ciphertext on the server.

- [ ] **Push notifications (privacy-preserving)**  
  Send generic “new message” alerts without leaking message content.

- [ ] **E2EE 1-to-1 chats**
  allow creating E2EE 1-to-1 (only two devices/participants) private chats.

---

## 🔒 Advanced / Security & Operations
Features that harden the system and prepare it for production scale.

- [ ] **Audit logs for suspicious activity**  
  Track repeated failed logins, brute-force attempts, or unusual account activity.

- [ ] **Rate limiting / brute-force protection**  
  Limit how many times an IP/account can hit sensitive endpoints like `/login`.

- [ ] setup CORS (cross origin resource sharing)

- [ ] **Performance/load testing**  
  Benchmark message delivery under heavy load. Optimize DB queries and network usage.

---

## Notes
- The MVP focuses on **secure 1-to-1 chat** with working encryption.  
- Nice-to-have features extend functionality to groups, media, and better UX.  
- Advanced features ensure the system is **secure and scalable** in production.
