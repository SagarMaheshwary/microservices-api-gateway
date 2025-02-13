# MICROSERVICES - API GATEWAY

API Gateway for the [Microservices](https://github.com/SagarMaheshwary/microservices) project.

### OVERVIEW

- Golang
- ZeroLog
- Gin framework - REST API
- gRPC – Client implementations for Authentication, Upload, and Video Catalog services
- Prometheus Client – Exports default and custom metrics for Prometheus server monitoring

### SETUP

Follow the instructions in the [README](https://github.com/SagarMaheshwary/microservices?tab=readme-ov-file#setup) of the main microservices repository to run this service along with others using Docker Compose.

### APIs

Check out the Postman collection and environment files in the **api/postman** directory for example requests.

| API                          | METHOD | BODY                                                                                                                                                                                                               | Headers                                | Description                                                                                                                                          |
| ---------------------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| /auth/register               | POST   | {"name": "string", "email": "string", "password", "string"}                                                                                                                                                        | -                                      | User registration - [authentication service](https://github.com/SagarMaheshwary/microservices-authentication-service)                                |
| /auth/login                  | POST   | { "email": "string", "password", "string"}                                                                                                                                                                         | -                                      | User login - authentication service                                                                                                                  |
| /auth/logout                 | POST   | -                                                                                                                                                                                                                  | Bearer token in "authorization" header | User logout - authentication service                                                                                                                 |
| /auth/profile                | GET    | -                                                                                                                                                                                                                  | Bearer token in "authorization" header | Get currently logged in user - authentication service                                                                                                |
| /videos                      | GET    | -                                                                                                                                                                                                                  | -                                      | List videos - [video catalog service](https://github.com/SagarMaheshwary/microservices-video-catalog-service)                                        |
| /videos/:id                  | GET    | -                                                                                                                                                                                                                  | -                                      | Get specified video details as well as DASH manifest url from cloudfront for streaming that video - video catalog service                            |
| /videos/upload/presigned-url | POST   | -                                                                                                                                                                                                                  | Bearer token in "authorization" header | Get S3 presigned url for uploading a video from frontend/postman - [upload service](https://github.com/SagarMaheshwary/microservices-upload-service) |
| /videos/upload/webhook       | POST   | {"video_id": "string - s3 upload id from presigned-url process", "thumbnail_id": "string - s3 upload id from presigned-url process", "title": "string - video title", "description": "string - video description"} | Bearer token in "authorization" header | Create a video - upload service                                                                                                                      |
| /health                      | GET    | -                                                                                                                                                                                                                  | -                                      | Service healthcheck endpoint                                                                                                                         |
| /metrics                     | GET    | -                                                                                                                                                                                                                  | -                                      | Prometheus metrics endpoint                                                                                                                          |
