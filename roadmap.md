# E2EE Chat App Roadmap

---

## ✅ MVP (minimum viable product)
These features give you a working secure 1-to-1 chat system.

- [X] **Auth & user profile**  
  Register, log in, and manage basic user information.

- [x] **Email verification**  
  After registration, send a token via email. User clicks to verify their account. Helps prevent fake accounts.

- [ ] **Password reset (forgot/change password)**  
  Generate a one-time token for password resets. Logged-in users can change their password directly.

- [ ] **Public key upload during registration**  
  Clients generate a key pair locally. The public key is uploaded to the server so others can start encrypted chats.

- [ ] **Secure session establishment (X3DH / Double Ratchet)**  
  Protocols to let two users agree on shared session keys securely. This is how chats stay end-to-end encrypted.

- [ ] **Create conversations (1-1 chat)**  
  Store metadata for chats between two users (conversation ID, participants).

- [ ] **Send/receive encrypted messages**  
  Clients encrypt before sending. The server only stores ciphertext, not plaintext.

- [ ] **Offline message queue (store & forward)**  
  If a user is offline, the server holds their encrypted messages and delivers them once they reconnect.

---

## 🚀 Nice-to-have
Makes the app more usable and reliable but not strictly needed for a demo.

- [ ] **Key rotation & backup strategy**  
  Allow users/devices to replace old keys and recover from lost devices without losing chats.

- [ ] **Device management**  
  APIs to link/unlink devices and list active sessions. Each device has its own key pair.

- [ ] **Message ordering & delivery receipts**  
  Ensure messages appear in the correct order and support “delivered/read” acknowledgements.

- [ ] **Group chats (basic)**  
  Store group metadata and allow multiple participants to exchange encrypted messages.

- [ ] **Group membership management**  
  APIs for inviting, removing, or leaving a group.

- [ ] **Group key distribution**  
  Manage group keys securely. When someone joins/leaves, rotate keys so ex-members lose access.

- [ ] **Message deletion & retention policies**  
  Support deleting messages (local or global) and optional auto-expiry after a set time.

- [ ] **Media messages**  
  Encrypt files (images, audio, etc.) client-side before upload. Store ciphertext on the server.

- [ ] **Push notifications (privacy-preserving)**  
  Send generic “new message” alerts without leaking message content.

---

## 🔒 Advanced / Security & Operations
Features that harden the system and prepare it for production scale.

- [ ] **Audit logs for suspicious activity**  
  Track repeated failed logins, brute-force attempts, or unusual account activity.

- [ ] **Rate limiting / brute-force protection**  
  Limit how many times an IP/account can hit sensitive endpoints like `/login`.

- [ ] **Performance/load testing**  
  Benchmark message delivery under heavy load. Optimize DB queries and network usage.

---

## Notes
- The MVP focuses on **secure 1-to-1 chat** with working encryption.  
- Nice-to-have features extend functionality to groups, media, and better UX.  
- Advanced features ensure the system is **secure and scalable** in production.
