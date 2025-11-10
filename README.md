# Sportradar Coding Exercise

This project is a sports event calendar API built with Go, Gin, PostgreSQL, and Docker. It includes a full backend with CRUD operations for events, sports, teams, and venues. It also features a simple frontend built with JavaScript and the Pico.css framework.

## Table of Contents

- [How to Run](#-how-to-run)
- [How to Test with Postman](#-how-to-test-with-postman)
- [Running Tests](#-running-tests)
- [API Endpoints](#api-endpoints)
- [Database Design](#database-design)
- [Assumptions & Design Decisions](#-assumptions--design-decisions)

---

## 🚀 How to Run

1.  **Prerequisites:**
    * Docker
    * Docker Compose

2.  **Clone the repository:**
    ```bash
    git clone git@github.com:vsennikov/sportradar-be-exercise.git
    cd sportradar-be-exercise
    ```

3.  **Create the environment file:**
    Copy the example file. This contains all the necessary passwords and port configurations.
    ```bash
    cp .env.example .env
    ```

4.  **Build and run the application:**
    This command will build the Go application, start the PostgreSQL database, and run application.
    ```bash
    docker compose up --build
    ```

5.  **Access the application:**
    * **Event Calendar (Frontend):** `http://localhost:8080/`
    * **API (Backend):** `http://localhost:8080/api/v1`

---

## 🧪 How to Test with Postman

This project includes a Postman collection to make testing the API easy.

1.  Open Postman.
2.  Click the **Import** button (usually in the top-left).
3.  Find the `postman/` folder in this project.
4.  Drag and drop the `sportradar_api.postman_collection.json` file into the Postman window.
5.  (Optional) Import the `local.postman_environment.json` file in the same way. This sets up the `{{baseUrl}}` variable for you.

---

## 🧪 Running Tests

This project includes unit and integration tests.

### Quick Start

```bash
# Run all tests (unit + integration)
go test ./...

# Run only unit tests (skip integration tests)
go test -short ./...

# Run tests with coverage
go test -cover ./...
```

For detailed information about the test suite see [TESTING.md](./TESTING.md) for complete testing documentation.

---

## API Endpoints

All endpoints are prefixed with `/api/v1`.

### Events

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/events` | Gets a paginated list of events. |
| `GET` | `/events/:id` | Gets a single event by its unique ID. |
| `POST` | `/events` | Creates a new event. (Returns new ID) |
| `PATCH` | `/events/:id` | Partially updates an existing event. |
| `DELETE`| `/events/:id` | Deletes an event. |

**Filtering & Pagination for `GET /events`:**

The `GET /events` endpoint supports the following query parameters:
* **`page`**: (Optional) The page number you want to view. *Example:* `?page=2`
* **`limit`**: (Optional) The number of events to show per page. *Example:* `?limit=5`
* **`sport_id`**: (Optional) Filters the list for a specific sport. *Example:* `?sport_id=1`
* **`date_from`**: (Optional) Filters for events on or after a date. *Example:* `?date_from=2025-01-01`

### Sports

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/sports` | Gets a list of all sports. |
| `GET` | `/sports/:id` | Gets a single sport by its unique ID. |
| `POST` | `/sports` | Creates a new sport. (Returns new ID) |
| `PUT` | `/sports/:id` | Replaces an existing sport. |
| `DELETE`| `/sports/:id` | Deletes a sport (Fails if in use). |

### Teams

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/teams` | Gets a list of all teams. |
| `GET` | `/teams/:id` | Gets a single team by its unique ID. |
| `POST` | `/teams` | Creates a new team. (Returns new ID) |
| `PATCH` | `/teams/:id` | Partially updates an existing team. |
| `DELETE`| `/teams/:id` | Deletes a team (Fails if in use). |

### Venues

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/venues` | Gets a list of all venues. |
| `GET` | `/venues/:id` | Gets a single venue by its unique ID. |
| `POST` | `/venues` | Creates a new venue. (Returns new ID) |
| `PATCH` | `/venues/:id` | Partially updates an existing venue. |
| `DELETE`| `/venues/:id` | Deletes a venue. |

---

## Database Design

Here is the Entity-Relationship Diagram (ERD) for the project:

![Database ERD](./erd.drawio.png)

---

## 🏛️ Assumptions & Design Decisions

This project was built with scalable architecture in mind, based on the following key decisions:

* **3-Layer Architecture:** The application is strictly separated into three layers: `Handlers` (Controllers), `Services` (Business Logic), and `Repositories` (Data Access). This follows the Single Responsibility Principle.
* **Dependency Inversion:** Interfaces (e.g., `EventService`, `EventRepository`) are defined at the `Service` layer. This separates business logic from the specific database or API implementation.
* **3-Model Separation:** To ensure clean separation, three types of models are used:
    1.  **DTOs:** For JSON binding/responses (e.g., `EventDTO`).
    2.  **Service Models:** "Clean" structs for business logic (e.g., `Event`).
    3.  **DB Models:** "Flat" structs for `sqlx` mapping (e.g., `eventDB`).
* **Database Strategy:**
    * **`sqlx` over ORM:** `sqlx` was chosen over a full ORM (like GORM) to demonstrate proficiency with raw, efficient SQL, as required by the exercise.
    * **N+1 Prevention:** The `baseEventSelectQuery` uses multiple `JOIN`s to fetch all required data in a single query, preventing the "N+1" problem.
* **Context Propagation:** `context.Context` is passed through every layer to handle request cancellation and prevent resource leaks (e.g., "zombie queries").
* **Simple Frontend:** The frontend is intentionally simple (plain JS and `Pico.css`) to meet the requirement without using a heavy framework (like React/Vue).
