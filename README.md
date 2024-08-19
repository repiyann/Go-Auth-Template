# Go Fiber Auth Template using GORM

## Description

This project is a comprehensive authentication template built with Go Fiber and GORM, designed to provide a robust foundation for developing secure and scalable applications. The template integrates essential components such as JWT for authentication, Mailjet for email services, and PostgreSQL for database management.

### Build With

* [Go Fiber](https://gofiber.io/)
* [GORM](https://gorm.io/index.html)
* [Mailjet](https://www.mailjet.com/)
* [Go JWT](https://github.com/golang-jwt/jwt)
* [PostgreSQL](https://www.postgresql.org/)
  
# Getting Started

## Dependencies

You need to have Go installed on your machine. You can download it from the [Go official site](https://go.dev/).

Additionally, you'll need to set up a Mailjet account for email functionalities.

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/yourusername/your-repository.git
    cd your-repository
    ```
2. Install Go dependencies:
    ```bash
    go mod tidy
    ```
3. Create a copy of your .env file
    ```bash
    cp .env.example .env
    ```
4. Edit the .env file to include your database information, Mailjet credentials, and other environment variables.
5. Set up your database and run migrations. Make sure you have the database set up according to your .env configuration.

### Mailjet Configuration

1. Ensure you have an account with Mailjet and have your API keys.
2. Add these keys to your .env file
    ```bash
    MJ_APIKEY_PUBLIC=your-mailjet-api-key
    MJ_APIKEY_PRIVATE=your-mailjet-api-secret
    ```

### JWT Configuration

1. Generate a secret key for JWT and add it to your .env file
    ```bash
    JWT_SECRET=your-secret-key
    ```

## Executing Program

1. Run the Go Fiber server:
    ```bash
    go run app/cmd/main.go
    ```
2. The port is configurable via the .env file. By default, it is set to 8080. You can interact with the API using tools like Postman or curl.

# Author

* GitHub: [@repiyann](https://github.com/repiyann)
* Instagram: [@repiyann](https://instagram.com/repiyann)
