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

| Done               | Method | URI         | authorize | comment |
| ------------------ | ------ | ----------- | --------- | ------- |
| :white_check_mark: | POST   | /auth/login | No        | Log in  |

### Users

| Done               | Method | URI                   | authorize | comment          |
| ------------------ | ------ | --------------------- | --------- | ---------------- |
| :white_check_mark: | GET    | /users                | Yes       | List Users       |
| :white_check_mark: | GET    | /users/:slug          | Yes       | Get User by Slug |
| :white_check_mark: | POST   | /users                | Yes       | Create User      |
| :white_check_mark: | PUT    | /users/:slug          | Yes       | Update User      |
| :white_check_mark: | DELETE | /users/:slug          | Yes       | Delete User      |
| :white_check_mark: | PUT    | /users/:slug/password | Yes       | Change Password  |
| :white_check_mark: | PUT    | /users/:slug/avatar   | Yes       | Upload Avatar    |

### Swaggo

```
Generate command: swag init

URL: http://localhost:8000/api/v1/swagger/index.html
```

### File Server

```

URL: http://localhost:8000/media
```

## FOR DEV

### GIT

Update

```
./git.sh "COMMENT"
```
