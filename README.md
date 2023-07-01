# First-API

This project provides a RESTful API that handles CRUD operations for user data. It's built in Go and uses the Gin framework and the GORM ORM tool. It is designed following the Repository pattern and uses JWT for authentication. The application also features data caching using Redis and rate limiting.

## Project Structure

The project is structured as follows:

- `api`: This directory contains the main API components including models, routes, services, middleware, and repository.

- `database`: This directory contains code for setting up and handling the database connection.

- `pkg`: This directory contains additional packages used in the project.

- `test`: This directory contains the unit tests for the project.

## Running the Project

To run this project locally, follow these steps:

1. Install Go: Ensure that you have Go installed on your machine.

2. Clone the repository: Clone this repository to your local machine.

3. Install dependencies: Navigate to the root directory of the project and run `go mod tidy` to install the necessary dependencies.

4. Set environment variables: Make sure to set up your environment variables in the `.env` file. This includes configuring the database connection, Redis server, and rate limiting settings.

5. Start Redis server: Ensure that you have Redis server running locally or update the code to connect to a remote Redis server.

6. Run the application: Run the application using the `go run` command.

## Environment Variables

To configure the application, you need to set the following environment variables in the `.env` file:

```plaintext
PORT=<port_number>
HOST=<host>
USERNAME=<username>
PASSWORD=<password>
DBNAME=<database_name>
AUTH_KEY=<authentication_key>
REDIS_HOST=<redis_host>
REDIS_PASSWORD=<redis_password>
USER_CACHE_DB=<user_cache_database_number>
RATELIMITER_CACHE_DB=<ratelimiter_cache_database_number>
USER_CACHE_EXPTIME=<user_cache_expiry_time>
RATELIMITER_CACHE_EXPTIME=<ratelimiter_cache_expiry_time>


Replace the placeholders with the actual values for your environment. Here's a brief explanation of each variable:

    PORT: The port number to run the application on.
    HOST: The host of the MySQL database.
    USERNAME: The username for the MySQL database connection.
    PASSWORD: The password for the MySQL database connection.
    DBNAME: The name of the MySQL database.
    AUTH_KEY: The secret key used for JWT authentication.
    REDIS_HOST: The host and port of the Redis server.
    REDIS_PASSWORD: The password for the Redis server (leave empty if not required).
    USER_CACHE_DB: The Redis database number for caching user data.
    RATELIMITER_CACHE_DB: The Redis database number for rate limiter data.
    USER_CACHE_EXPTIME: The expiry time in milliseconds for the user cache.
    RATELIMITER_CACHE_EXPTIME: The expiry time in milliseconds for the rate limiter cache.

Make sure to adjust these values according to your environment.

## Features

### User routes

The API provides the following user-related routes:

- GET `/v1/user`: Retrieve all users.

- POST `/v1/user`: Create a new user.

- PUT `/v1/user/:id`: Update a user.

- DELETE `/v1/user/:id`: Delete a user.

- GET `/v1/user/filter`: Retrieve a user based on certain filters.

- POST `/v1/user/login`: Login route to authenticate a user.

## API Documentation

Please refer to the [API documentation](https://documenter.getpostman.com/view/27547278/2s93sdZsPu) for detailed information about the available endpoints, request/response formats, and examples.

### Caching

The application uses Redis for caching user data. This helps improve performance by reducing the need to fetch data from the database repeatedly.

### Rate Limiting

The API implements rate limiting to prevent abuse and protect against excessive requests.

## Testing

The project includes a suite of unit tests, found in the `test` directory. To run these tests, navigate to the root directory of the project and run `go test ./...`.

## Future Improvements

The project is a simple demonstration of a RESTful API in Go. Possible future improvements include adding more endpoints, improving error handling, adding integration tests, and enhancing security features.

## Contributions

Contributions, issues, and feature requests are welcome! Feel free to check the [issues page](#) for any open tasks or create a new issue for suggestions.

## License

