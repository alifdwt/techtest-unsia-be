# Technical Test UNSIA – Quiz Backend API

Backend API untuk sistem kuis (LMS sederhana) dengan **hybrid grading**:

- **Pilihan ganda** dinilai otomatis
- **Essay** dinilai manual (oleh dosen / admin)

Project ini difokuskan pada **desain skema database, query SQL, dan alur API**, sesuai dengan kebutuhan technical test backend.

---

## ✨ Fitur Utama

- Multiple quiz
- Batas waktu pengerjaan (timer berbasis backend)
- Batas percobaan (attempt limit)
- Reuse attempt selama masih `in_progress`
- Autosave jawaban
- Hybrid grading (auto + manual)
- Swagger API documentation

---

## 🏗️ Arsitektur Singkat

**Stack:**

- Go 1.25
- Fiber
- PostgreSQL 16
- SQLC (type-safe query)
- golang-migrate
- Swagger (swaggo)

**Struktur layer:**

- `handler` → HTTP layer
- `service` → business logic
- `db/query` → SQL
- `internal/db` → hasil generate SQLC

---

## 🗄️ Database Schema (ERD)

ERD direpresentasikan menggunakan **dbdiagram.io (DBML)** agar relasi lebih jelas dan tidak bergantung pada tooling lokal.

![ERD](./img/unsia-backend-erd.png)

**Catatan desain penting:**

- `answers.selected_option_id` bersifat **kondisional** (NULL untuk essay)
- Tidak ada tabel `users`, karena user diasumsikan berasal dari sistem autentikasi eksternal
- Relasi dijaga dengan foreign key di level database

---

## 📦 Migrasi Database

Menjalankan PostgreSQL via Docker:

```bash
docker compose up -d
```

Menjalankan migration:

```bash
migrate \
  -path db/migration \
  -database "postgres://postgres:postgres@localhost:5433/quizdb?sslmode=disable" \
  up
```

Seed data akan membuat:

- 1 quiz
- beberapa soal pilihan ganda & essay

---

## 🚀 Menjalankan Server

```bash
go run ./cmd/server
```

Server akan berjalan di:

```
http://localhost:8080
```

---

## 📘 Swagger API Documentation

Swagger tersedia di:

```
http://localhost:8080/swagger/index.html
```

![Swagger](./img/unsia-backend-swagger.png)

---

## 🔁 Alur Sistem Quiz

### 1. Start Quiz

```http
POST /start
```

- Membuat attempt baru **atau** reuse attempt `in_progress`
- Mengembalikan soal + deadline
- Attempt dihitung terhadap `max_attempts`

---

### 2. Submit Jawaban (Autosave / Final)

```http
POST /submit
```

Digunakan untuk:

- Autosave jawaban pilihan ganda
- Autosave jawaban essay
- Final submit (`final = true`)

Aturan penting:

- MC → `selected_option_id`
- Essay → `essay_answer`
- Essay **tidak** auto-graded

---

### 3. Lihat Hasil Quiz

```http
GET /result?attempt_id=UUID
```

Kondisi hasil:

- `in_progress` → tidak bisa diakses
- `waiting_assessment` → nilai MC saja
- `graded` → nilai final

---

## 📝 Manual Grading (via SQL)

Penilaian essay dilakukan langsung melalui SQL (misalnya via DBeaver).

### Melihat essay yang belum dinilai

```sql
SELECT
  a.id AS answer_id,
  q.question_text,
  a.essay_answer,
  q.points AS max_points
FROM answers a
JOIN questions q ON q.id = a.question_id
JOIN quiz_attempts qa ON qa.id = a.attempt_id
WHERE q.type = 'essay'
  AND a.score IS NULL
  AND qa.status = 'waiting_assessment';
```

### Memberi nilai essay

```sql
UPDATE answers
SET score = 15
WHERE id = 'ANSWER_ID';
```

### Menyelesaikan attempt

```sql
UPDATE quiz_attempts
SET status = 'graded'
WHERE id = 'ATTEMPT_ID';
```

---

## 📊 Query Rekap Nilai (Tambahan Tugas SQL)

Contoh query rekap nilai siswa per quiz:

```sql
SELECT
  q.title,
  SUM(a.score) AS final_score,
  CASE
    WHEN SUM(a.score) >= 60 THEN 'LULUS'
    ELSE 'TIDAK LULUS'
  END AS status_kelulusan
FROM quiz_attempts qa
JOIN quizzes q ON q.id = qa.quiz_id
JOIN answers a ON a.attempt_id = qa.id
WHERE qa.user_id = 'USER_ID'
  AND qa.status = 'graded'
GROUP BY q.title;
```

---

## 🎯 Catatan Akhir

Project ini **tidak bertujuan membangun LMS lengkap**, melainkan:

- menunjukkan desain database yang konsisten
- query SQL yang realistis
- API backend yang defensif dan scalable

Segala percakapan saya dengan ChatGPT terkait project ini dapat diakses di link berikut:
https://chatgpt.com/share/697b5825-69e8-800d-872c-4a5fc9b1e111

Terima kasih 🙌
