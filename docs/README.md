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

If the response has any status code that isn't between 200 and 399, you can safely handle it as an error (see above). Any response with this body:

```go
{
    ok string
}
```

will always be a success (you can safely assert that `ok == "true"`). If any errors happen, the response will have a different body (see above).

## General Authorization

Use the `Authorization: Bearer <JWT TOKEN>` header. (don't actually put `<` and `>` in the header). JWT tokens are good because even though anyone can read them, only server generated tokens are valid due to their signature. The tokens passed have this payload:

```go
{
    iat int, // <- issued at (date of creation)
    name string, // <- username
    sub string, // <- user id
}
```

So that even if they're decoded, they don't contain any sensitive information. That being said, do NOT send tokens to other people, as they can still use them in the API.

## Miscelaneous

**GET** `/api/health` or `/api/health`

Response:

```go
{
    ok string // <- as said above, this is always "true"
}
```

**GET** `/api/`

Returns a list of all the routes of the server.

Response:

```go
[
    {
        name string,
        path string,
        method string
    }
]
```

## User management

### Sign up

**POST** `/api/auth/signup`

An email will be sent to the user with a confirmation link. The user will be able to log in only after confirming their email.
Note the confirmation link will be valid for only 1 hour. When the link expires, the user will have to sign up again. The database isn't actually changed until the user confirms their email.

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
    ok string
}
```

### Confirm email

**GET** `/api/auth/confirm_email?token=JWT_TOKEN&email=EMAIL`

Response:

```go
{
    ok string
}
```

### Log in

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

### Check if token is valid (confirm identity)

**POST** `/api/auth/confirm_identity`

No body required, but it needs to have [authorization](#general-authorization).

Response:

```go
{
    ok string
}
```

### Delete account

**POST** `/api/auth/delete`

Both the token and the email are required for this request. The email is sent via the body:

```go
{
    email string
}
```

And the token is sent via the [authorization](#general-authorization) header.

Response:

```go
{
    ok string
}
```
