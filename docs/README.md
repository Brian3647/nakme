# API usage & documentation

## General usage

All routes in `/api` return a JSON response (either the specified one below or [`ApiError`](#common-types-in-requests-and-responses)), and all the requests should have a JSON body.

HTTP codes follow the standard (RFC 9110 _[mdn docs](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status)_), but `400 Bad Request`, `500 Internal Server Error` and `200 OK` might be used even if there's a more specific code for what the response is trying to convey.

Optional strings (marked with `// <- optional`) can be sent as empty strings (`""`).

## Common types in requests and responses

- `token` is a user's token, used for authentication. It's a JSON Web Token (JWT) string: https://jwt.io/.
- `id` is a user's ID, used for identification. It's a number (a big one, be careful when using smaller int sizes) that is passed as a string.
- `ApiError` (in responses) is returned when an error occurs. It has the following structure:

```go
{
    message string
}
```

If the response has any status code that isn't between 200 and 399, you can safely handle it as an error.

## Authentication

## Sign up

**POST** `/api/auth/signup`

An emaill will be sent to the user with a confirmation link. The user will be able to log in only after confirming their email.
Note the confirmation link will be valid for only 1 hour.

Body:

```go
{
    username string, // <- must be at least 1 character long, max 80
    password string, // <- must be at least 8 characters long, max 80
    email string // <- must be unique
}
```

Response:

```go
{
    ok string // <- always "true", it will return an error if something went wrong
}
```

### Confirm email

**POST** `/api/auth/confirm_email`

Body:

```go
{
    token string
}
```

Response:

```go
{
    ok string // <- always "true", it will return an error if the token is invalid
}
```

## Log in

**POST** `/api/auth/login`

Body:

```go
{
    email string,
    password string
}
```

Response:

```go
{
    username string,
    token string
}
```

## Check if token is valid (confirm identity)

**POST** `/api/auth/confirm_identity`

Body:

```go
{
    token string
}
```

Response:

```go
{
    ok string // <- always "true", it will return an error if the token is invalid
}
```
