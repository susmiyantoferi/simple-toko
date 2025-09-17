# Simple toko

Teknologi Yang Di Gunakan :
- ![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
- ![GORM](https://img.shields.io/badge/GORM-0D8ABC?style=for-the-badge&logo=go&logoColor=white)
- ![MySQL](https://img.shields.io/badge/MySQL-4479A1?style=for-the-badge&logo=mysql&logoColor=white)
- ![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
- ![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white)
- ![k6](https://img.shields.io/badge/k6-1B00FF?style=for-the-badge&logo=k6&logoColor=white)
- ![JWT](https://img.shields.io/badge/JWT-000000?style=for-the-badge&logo=json-web-tokens&logoColor=white)
- ![CI/CD GitHub Actions](https://img.shields.io/badge/CI/CD_GitHub_Actions-6f42c1?style=for-the-badge&logo=github-actions&logoColor=white)

Business flow :

![business-flow](./pict/business_flow.png)

Batasan :
- confirm order manual by admin
- confirm payment manual by admin
- payment manual 
- belum terkoneksi dengan payment gateway
- belum terkoneksi dengan ekspedisi pengiriman

ERD :

![erd-diagram](./pict/simpletoko_ERD.png)

## Fitur utama:
- **User Management :** Registrasi customer, admin, login, refresh token, dan autentikasi menggunakan JWT.
- **CRUD :** Product, inventory, order, address, user, payment.
- **Create order :** customer bisa memilih lebih dari satu barang, customer bisa memilih dan mengupdate address
- **Payment :** customer dapat mengupload bukti pembayaran, admin dapat melihat atau download
- **Confirm order & payment  :** confirm by admin only
- **RBAC :** customer hanya bisa melakukan create order, update address, upload payment, melihat product, melihat order. admin dapat full akses fitur
- **Report :** penjualan perbulan, product terlaris dan kurang laris

## Set up local :

- clone repo
- ubah nama file .env.sample jadi .env
- ubah isi file .env sesuai dengan nama database, username, password, host, port, redis address dan redis password masing masing
  ```bash
  DB_NAME=db_toko
  DB_USERNAME=root
  DB_PASSWORD=
  DB_PORT=3306
  DB_HOST=127.0.0.1
  PORT_APP=:8080

  REDIS_ADDRS=127.0.0.1:6379
  REDIS_PWD=
  ```
- un-comment code berikut di file config/db.go :
  ```bash
  // "github.com/joho/godotenv"
  // err := godotenv.Load()
  // if err != nil{
  // 	log.Fatal("error load env")
  // }
  ```
- jalankan di terminal go run main.go

  ```bash
  go run main.go

  ```

- URL Documentation API : https://documenter.getpostman.com/view/22397647/2sB3BLjnnp
