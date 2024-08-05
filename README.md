# go-kv
Write a [key-value database](https://en.wikipedia.org/wiki/Key%E2%80%93value_database) in Go, implementing the following HTTP endpoints:

## `GET /[key]`
This returns the value for the provided key. If the key does not exist, a 404 is returned.

## `PUT /[key]`
This sets the value for the provided key to the contents of the request body. If the key already exists, it is updated in-place.

## `DELETE /[key]`
This deletes the value for the provided key. If the key does not exist, a 404 is returned.

## `GET /`
This returns a list of all keys as a JSON array.

# Notes
* You will need to keep paralellism in mind - the server should not encounter any race conditions when serving requests in parallel.
* The database does not need to be persisted to disk, it can be in-memory only.
* Each endpoint should be unit tested, including tests for parallelism.

