# Golang API Example

##### Golang v.1.14

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
| :white_check_mark: | GET    | /healthcheck | No        | ครวจสอบสถานะ Service |

### Auth

| Done               | Method | URI    | authorize | comment |
| ------------------ | ------ | ------ | --------- | ------- |
| :white_check_mark: | POST   | /login | No        | Log in  |

### Users

| Done               | Method | URI    | authorize | comment     |
| ------------------ | ------ | ------ | --------- | ----------- |
| :white_check_mark: | GET    | /users | Yes       | List Users  |
| :white_check_mark: | POST   | /users | Yes       | Create User |

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
