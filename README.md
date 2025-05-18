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


### 📝 Request Body Payloads

Bagian ini merinci struktur JSON yang diharapkan untuk permintaan `POST` dan `PUT`.

---

#### 📌 User Registration (`POST /register`)

- **`username`** (string, wajib): Nama pengguna unik.  
- **`password`** (string, wajib): Kata sandi pengguna.

```json
{
  "username": "userbaru",
  "password": "PasswordAman123!"
}
```

---

#### 📚 Create Book (`POST /books`)

- **`title`** (string, wajib): Judul buku.  
- **`category_id`** (integer, wajib): ID kategori tempat buku berada.  
- **`description`** (string, opsional): Deskripsi singkat buku.  
- **`image_url`** (string, opsional): URL gambar sampul buku.  
- **`release_year`** (integer, wajib): Tahun rilis buku.  
- **`price`** (integer, wajib): Harga buku (dalam mata uang lokal).  
- **`total_page`** (integer, wajib): Jumlah halaman buku.  
- **`thickness`** (string, opsional): Ketebalan buku (misalnya, "tipis", "sedang", "tebal").

```json
{
  "title": "Filosofi Teras",
  "category_id": 4,
  "description": "Buku tentang filosofi stoikisme untuk kehidupan sehari-hari.",
  "image_url": "https://example.com/filosofi_teras.jpg",
  "release_year": 2018,
  "price": 85000,
  "total_page": 300,
  "thickness": "sedang"
}
```

---

#### 📝 Update Book (`PUT /books/:id`)

Semua bidang bersifat opsional. Kirimkan hanya bidang yang ingin Anda perbarui.

```json
{
  "description": "Versi revisi dengan bab tambahan.",
  "price": 90000
}
```

---

#### 🗂️ Create Category (`POST /categories`)

- **`name`** (string, wajib): Nama kategori baru.

```json
{
  "name": "Self-Improvement"
}
```

---

#### 🛠️ Update Category (`PUT /categories/:id`)

- **`name`** (string, opsional): Nama kategori yang diperbarui.

```json
{
  "name": "Pengembangan Diri"
}
```

---

#### ⭐ Create Review (`POST /books/:book_id/reviews`)

- **`rating`** (integer, wajib): Penilaian buku (1-5).  
- **`comment`** (string, opsional): Komentar atau ulasan tentang buku.  
  > Catatan: `user_id` secara otomatis diambil dari pengguna yang terautentikasi.

```json
{
  "rating": 5,
  "comment": "Buku yang sangat menginspirasi dan mudah dipahami!"
}
```

---

#### 🔁 Update Review (`PUT /reviews/:id`)

Semua bidang bersifat opsional. Kirimkan hanya bidang yang ingin Anda perbarui.

```json
{
  "rating": 4,
  "comment": "Konten bagus, tapi kadang agak berulang."
}
```


## 📄 License

This project is licensed under the MIT License.