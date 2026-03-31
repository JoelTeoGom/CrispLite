# CrispLite — TODO

## Users (REST)
- [ ] `POST /api/users` — Create User (register)
- [ ] `POST /api/users/:id/contacts` — Add Contact (creates empty conversation between 2 users, enforce user_a < user_b to avoid dupes)
- [ ] `DELETE /api/users/:id/contacts/:contactId` — Remove Contact (soft-delete: hides conversation for deleter, keeps it for the other user. Needs contacts table with active flag)

## Chat (REST)
- [ ] `GET /api/users/:id/conversations?cursor=&limit=50` — List Conversations (paginated, with last message, ordered by last_message_at DESC)
- [ ] `GET /api/conversations/:id/messages?cursor=&limit=50` — Load Messages (paginated, ordered by timestamp DESC)

## Chat (WebSocket)
- [ ] `WS /ws/v1/chat` — Already exists. Handles incoming messages and pushes real-time updates (new messages, conversation reorder in frontend)

## DB
- [ ] Add contacts table to migration (user_id, contact_id, active, created_at)

## Cross-cutting
- [ ] Auth (JWT) — needed for all endpoints and WS

## Future
- [ ] Read receipts / typing indicators / online status
- [ ] Search / message filtering
- [ ] Group conversations
- [ ] Media / attachments
