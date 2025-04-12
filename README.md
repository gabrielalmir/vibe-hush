# Hush

**Hush** is a tool in the making to bring simple, efficient caching to legacy systems. Built with **Go** to run anywhere, it’ll offer a REST API and CLI to store data in memory or on disk, speeding up systems where performance lags.

Picture bridging the gap between your legacy system and modern tech with an open-source project until established tools become viable. That’s Vibe Hush.

Status: Concept in development.

---

## Vibe Hush Development Checklist

### Planning
- [x] Define concept: Lightweight caching for legacies with Go
- [x] Choose tech: Go for backend/CLI, SQLite as fallback
- [ ] Detail scope (REST/CLI features)
- [ ] Plan repository structure (monorepo vibe/vibe-hush)

### MVP Implementation
- [ ] Set up basic Go project (main.go, folder structure)
- [ ] Implement in-memory caching (map + basic TTL)
- [ ] Build initial REST endpoints
  - [ ] `PUT /cache/{key}` to store
  - [ ] `GET /cache/{key}` to retrieve
- [ ] Add basic CLI
  - [ ] `vibe hush set key "data" --ttl=300`
  - [ ] `vibe hush get key`
- [ ] Test MVP in a simulated legacy env (e.g., Windows, no Docker)

### Expansion
- [ ] Add SQLite support as disk fallback
- [ ] Implement auto-expiration for keys (goroutines or timers)
- [ ] Create simple TUI with Bubble Tea (cache viewer)
- [ ] Document REST API (endpoints and examples)

### Finalization
- [ ] Test portability (Windows, macOS, Linux)
- [ ] Write detailed README (install, usage, contribution)
- [ ] Publish repo on GitHub
- [ ] Release v0.1 with compiled binary

---

## Next Steps
Stay tuned for updates! Vibe Hush is starting as a concept, and each step here brings us closer to a working MVP. Got ideas? Open an issue on GitHub (coming soon)!

#OpenSource #Golang #Caching
