# REST API

This specification is based on the LXD specification. See [https://github.com/lxc/lxd/blob/master/doc/rest-api.md](https://github.com/lxc/lxd/blob/master/doc/rest-api.md) for more information.

## API Path

The default path to the API is `/api/<version>`.

## API Versioning

The currently supported api versions are:

- `v1`

## Return Values

### Standard Result Object

The API is guaranteed to return an object for all requests unless explicitly defined. An example is shown below.

```json
{
    "type": "sync",
    "status": "Success", // "Success", "Error"
    "status_code": 200, // e.g., 400
    "metadata": {
        // Other result data
    }
}
```

### Type

The only currently supported type are synchronous commands, so `"sync"` will always be returned.

### Status

Goshelf uses two `status` keys:

- `"Success"`
- `"Failure"`


### Status Code

The HTTP status code returned is also returned in the response body as per LXD specification. Goshelf uses the following status codes:

Code | Meaning
--- | ---
200 | Success
400 | Failure

## Standard HTTP Methods

The bookshelf API supports GET, PUT, POST, and DELETE HTTP methods.
