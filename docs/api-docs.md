# TaskForge API Documentation

Dokumentasi lengkap RESTful API untuk backend **TaskForge**.

- **Base URL**: `http://localhost:3000/api`
- **Content-Type**: `application/json`
- **Format Respons**: JSend-compliant JSON

---

## 1. Mekanisme Keamanan

### A. Double Submit CSRF Protection
Untuk seluruh request yang memodifikasi data (`POST`, `PUT`, `DELETE`), API ini dilindungi oleh middleware CSRF:
1. Panggil `GET /api/auth/csrf` sebelum melakukan request formulir/login/registrasi.
2. Server akan mengembalikan `csrf_token` dalam response body dan otomatis menyetel cookie `csrf_`.
3. Sertakan header berikut pada setiap request `POST`, `PUT`, `DELETE`:
   ```http
   X-CSRF-Token: <nilai_csrf_token>
   ```
4. Pastikan cookie `csrf_` ikut terkirim bersama request (ditangani otomatis oleh browser jika `withCredentials: true` atau Postman).

### B. Autentikasi JWT & HTTPOnly Cookie
Setelah login berhasil, server akan menyetel dua cookie dengan flag `HttpOnly`:
- **`access_token`**: Digunakan untuk mengautentikasi setiap protected endpoint (masa aktif default: 15 menit).
- **`refresh_token`**: Digunakan untuk memperpanjang `access_token` tanpa perlu login ulang (masa aktif default: 7 hari).

> **Fallback:** Klien non-browser (seperti mobile app) juga dapat mengirimkan token via header:
> ```http
> Authorization: Bearer <access_token>
> ```

---

## 2. Format Standar Respons

### Respons Berhasil (Success)
```json
{
  "status": "success",
  "code": 200,
  "message": "Pesan deskriptif keberhasilan",
  "data": { ... }
}
```

### Respons Gagal (Error)
```json
{
  "status": "error",
  "code": 400,
  "message": "Pesan deskriptif kesalahan",
  "errors": "Detail error atau null"
}
```

---

## 3. Daftar Endpoint

### Modul Auth (`/api/auth`)

Tabel ikhtisar endpoint autentikasi:

| Method | Endpoint | Akses | Header Wajib | Request Body | Status Sukses | Deskripsi |
| :---: | :--- | :---: | :--- | :---: | :---: | :--- |
| `GET` | `/api/auth/csrf` | Publik | - | - | `200` | Mengambil CSRF token & cookie `csrf_` |
| `POST` | `/api/auth/register` | Publik | `X-CSRF-Token` | `{ name, email, password }` | `201` | Registrasi user baru (default role: `user`) |
| `POST` | `/api/auth/login` | Publik | `X-CSRF-Token` | `{ email, password }` | `200` | Login & set cookie `access_token` & `refresh_token` |
| `POST` | `/api/auth/refresh` | Publik | `X-CSRF-Token` *(Cookie: `refresh_token`)* | - | `200` | Memperbarui `access_token` baru |
| `POST` | `/api/auth/logout` | Publik | `X-CSRF-Token` | - | `200` | Menghapus session cookie di client |

---

#### 1. Get CSRF Token
Mengambil token CSRF baru untuk digunakan pada request berikutnya.

- **Method**: `GET`
- **URL**: `/api/auth/csrf`
- **Akses**: Publik
- **Headers**: Tidak ada

**Respons Sukses (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "CSRF token berhasil diambil",
  "data": {
    "csrf_token": "d2a1941c-f74b-4601-87dc-711cf81825cc"
  }
}
```
*Catatan: Header response menyertakan `Set-Cookie: csrf_=...; Path=/; SameSite=Lax`.*

---

#### 2. Register
Mendaftarkan pengguna baru dengan role default `user`.

- **Method**: `POST`
- **URL**: `/api/auth/register`
- **Akses**: Publik
- **Headers**:
  - `Content-Type: application/json`
  - `X-CSRF-Token: <token>`

**Request Body:**
```json
{
  "name": "Budi Santoso",
  "email": "budi@example.com",
  "password": "password123"
}
```

**Validasi:**
- `name`: Wajib diisi.
- `email`: Wajib format email valid dan belum terdaftar.
- `password`: Minimal 6 karakter.

**Respons Sukses (201 Created):**
```json
{
  "status": "success",
  "code": 201,
  "message": "Registrasi berhasil",
  "data": {
    "public_id": "8f3b2075-e8d9-4b21-9562-b13c7bb61c6b",
    "name": "Budi Santoso",
    "email": "budi@example.com",
    "role": "user",
    "created_at": "2026-09-05T22:30:00Z",
    "updated_at": "2026-09-05T22:30:00Z"
  }
}
```

**Respons Error Contoh (400 Bad Request):**
```json
{
  "status": "error",
  "code": 400,
  "message": "email sudah digunakan",
  "errors": null
}
```

---

#### 3. Login
Melakukan autentikasi menggunakan email dan password.

- **Method**: `POST`
- **URL**: `/api/auth/login`
- **Akses**: Publik
- **Headers**:
  - `Content-Type: application/json`
  - `X-CSRF-Token: <token>`

**Request Body:**
```json
{
  "email": "admin@taskforge.com",
  "password": "password123"
}
```

**Respons Sukses (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Login berhasil",
  "data": {
    "public_id": "cc3b4487-901b-45b6-8c32-580de73a8aeb",
    "name": "Administrator",
    "email": "admin@taskforge.com",
    "role": "admin",
    "created_at": "2026-09-04T20:21:41Z",
    "updated_at": "2026-09-04T20:21:41Z"
  }
}
```
*Catatan: Server menyetel cookie `access_token` dan `refresh_token` (`HttpOnly`, `SameSite=Lax`).*

**Respons Error Contoh (401 Unauthorized):**
```json
{
  "status": "error",
  "code": 401,
  "message": "email atau password salah",
  "errors": null
}
```

---

#### 4. Refresh Token
Menerbitkan `access_token` baru ketika token lama telah kadaluarsa.

- **Method**: `POST`
- **URL**: `/api/auth/refresh`
- **Akses**: Publik (memerlukan cookie `refresh_token`)
- **Headers**:
  - `X-CSRF-Token: <token>`

**Request Body:** Tidak ada.

**Respons Sukses (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Token berhasil diperbarui",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Respons Error Contoh (401 Unauthorized):**
```json
{
  "status": "error",
  "code": 401,
  "message": "Refresh token tidak valid atau kadaluarsa",
  "errors": null
}
```

---

#### 5. Logout
Menghapus session cookie di client.

- **Method**: `POST`
- **URL**: `/api/auth/logout`
- **Akses**: Publik
- **Headers**:
  - `X-CSRF-Token: <token>`

**Request Body:** Tidak ada.

**Respons Sukses (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Logout berhasil",
  "data": null
}
```
*Catatan: Cookie `access_token` dan `refresh_token` dihapus dari browser.*

---

### Modul User (`/api/users`)

Seluruh endpoint di bawah modul User memerlukan autentikasi (`Authenticate`).

Tabel ikhtisar endpoint user:

| Method | Endpoint | Akses / Guard | Parameter | Request Body | Status Sukses | Deskripsi |
| :---: | :--- | :---: | :--- | :---: | :---: | :--- |
| `GET` | `/api/users/me` | Authenticated | - | - | `200` | Profil user yang sedang aktif login |
| `GET` | `/api/users` | Khusus `admin` | Query: `?page=1&limit=10` | - | `200` | Daftar seluruh user dengan paginasi |
| `GET` | `/api/users/:id` | Authenticated | Path: `:id` (UUID) | - | `200` | Detail user spesifik berdasarkan ID |
| `PUT` | `/api/users/:id` | Self / Admin | Path: `:id` (UUID) | `{ name?, email? }` | `200` | Memperbarui nama atau email akun |
| `DELETE` | `/api/users/:id` | Khusus `admin` | Path: `:id` (UUID) | - | `200` | Menghapus user (proteksi akun sendiri) |

---

#### 1. Get Current User Profile (Me)
Mendapatkan informasi profil dari user yang sedang login.

- **Method**: `GET`
- **URL**: `/api/users/me`
- **Akses**: User Terotentikasi (semua role)
- **Cookie**: `access_token` (atau Header `Authorization: Bearer <token>`)

**Respons Sukses (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Berhasil mendapatkan data profil",
  "data": {
    "public_id": "cc3b4487-901b-45b6-8c32-580de73a8aeb",
    "name": "Administrator",
    "email": "admin@taskforge.com",
    "role": "admin",
    "created_at": "2026-09-04T20:21:41Z",
    "updated_at": "2026-09-04T20:21:41Z"
  }
}
```

---

#### 2. Get All Users (Pagination)
Melihat daftar seluruh pengguna terdaftar dalam sistem.

- **Method**: `GET`
- **URL**: `/api/users`
- **Akses**: Khusus **`admin`**
- **Query Parameters**:
  - `page` (opsional, default: `1`): Nomor halaman.
  - `limit` (opsional, default: `10`, max: `100`): Jumlah data per halaman.
- **Contoh URL**: `/api/users?page=1&limit=10`

**Respons Sukses (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Berhasil mengambil data pengguna",
  "data": {
    "users": [
      {
        "public_id": "cc3b4487-901b-45b6-8c32-580de73a8aeb",
        "name": "Administrator",
        "email": "admin@taskforge.com",
        "role": "admin",
        "created_at": "2026-09-04T20:21:41Z",
        "updated_at": "2026-09-04T20:21:41Z"
      },
      {
        "public_id": "8f3b2075-e8d9-4b21-9562-b13c7bb61c6b",
        "name": "Budi Santoso",
        "email": "budi@example.com",
        "role": "user",
        "created_at": "2026-09-05T22:30:00Z",
        "updated_at": "2026-09-05T22:30:00Z"
      }
    ],
    "total_data": 2,
    "current_page": 1,
    "total_pages": 1,
    "limit": 10
  }
}
```

**Respons Error Contoh (403 Forbidden):**
```json
{
  "status": "error",
  "code": 403,
  "message": "Akses ditolak: Anda tidak memiliki izin untuk tindakan ini",
  "errors": null
}
```

---

#### 3. Get User By ID
Mengambil data detail pengguna tertentu berdasarkan UUID `public_id`.

- **Method**: `GET`
- **URL**: `/api/users/:id`
- **Akses**: User Terotentikasi
- **Path Parameters**:
  - `id`: UUID pengguna (misal: `8f3b2075-e8d9-4b21-9562-b13c7bb61c6b`)

**Respons Sukses (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Berhasil mendapatkan data user",
  "data": {
    "public_id": "8f3b2075-e8d9-4b21-9562-b13c7bb61c6b",
    "name": "Budi Santoso",
    "email": "budi@example.com",
    "role": "user",
    "created_at": "2026-09-05T22:30:00Z",
    "updated_at": "2026-09-05T22:30:00Z"
  }
}
```

---

#### 4. Update User Profile
Mengubah informasi nama atau email pengguna.

- **Method**: `PUT`
- **URL**: `/api/users/:id`
- **Akses**: Pemilik akun sendiri atau **`admin`**
- **Headers**:
  - `Content-Type: application/json`
  - `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `id`: UUID pengguna target

**Request Body:**
```json
{
  "name": "Budi Santoso Baru",
  "email": "budi.baru@example.com"
}
```

**Respons Sukses (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Berhasil memperbarui data user",
  "data": {
    "public_id": "8f3b2075-e8d9-4b21-9562-b13c7bb61c6b",
    "name": "Budi Santoso Baru",
    "email": "budi.baru@example.com",
    "role": "user",
    "created_at": "2026-09-05T22:30:00Z",
    "updated_at": "2026-09-05T23:15:00Z"
  }
}
```

**Respons Error Contoh (403 Forbidden jika mengubah akun orang lain):**
```json
{
  "status": "error",
  "code": 403,
  "message": "Akses ditolak: Anda tidak memiliki izin untuk mengubah data user lain",
  "errors": null
}
```

---

#### 5. Delete User
Menghapus pengguna dari database.

- **Method**: `DELETE`
- **URL**: `/api/users/:id`
- **Akses**: Khusus **`admin`**
- **Headers**:
  - `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `id`: UUID pengguna target

> **Pengaman Khusus:** Admin tidak dapat menghapus akun miliknya sendiri yang sedang login.

**Respons Sukses (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "User berhasil dihapus",
  "data": null
}
```

**Respons Error Contoh (400 Bad Request jika admin menghapus dirinya sendiri):**
```json
{
  "status": "error",
  "code": 400,
  "message": "Tidak dapat menghapus akun Anda sendiri yang sedang aktif",
  "errors": null
}
```

---

## 4. Matriks Hak Akses (Role-Based Access Control)

| Endpoint | Publik | Role `user` | Role `admin` |
| :--- | :---: | :---: | :---: |
| `GET /api/auth/csrf` | ✅ | ✅ | ✅ |
| `POST /api/auth/register` | ✅ | ✅ | ✅ |
| `POST /api/auth/login` | ✅ | ✅ | ✅ |
| `POST /api/auth/refresh` | ✅ | ✅ | ✅ |
| `POST /api/auth/logout` | ✅ | ✅ | ✅ |
| `GET /api/users/me` | ❌ | ✅ | ✅ |
| `GET /api/users` | ❌ | ❌ | ✅ |
| `GET /api/users/:id` | ❌ | ✅ | ✅ |
| `PUT /api/users/:id` | ❌ | ✅ *(hanya akun sendiri)* | ✅ *(semua akun)* |
| `DELETE /api/users/:id` | ❌ | ❌ | ✅ *(kecuali diri sendiri)* |
