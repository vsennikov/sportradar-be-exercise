## Database Design

Here is the Entity-Relationship Diagram (ERD) for the project:

![Database ERD](./erd.drawio.png)

## API Endpoints

All endpoints are prefixed with `/api/v1`.

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/events` | Gets a paginated list of events. |
| `GET` | `/events/:id` | Gets a single event by its unique ID. |
| `POST` | `/events` | Creates a new event. |

### List Events (Filtering & Pagination)

The `GET /events` endpoint supports the following query parameters:

* **`page`**: (Optional) The page number you want to view.
    * *Example:* `?page=2`
* **`limit`**: (Optional) The number of events to show per page.
    * *Example:* `?limit=5`
* **`sport_id`**: (Optional) Filters the list for a specific sport.
    * *Example:* `?sport_id=1`

---

## 🧪 How to Test with Postman

This project includes a Postman collection to make testing easy.

1.  Open Postman.
2.  Click the **Import** button.
3.  Find the `postman/` folder in this project.
4.  Drag and drop the `sportradar_api.postman_collection.json` file into the Postman window.
5.  (Optional) Import the `local.postman_environment.json` file in the same way. This sets up the `{{baseUrl}}` variable for you.