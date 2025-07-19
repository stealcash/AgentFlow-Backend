# AgentFlow Backend

 **AgentFlow-Backend** is the backend service for the **AgentFlow** platform — built with **Go (Golang)**, **Gin**, and **PostgreSQL**.
It powers core APIs, authentication, chatbot logic, analytics, and more.

---

## ⚙️ Requirements

* **Go 1.20+**
* **PostgreSQL 13+**
* **Git**

---

##  Getting Started

### 1️⃣ Clone the repository

```bash
git clone https://github.com/stealcash/AgentFlow-Backend.git
cd AgentFlow-Backend
```

### 2️⃣ Install dependencies

```bash
go mod tidy
```

### 3️⃣ Setup your `.env` and `config.toml`

* Copy `.env.example` to `.env`:

  ```bash
  cp .env.example .env
  ```

* Copy the sample config to `config.toml`:

  ```bash
  cp config/config-sample.toml config/config.toml
  ```

* Edit `.env` and `config/config.toml` with your local database credentials and other settings.

---

### 4️⃣ Prepare your database

Make sure PostgreSQL is running, and your database matches your `.env` and `config.toml` settings.

Run migrations if required.

---

### 5️⃣ Run the server

```bash
go run main.go
```

The server will start on `http://localhost:PORT` (defined in `.env`).

---



## 🖥️ Supporting Frontend

You can find the official **AgentFlow Frontend** repository here:
➡️ [AgentFlow-Frontend](https://github.com/stealcash/AgentFlow-Frontend)
