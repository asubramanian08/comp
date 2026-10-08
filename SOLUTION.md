# Shared Counter — SvelteKit + Go (GraphQL)

A button in the browser increments one counter that every visitor shares. The page is a
SvelteKit app; the counter lives in a Go server that speaks GraphQL.

```
Browser (page from :5173)                       Go server (:8080)
  click ──► POST /graphql                         withCORS ──► gqlgen ──► Increment()
            {"query":"mutation { increment }"}                            lock; count++; unlock
  shows ◄── {"data":{"increment":42}}  ◄──────────────────────────────────┘
```

## Run it

Needs Go ≥ 1.24 and Node ≥ 22.17.

```sh
# terminal 1 — API on http://localhost:8080/graphql, GraphQL Playground on http://localhost:8080/
cd backend && go run .

# terminal 2 — page on http://localhost:5173
cd frontend && npm install && npm run dev
```

Tests: `cd backend && go test -race ./...` (1,000 simultaneous increments must return 1..1000, each once).

Settings (all optional): `PORT` (default `8080`) and `ALLOWED_ORIGIN` (default
`http://localhost:5173`) for the Go server; `VITE_API_URL` (default
`http://localhost:8080/graphql`) for the frontend.

After editing `backend/schema.graphqls`, regenerate with `cd backend && go generate ./...`.

## Files

| File | What it is |
|---|---|
| `backend/schema.graphqls` | The API: `query { counter }` reads the value, `mutation { increment }` adds 1 and returns the new value |
| `backend/resolver.go` | The shared state: the counter and the mutex that guards it |
| `backend/schema.resolvers.go` | The two functions that read / increment the counter |
| `backend/main.go` | Starts the HTTP server, adds CORS headers |
| `backend/main_test.go` | Concurrency test |
| `backend/gqlgen.yml`, `backend/generated/` | Code-generator settings and its output (never edited by hand) |
| `frontend/src/routes/+page.svelte` | The whole UI and the `fetch` call |

## Answers to the README's questions

- **Backend:** Go + GraphQL. Go runs each request on its own goroutine and ships `sync.Mutex`
  and a race detector, so the concurrency requirement is first-class. GraphQL gives a typed
  schema that documents the API. For a single operation a JSON endpoint would be simpler;
  GraphQL was chosen to learn its trade-offs.
- **What the frontend calls:** `POST /graphql` with `{"query":"mutation { increment }"}`.
- **HTTP method for a JSON API:** `POST` — repeating it changes the result (not idempotent), so
  not `PUT`; it changes state, so never `GET`.
- **Request data:** none beyond the operation. The server decides the +1, so clients can't
  overwrite each other's clicks.
- **Response:** `{"data":{"increment":42}}`, or `{"errors":[…]}` on failure.
- **Numeric type:** `int64`. GraphQL's `Int` is specified as 32-bit; gqlgen doesn't enforce that
  on output, and a click counter won't reach 2³¹, but a custom 64-bit scalar would be the strict fix.
- **Concurrency:** a `sync.Mutex` around read-add-write, with the new value read *before*
  unlocking so two requests can never be handed the same number.
- **How the frontend talks to the backend:** the browser `fetch`es the Go server directly over HTTP.
- **CORS:** needed — `:5173` and `:8080` are different origins, and a JSON POST triggers an
  `OPTIONS` preflight. `main.go` answers it, allowing only `ALLOWED_ORIGIN`.
- **On failure:** the page shows the error and keeps the last number. It does not retry
  automatically, because the failed request may already have been counted. Clicks are never
  blocked; replies can arrive out of order, so the page shows the larger of the old and new value.
- **In-memory limits:** the count resets on restart, and two server copies would each have their
  own counter. Redis `INCR` or a database row would fix both. Other tabs also only see new
  values when they click (polling, Server-Sent Events, or GraphQL subscriptions would fix that).
