# ReadRate - Book Review API

Ini adalah RESTful API yang dibangun dengan Golang menggunakan Gin framework dan PostgreSQL sebagai penyimpanan data. Aplikasi ini dirancang untuk mengelola **buku** dan **kategori**, serta memungkinkan **pengguna untuk mendaftar, login, dan memberikan ulasan (review) untuk buku**.

---

## 🧰 Tech Stack

* **Go (Golang)**
* **Gin Web Framework**
* **PostgreSQL**
* **`bcrypt`** untuk hashing password
* **`.env` Configuration** dengan `godotenv`

---

## 📦 Setup

1.  **Clone the repository:**

    ```bash
    git clone [https://github.com/KonekkoNekko/sb-go-readrate-nabiel.git](https://github.com/KonekkoNekko/sb-go-readrate-nabiel.git)
    cd sb-go-readrate-nabiel
    ```

2.  **Set up your `.env` file** inside `config/.env`. Pastikan detail koneksi database Anda sesuai:

    ```env
    DATABASE_URL="postgresql://user:password@host:port/dbname" # Contoh Railway URL
    PGDATABASE="your_database"
    PGHOST="your_host"
    PGPASSWORD="your_password"
    PGPORT="5432"
    PGUSER="your_user"
    ```

3.  **Run the application:**
    ```bash
    go mod tidy # Pastikan semua dependensi terinstall
    go run .
    ```
    Aplikasi akan terhubung ke database dan menjalankan migrasi skema secara otomatis saat pertama kali dijalankan, membuat tabel yang diperlukan (users, categories, books, reviews).

---

## 📘 API Documentation

Aplikasi ini menggunakan **Basic Authentication** untuk sebagian besar *endpoint* yang memerlukan otorisasi (misalnya, membuat, memperbarui, atau menghapus data buku, kategori, atau ulasan). Beberapa *endpoint* bersifat publik dan dapat diakses tanpa autentikasi.

### 🔐 Authentication & Users Endpoints

| Method | Endpoint | Description |
| :----- | :------- | :---------- |
| `POST` | `/register` | Register a new user (password is hashed) |
| `AUTH` | Basic Auth | Provide `Username` and `Password` in the `Authorization` header for protected routes |

### 📚 Books Endpoints

| Method | Endpoint | Description | Authentication Required |
| :----- | :------- | :---------- | :---------------------- |
| `GET` | `/books` | Get all books | ❌ |
| `GET` | `/books/:id` | Get a book by ID | ❌ |
| `POST` | `/books` | Create a new book | ✅ |
| `PUT` | `/books/:id` | Update a book by ID | ✅ |
| `DELETE` | `/books/:id` | Delete a book by ID | ✅ |

### 🗂️ Categories Endpoints

| Method | Endpoint | Description | Authentication Required |
| :----- | :------- | :---------- | :---------------------- |
| `GET` | `/categories` | Get all categories | ❌ |
| `GET` | `/categories/:id` | Get a category by ID | ❌ |
| `GET` | `/categories/:id/books` | Get books under a specific category | ❌ |
| `POST` | `/categories` | Create a new category | ✅ |
| `PUT` | `/categories/:id` | Update a category by ID | ✅ |
| `DELETE` | `/categories/:id` | Delete a category by ID | ✅ |

### ⭐ Reviews Endpoints

| Method | Endpoint | Description | Authentication Required |
| :----- | :------- | :---------- | :---------------------- |
| `GET` | `/books/:id/reviews` | Get all reviews for a specific book | ❌ |
| `GET` | `/users/:user_id/reviews` | Get all reviews written by a specific user | ❌ |
| `POST` | `/books/:book_id/reviews` | Create a new review for a book | ✅ |
| `PUT` | `/reviews/:id` | Update an existing review (only by owner) | ✅ |
| `DELETE` | `/reviews/:id` | Delete a review (only by owner) | ✅ |

---

## 📄 License

This project is licensed under the MIT License.