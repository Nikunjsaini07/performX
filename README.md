# PerformX: FIFA World Cup 2026 Archive 🏆

**PerformX** is the definitive community archive for the **FIFA World Cup 2026**. 
Score matches and players out of 10, write in-depth reviews, and see exactly what the world is watching. Think of it as Letterboxd or IMDb, but built exclusively for the beautiful game. 

Because of the massive amount of historical and live statistical data required to power this platform, running the full suite locally requires significant database compute. Therefore, this document serves as a comprehensive architectural blueprint of the system.

---

## 🏗️ High-Level System Architecture

PerformX is built as a highly scalable monorepo, strictly separating the client interface, the RESTful API, and internal data-management tooling.

- **`/frontend`**: A modern, server-side rendered (SSR) web interface built with **Next.js** and **TailwindCSS**.
- **`/backend`**: A blazing fast, highly concurrent REST API built purely in **Go (Golang)**.
- **`/admin-tool`**: Internal tooling for managing platform data, API keys, and administrative overrides.

---

## ⚙️ Backend Architecture (Deep Dive)

The backend is engineered for raw speed and minimal memory overhead, avoiding heavy web frameworks or bulky ORMs.

### 1. The HTTP Layer (Standard Library)
Instead of using external routing frameworks like Gin or Fiber, PerformX utilizes Go's native `net/http` `ServeMux`. 
- **Modularity:** Endpoints are grouped logically by domain (`users`, `matches`, `performances`, `teams`, `players`).
- **Middleware:** Custom middleware functions wrap protected routes for `RequireAuth` (JWT) and `RequireAdmin` authorization, ensuring robust, role-based access control.

### 2. Database & Connection Pooling
The platform is backed by **PostgreSQL** (specifically optimized for Neon Serverless Postgres). Because connection establishment is expensive, the backend implements aggressive **Connection Pooling** via `pgxpool`:
- **Max Connections:** 10 concurrent connections to prevent database locking under heavy load.
- **Min Connections:** 2 idle connections always kept alive for zero-latency startup.
- **Lifecycle Management:** Idle connections are gracefully killed after 15 minutes to save memory.

### 3. Database Queries (SQLC)
I decided to skip heavy ORMs like GORM for this project and went with **SQLC** instead. 
- All the database interactions are just pure SQL written in standard `.sql` files.
- I run SQLC to read those files and automatically generate all the Go structs and functions for me (`internal/db/`).
- It's much faster than an ORM, and it catches SQL typos right when I compile the code, which is super nice.

### 4. Authentication (Stateless JWT)
User sessions are entirely stateless. Upon login, the server issues a **JSON Web Token (JWT)**.
- The token is signed using a secure `JWT_SECRET` injected at runtime via OS environment variables.
- The `authMiddleware` intercepts incoming requests, validates the signature, and injects the user's identity directly into the request Context.

---

## 💾 Core Data Models & Features

The PerformX database schema is highly relational and optimized for complex joins required to render the frontend stat pages for all 48 participating countries:

### ⚽ Matches & Teams
- **Teams Page**: Tracks all 48 participating countries.
- **Match Details**: Displays stage, kickoff time, home vs away score, and a summary recap. Matches receive a community-driven average rating out of 10.

### 🏃‍♂️ Players & Performances
- **Performance Details**: Captures an individual player's specific contribution to a single match.
- **Granular Stats**: Tracks Minutes Played, Goals, Assists, Saves, Tackles, Clearances, and Accurate Passes.
- Each individual performance can be rated out of 10 independently of the overall match rating.

### 💬 The Social Graph (Ratings & Reviews)
- **1-10 Ratings**: Users can submit precise 1 to 10 scale ratings on both entire *Matches* and individual *Performances*.
- **Markdown Reviews**: Users can write long-form, Markdown-supported reviews with titles.
- **Threaded Comments**: Users can jump into the comments section of a specific review to discuss it with the community.
- **Likes**: A polymorphic like system allows users to "Like" reviews and specific comments.

---

## 🛠️ Internal Data Tooling (CLI)

Because PerformX requires a massive dataset (real-world FIFA 2026 data) to function, the backend contains custom CLI tools located in `/backend/cmd/` to manage the database lifecycle. 

Right now, these tools are **the most important part of the backend setup** because the app cannot function locally without populating the database first:

- **`cmd/api`**: The main REST API server.
- **`cmd/migrate`**: Automatically applies or rolls back `.sql` schema changes.


---

## 🌐 Frontend Architecture

The `/frontend` directory houses the client application.
- **Framework:** Next.js (React) for Server-Side Rendering (SSR) and superior SEO, vital for public-facing player and match profiles.
- **Styling:** TailwindCSS for utility-first, highly responsive UI design.
- **Data Fetching:** Interfaces directly with the Go REST API, utilizing caching where appropriate for static team/player data, while dynamically fetching volatile data (live 1-10 ratings, Markdown reviews, and comments).

---

## 🚀 Deployment & Environment

PerformX is designed as a secure **12-Factor App** with distinct deployment environments for the frontend, backend, and database:

- **Frontend (Vercel)**: The Next.js client is seamlessly deployed to Vercel, which optimizes the Server-Side Rendering (SSR) and caches static profiles globally at the edge.
- **Database (Neon Serverless)**: The PostgreSQL database is hosted on Neon, allowing the database compute to scale instantly based on traffic spikes during live FIFA matches.
- **Backend (Render / AWS)**: The Go REST API is deployed in a stateless container.
- **Environment Variables**: The system relies entirely on OS-level secrets. Key environment variables required include:
  - **Core**: `DATABASE_URL`, `JWT_SECRET`, `PORT`, `API_URL`, `FRONTEND_URL`
  - **OAuth & Auth**: `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`
  - **Media & Storage**: `CLOUD_NAME`, `CLOUDINARY_API`, `CLOUDINARY_SECRET`
  - **AI & Search Services**: `GROG_API_KEY`, `COHERE_API_KEY`, `SERPER_API_KEY`
  - **Email**: `BREVO_API`, `SENDER_EMAIL`
  - **Football Data (Fotmob)**: `RAPID_API` keys for live/historical match ingestion.
  
  While a `.env` file is used for local development, production providers (Vercel/Neon/Render) inject these directly into secure system memory at runtime. The server is configured to fail fast (`log.Fatal`) if critical secrets are missing.
