# backend-chuba-tickitz

![badge golang](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![badge postgresql](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![badge redis](https://img.shields.io/badge/redis-%23DD0031.svg?&style=for-the-badge&logo=redis&logoColor=white)

Backend project for [frontend Chuba tizkitz](https://github.com/M16Yusuf/chuba-tickitz). This project about build ticketing app for movies on cinemas. Build using gin gonic framework as backend, PostgreSQL as database, and redis as cache sistem.

## 🔧 Tech Stack

- [Go](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Redis](https://redis.io/docs/latest/operate/oss_and_stack/install/archive/install-redis/install-redis-on-windows/)
- [JWT](https://github.com/golang-jwt/jwt)
- [argon2](https://pkg.go.dev/golang.org/x/crypto/argon2)
- [migrate](https://github.com/golang-migrate/migrate)
- [Docker](https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository)
- [Swagger for API docs](https://swagger.io/) + [Swaggo](https://github.com/swaggo/swag)

## 🗝️ Environment

```bash
# database
DB_USER=<your_database_user>
DB_USER_PASS=<your_database_password>
DB_NAME=<your_database_name
DB_HOST=<your_database_host>
DB_PORT=<your_database_port>

# JWT hash
JWT_SECRET=<your_secret_jwt>
JWT_ISSUER=<your_jwt_issuer>

# Redish
RDB_HOST=<your_redis_host>
RDB_PORT=<your_redis_port>
RDB_USER=<your_redis_user>
RDB_PWD=<your_redis_password>

```

## ⚙️ Installation

1. Clone the project

```sh
$ https://github.com/M16Yusuf/backend-chuba-tickitz.git
```

2. Navigate to project directory

```sh
$ cd backend-chuba-tickitz
```

3. Install dependencies

```sh
$ go mod tidy
```

4. Setup your [environment](##-environment)

5. Install [migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation) for DB migration

6. Do the DB Migration

```sh
$ migrate -database YOUR_DATABASE_URL -path ./db/migrations up
```

7. Run the project

```sh
$ go run ./cmd/main.go
```

## 🚧 API Documentation

| Method | Endpoint                 | Body                                                                                                                                               | Description                       |
| ------ | ------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------- |
| GET    | /img                     |                                                                                                                                                    | Static File                       |
| POST   | /auth                    | email:string, password:string                                                                                                                      | Login                             |
| POST   | /auth/new                | email:string, password:string                                                                                                                      | Register                          |
| DELETE | /auth                    | header: Authorization (token jwt)                                                                                                                  | Logout                            |
| GET    | /movies                  | page:integer, search:string, genres: array                                                                                                         | Filter movies by genres or title  |
| GET    | /movies/popular          | page:integer                                                                                                                                       | Get popular movies                |
| GET    | /movies/upcoming         | page:integer                                                                                                                                       | Get upcoming movies               |
| GET    | /movies/{movie_id}       | movie_id:integer                                                                                                                                   | Get detail movie by id movie      |
| GET    | /users/                  | header: Authorization (token jwt)                                                                                                                  | Get user data                     |
| PATCH  | /users/                  | header: Authorization (token jwt), first_name: string, last_name: string, phone: string                                                            | Update user data                  |
| PATCH  | /users/avatar            | header: Authorization (token jwt), avatar:file                                                                                                     | Update avatar/profile user        |
| PATCH  | /users/password          | header: Authorization (token jwt), password:string                                                                                                 | update new password user          |
| GET    | /histories               | header: Authorization (token jwt)                                                                                                                  | Get histories data from a user    |
| GET    | /seats/{schedule_id}     | header: Authorization (token jwt), schedule_id: integer                                                                                            | Get booked seat from schedule_id  |
| GET    | /schedules/{movieid}     | header: Authorization (token jwt), movie_id: integer                                                                                               | Get all schedule from a movie     |
| GET    | /order                   | header: Authorization (token jwt), schedule_id:integer, payment_id:integer, total_price:integer, seat: array                                       | make new order for a user         |
| GET    | /admin/movies            | header: Authorization (token jwt), page:integer                                                                                                    | Get all movie list (admin)        |
| POST   | /admin/movies            | header: Authorization (token jwt), title, poster_path, backdrop_path, overview, duration, []actors{actor_name}, director{name}, []genres{genre_id} | make new data movie               |
| DELETE | /admin/movies/{movie_id} | header: Authorization (token jwt), movie_id: integer                                                                                               | soft delete a movie from database |

## 📄 LICENSE

MIT License

Copyright (c) 2025 Muhammad Yusuf m16yusuf

## 📧 Contact Info

[![LinkedIn](https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/m16yusuf/)
[![Instagram](https://img.shields.io/badge/Instagram-E4405F?style=for-the-badge&logo=Instagram&logoColor=white)](https://www.instagram.com/M16Yusuf/)
[![Twitter](https://img.shields.io/badge/Twitter-0077b5?style=for-the-badge&logo=Twitter&logoColor=white)](https://twitter.com/M16Yusuf)
[![Facebook](https://img.shields.io/badge/Facebook-1877F2?style=for-the-badge&logo=facebook&logoColor=white)](https://facebook.com/m16yusuff)

## 🎯 Related Project

[Frontend chuba-tickitz](https://github.com/M16Yusuf/chuba-tickitz)

[backend chuba-tickitz](https://github.com/M16Yusuf/backend-chuba-tickitz)

[ERD chuba-tickitz](https://github.com/M16Yusuf/sql-chuba-tickitz)
