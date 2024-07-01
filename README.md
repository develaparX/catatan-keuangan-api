# Live Code Keuangan


Organisasi XYZ akan mencatat semua pengeluaran mereka, baik itu uang masuk maupun uang keluar. Berikut contoh data
pengeluaran:

|    Tanggal      |      Masuk      |      Keluar     |      Saldo      | Tipe Transaksi  |    Keterangan   |
| --------------- | --------------- | --------------- | --------------- | --------------- | --------------- |
| 1 Dec 2023      |      5,000,000  |                 |       5,000,000 |   KREDIT        | Saldo awal bulan desember 2023    |
| 1 Dec 2023      |                 |       1,000,000 |       4,000,000 |   DEBIT         | Bayar kontrakan kantor    |
| 2 Dec 2023      |        500,000  |                 |       4,500,000 |   KREDIT        | Dapet bonus    |

Buatlah program untuk mencatat pengeluaran seperti tabel di atas.

Requirement:
1. Buat service catatan keuangan.
2. Buatlah sebuah api login menggunakan middleware Bearer auth, dimana response dari api login berupa token yang akan digunakan sebagai authorization api transaksi
3. Buat sebuah DML untuk insert data user pertama
4. Password di dalam table user harus dalam bentuk encryption
5. Berikan sebuah informasi apabila pengeluaran melebihi saldo yang tersedia.
6. Buatlah service sebagai berikut:
   - `CREATE` pengeluaran
   - `LIST` Pengeluaran (pagination dan atau bisa filter dengan range tanggal)
   - `GET` by ID
   - `GET` berdasarkan tipe transaksi
7. Jangan lupa buatkan Dokumentasi untuk menjalankan atau bisa EXPORT Collection nya.
8. EXPORT Database atau sertakan DDL dan DML nya didalam yang sudah ada isian data pada Tabel.


Pengumpulan:
1. Buat project baru di `livecode` dengan nama `livecode-catatan-keuangan`.
2. Pesan dalam Commit harus jelas.

Berikut DDL nya:
```sql
CREATE DATABASE IF NOT EXISTS catatan_keuangan_db;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE transaction_type AS ENUM ('CREDIT', 'DEBIT');

CREATE TABLE expenses (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    date DATE NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    transaction_type transaction_type,
    balance DOUBLE PRECISION NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE users (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    email varchar(255) NOT NULL,
    password varchar(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
```

Berikut endpoint, request dan response:
1. `LOGIN`
   - `POST` -> `/api/v1/users/login`
   - `HEADER` -> `Authorization` 'Basic base64(clientId +:+ client secret)'
   - Request:
     ```json
     {
        "email": "test@gmail.com",
        "password": "test"
     }
     ``` 
   - Response:
     ```json
     "responseCode":"2000101"
     "data": {
        "token": "string"
     }
     ```
3. `CREATE` pengeluaran
   - `POST` -> `/api/v1/expenses`
   - `HEADER` -> `Authorization` 'Bearer token'
   - Request:
     ```json
     {
        "amount": 250000,
        "transactionType": "CREDIT",
        "description": "Tambahan jajan"
     }
     ``` 
   - Response:
     ```json
     "responseCode":"2000101"
     "data": {
        "id": "a81bc81b-dead-4e5d-abff-90865d1e13b1",
        "date": "2023-12-08 05:17:42.583767+07",
        "amount": 250000,
        "transactionType": "CREDIT",
        "description": "Tambahan jajan",
        "createdAt": "2023-12-08 05:17:42.583767+07",
        "updatedAt": "2023-12-08 05:17:42.583767+07"
     }
     ``` 
3. `LIST` Pengeluaran (pagination dan atau bisa filter dengan range tanggal)
   - `GET` -> `/api/v1/expenses`
   - `HEADER` -> `Authorization` 'Bearer token'
   - Query Params:
     ```
     ?page=1&size=5
     ?page=1&size=5&startDate=2023-12-08&endDate=2023-12-08
     ``` 
   - Response:
     ```json
     "responseCode": "2000101",
     "data": [
       {
          "id": "a81bc81b-dead-4e5d-abff-90865d1e13b1",
          "date": "2023-12-08 05:17:42.583767+07",
          "amount": 250000,
          "transactionType": "CREDIT",
          "balance": 50000000,
          "description": "Tambahan jajan",
          "createdAt": "2023-12-08 05:17:42.583767+07",
          "updatedAt": "2023-12-08 05:17:42.583767+07"
       }
     ],
     "paging": {
        "page": 1,
        "totalData": 10
     }
     ``` 
5. `GET` by ID
   - `GET` -> `/api/v1/expenses/a81bc81b-dead-4e5d-abff-90865d1e13b1`
   - `HEADER` -> `Authorization` 'Bearer token'
   - Response:
     ```json
     "responseCode": "2000101",
     "data": {
        "id": "a81bc81b-dead-4e5d-abff-90865d1e13b1",
        "date": "2023-12-08 05:17:42.583767+07",
        "amount": 250000,
        "transactionType": "CREDIT",
        "balance": 50000000,
        "description": "Tambahan jajan",
        "createdAt": "2023-12-08 05:17:42.583767+07",
        "updatedAt": "2023-12-08 05:17:42.583767+07"
     }
7. `GET` berdasarkan tipe transaksi
   - `GET` -> `/api/v1/expenses/type/:type`
   - `HEADER` -> `Authorization` 'Bearer token'
   - Params:
     ```
     /credit
     /debit
     ``` 
   - Query Params:
     ```
     ?page=1&size=5
     ``` 
   - Response:
     ```json
     "responseCode": "2000101",
     "data": [
       {
          "id": "a81bc81b-dead-4e5d-abff-90865d1e13b1",
          "date": "2023-12-08 05:17:42.583767+07",
          "amount": 250000,
          "transactionType": "CREDIT",
          "balance": 50000000,
          "description": "Tambahan jajan",
          "createdAt": "2023-12-08 05:17:42.583767+07",
          "updatedAt": "2023-12-08 05:17:42.583767+07"
       }
     ],
     "paging": {
        "page": 1,
        "totalData": 10
     }
     ``` 
