# Golang API Starter

##### Golang v.1.14

Install:

```
git clone https://github.com/judascrow/go-apix.git
rm .git
rm go.mod
rm go.sum
go mod init <New Mudule>
git init
cp .env.example .env
```

## Features

- RESTful API
- Gin
- Gorm
- Gin-Swagger
- Log-Zap
- Gin-Jwt
- Casbin

## COMMAND

Migrate & Seed :

```
go run main.go create seed
```

Run :

```
go run main.go
```

## API ROUTER URI

```
Host: localhost:8000
Base Path: /api/v1
```

### Healthcheck

| Done               | Method | URI          | authorize | comment              |
| ------------------ | ------ | ------------ | --------- | -------------------- |
| :white_check_mark: | GET    | /healthcheck | No        | Check Status Service |

### Auth

| Done               | Method | URI    | authorize | comment |
| ------------------ | ------ | ------ | --------- | ------- |
| :white_check_mark: | POST   | /login | No        | Log in  |

### Users

| Done               | Method | URI                   | authorize | comment          |
| ------------------ | ------ | --------------------- | --------- | ---------------- |
| :white_check_mark: | GET    | /users                | Yes       | List Users       |
| :white_check_mark: | GET    | /users/:slug          | Yes       | Get User by Slug |
| :white_check_mark: | POST   | /users                | Yes       | Create User      |
| :white_check_mark: | PUT    | /users/:slug          | Yes       | Update User      |
| :white_check_mark: | PUT    | /users/:slug/password | Yes       | Change Password  |
| :white_check_mark: | DELETE | /users/:slug          | Yes       | Delete User      |

### Swaggo

```
Generate command: swag init

URL: http://localhost:8000/api/v1/swagger/index.html
```

## FOR DEV

### GIT

Update

```
./git.sh "COMMENT"
```

### Project Structure

```
.

│
├───api
│   │   server.go
│   │
│   ├───controllers
│   │       users.go
│   │
│   ├───infrastructure
│   │       db.go
│   │
│   ├───middlewares
│   │   └───jwt
│   │           auth_jwt.go
│   │
│   ├───models
│   │       base.go
│   │       casbin.go
│   │       role.go
│   │       user.go
│   │
│   ├───routes
│   │       auth.go
│   │       router.go
│   │
│   ├───seeds
│   │       seeder.go
│   │
│   ├───services
│   │       shared.go
│   │       users.go
│   │
│   └───utils
│       ├───messages
│       │       errors.go
│       │       messages.go
│       │
│       └───responses
│               responses.go

│   .env
│   .env.example
│   .gitignore
│   auth.conf
│   git.sh
│   go.mod
│   go.sum
│   LICENSE
│   main.go
│   README.md
```
