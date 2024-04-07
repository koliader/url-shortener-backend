# Next.js URL Shortener Backend

This is the backend for a URL shortener application built with Golang, Gin, SQLC, and PostgreSQL.

## Features

- **Authentication:**
  - Supports both simple authentication (username and password) and GitHub OAuth.
- **URL Shortening:**
  - Users can shorten URLs anonymously or after authentication.
- **Profile:**
  - Users can view their own shortened URLs and total clicks for each URL.
  - Users can update the click count for each URL.

## Endpoints

- **Authentication:**
  - `POST /auth/register`: Register a new user.
  - `POST /auth/login`: Log in with username and password.
  - `GET /auth/github/:code`: Authenticate with GitHub OAuth.
- **URLs:**
  - `POST /urls/guest`: Shorten a URL anonymously.
  - `GET /urls/:code`: Get the original URL associated with the provided code.
  - `PUT /urls/:code`: Update the click count for a URL.
  - `POST /urls`: Shorten a URL after authentication.
  - `GET /urls/myUrls`: Get a list of shortened URLs created by the authenticated user.
- **Users:**
  - `GET /users/:username`: Get user profile by username.

## License

This project is licensed under the [MIT License](LICENSE).
