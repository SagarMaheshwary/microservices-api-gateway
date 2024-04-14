# MICROSERVICES - API GATEWAY

API Gateway for Microservices project.

### TECHNOLOGIES

- Golang (1.22)
- Gin framework
- gRPC

### SETUP

After cloning the project, cd into the project directory and copy **.env.example** to **.env** and update the required variables.

Create executable and start the server:

```bash
go build cmd/server/main.go && ./main
```

Or install "[air](https://github.com/cosmtrek/air)" and run it to autoreload when making file changes:

```bash
air -c .air-toml
```

### APIs

Checkout Postman collection and environment files for examples in **api/postman** directory.

| API            | METHOD | BODY                                                        | Headers                                | Description                                                                                                                 |
| -------------- | ------ | ----------------------------------------------------------- | -------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| /auth/register | POST   | {"name": "string", "email": "string", "password", "string"} | -                                      | User registration via [authentication service](https://github.com/SagarMaheshwary/microservices-authentication-service) RPC |
| /auth/login    | POST   | { "email": "string", "password", "string"}                  | -                                      | User login via authentication service RPC                                                                                   |
| /auth/logout   | POST   | -                                                           | Bearer token in "authorization" header | User logout via authentication service RPC                                                                                  |
