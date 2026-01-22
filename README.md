# Lathi - Javanese Interactive Visual Novel Learning Platform

**Lathi** is a Javanese language learning platform based on an Interactive Visual Novel. It combines branching narratives, interactive dialogues, and cinematic visual elements to deliver an immersive learning experience.

This backend repository powers the Lathi platform, providing a robust API for user management, gamified story progression, vocabulary mastery, and global leaderboards.

## 📜 System Description

Lathi transforms the traditional language learning process into an engaging game. Users play through chapters of a story, interacting with characters using Javanese. The system tracks:

- **Narrative Progression:** Users make choices that affect the story flow and character moods ("Hearts").
- **Vocabulary Acquisition:** Players unlock new Javanese words (Ngoko & Krama) by encountering them in the story.
- **Gamification:** Users earn titles (Cantrik, Abdi, Priyayi), badges, and scores based on their completion and vocabulary collection.

## 🚀 Key Features

### 📖 Interactive Story Mode

- **Visual Novel Engine:** API supports chapters, slides, background images, and character sprites.
- **Branching Choices:** User decisions impact the "Mood/Heart" system and conversation outcomes.
- **Session Tracking:** Saves progress (current slide, hearts, history log) to allow resuming anytime.
- **Unlockables:** Automatically unlocks vocabulary entries upon encountering them in dialogue.

### 📚 Dictionary

- **Searchable Database:** Look up words in Ngoko, Krama, and Indonesian.
- **Collection System:** Tracks which words a user has "unlocked" through gameplay.

### 🏆 Gamification & Social

- **Global Leaderboard:** Ranks users based on a composite score of chapters completed and words collected.
- **Badges System:** Awards badges for specific achievements (e.g., "Perfect Heart", "Vocab Collector").
- **Dynamic Titles:** User titles update automatically based on progress (Cantrik -> Abdi -> Priyayi).

## 🛠 Tech Stack

### Core

- **Language:** [Go 1.25+](https://go.dev/)
- **Framework:** [Fiber v2](https://gofiber.io/)

### Database & Storage

- **Relational DB:** [PostgreSQL](https://www.postgresql.org/)
- **ORM:** [GORM](https://gorm.io/)
- **Caching:** [Redis](https://redis.io/)
- **Object Storage:** [MinIO](https://min.io/)

### Infrastructure & Tools

- **Authentication:** [JWT](https://jwt.io/)
- **Config Management**: [Viper](https://github.com/spf13/viper)
- **Validator**: [Validator v10](https://github.com/go-playground/validator)
- **Hashing and Encryption**: [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- **Mailer**: [Gomail](https://github.com/go-gomail/gomail)
- **Task Scheduler**: [Cron](https://github.com/robfig/cron)
- **Containerization:** [Docker](https://www.docker.com/)
- **Documentation:** [OpenAPI 3.1](https://www.openapis.org/)

## 📂 Project Structure

The project follows a modular structure based on Clean Architecture principles:

```bash
.
├── cmd
│   ├── api          # Main entry point for the REST API
│   └── bootstrap    # App initialization
├── db
│   ├── migration    # Database schema migrations
│   └── seed         # Data seeders
├── docs             # OpenAPI documentation
├── internal
│   ├── app          # Application modules
│   │   ├── auth
│   │   ├── dictionary
│   │   ├── leaderboard
│   │   ├── story
│   │   └── user
│   ├── config       # Environment configuration
│   ├── domain       # Domain entities, DTOs, and interface contracts
│   └── infra        # Infrastructure implementations
│   └── middleware   # Application middlewares
└── pkg              # Shared utilities
```

## ⚙️ Setup and Installation

### Prerequisites

- Docker and Docker Compose
- Make (optional, for using Makefile commands)

### 1. Clone the Repository

```bash
git clone https://github.com/Ablebil/lathi-be.git
cd lathi-be
```

### 2. Environment Configuration

Copy the example environment file and configure it.

```bash
cp .env.example .env
```

_Ensure you fill in the SMTP details and MinIO credentials in `.env` for full functionality._

**Production Deployment Note:**

> When deploying to a production environment, you must change `APP_ENV` to `production` in your `.env` file.
> This setting enables security features (such as Secure Cookies) and optimizes logging performance.

### 3. Run with Docker

The easiest way to run the application is using Docker Compose, which sets up the App, PostgreSQL, Redis, and Swagger UI.

```bash
make run
# OR directly
docker compose up
```

The API will be available at `http://localhost:8080`.

### 4. Database Migration & Seeding

To populate the database with initial data:

```bash
# Run migration
make migrate-up
# OR directly
docker compose exec app /app/server migrate -action up

# Run seeder (all domains)
make seed-all
# OR directly
docker compose exec app /app/server seed

# Run seeder (specific domain)
docker compose exec app /app/server seed domain user
```

## 📖 API Documentation

Once the application is running, you can access the interactive Swagger API documentation at http://localhost:8081

## 🔌 API Endpoints Summary

### Auth

| Method | Endpoint                | Description              |
| ------ | ----------------------- | ------------------------ |
| POST   | `/api/v1/auth/register` | Register a new user      |
| POST   | `/api/v1/auth/verify`   | Verify email address     |
| POST   | `/api/v1/auth/login`    | Login and receive tokens |
| POST   | `/api/v1/auth/refresh`  | Refresh access token     |
| POST   | `/api/v1/auth/logout`   | Logout user              |

### Story

| Method | Endpoint                               | Description                     |
| ------ | -------------------------------------- | ------------------------------- |
| GET    | `/api/v1/stories/chapters`             | List all chapters and progress  |
| GET    | `/api/v1/stories/chapters/:id/content` | Get chapter content             |
| GET    | `/api/v1/stories/chapters/:id/session` | Get chapter progress            |
| POST   | `/api/v1/stories/chapters/:id/start`   | Start a chapter session         |
| POST   | `/api/v1/stories/action`               | Submit choice/next slide action |

### Dictionary

| Method | Endpoint               | Description                |
| ------ | ---------------------- | -------------------------- |
| GET    | `/api/v1/dictionaries` | Search and list vocabulary |

### User

| Method | Endpoint                | Description               |
| ------ | ----------------------- | ------------------------- |
| GET    | `/api/v1/users/profile` | Get user stats and badges |
| PATCH  | `/api/v1/users/profile` | Update profile info       |
| DELETE | `/api/v1/users/account` | Delete account            |

### Leaderboard

| Method | Endpoint               | Description          |
| ------ | ---------------------- | -------------------- |
| GET    | `/api/v1/leaderboards` | Get top global users |

## 📝 License

This project is licensed under the MIT License.
