# API Usage Guide
 ## REGISTER USER

### Step 1: Register a New User
```bash
# Define the request body
BODY='{"username": "John Doe", "email": "john@example.com", "password": "mangotree"}'

# Register the user
curl -d "$BODY" http://localhost:4000/api/v1/register/user
```

### Step 2: Activate the User **NOTE: the TOKEN VALUE is sent through email**
```bash
# Replace "TOKEN_VALUE" with the token sent via email
curl -X PUT -d '{"token": "TOKEN_VALUE"}' http://localhost:4000/api/v1/users/activated
```
 ## AUTHENTICATE THE USER

### Step 1: Request authentication by logging in to authentication endpoint

```bash
# Define the request body
BODY='{"email": "john@example.com", "password": "mangotree"}'

# Trigger the authentication endpoint
curl -d "$BODY" http://localhost:4000/api/v1/tokens/authentication
```

### Step 2: Use the "BEARER TOKEN"
 ```bash
# Replace "BEARER_TOKEN" with the token returned in the previous step
curl -i -H "Authorization: Bearer BEARER_TOKEN" http://localhost:4000/api/v1/healthcheck
```

### Step 3: fetch a specific user Information NOTE: not fetching password
``` bash
# Replace ":uid" with the user ID
curl -i http://localhost:4000/api/v1/users/1
```
 ## READING LIST SECTION 

### fetch all reading list
``` bash
curl -i http://localhost:4000/api/v1/lists
```

### fetch a specific list
``` bash
# Replace ":rl_id" with the reading list ID
curl -X GET http://localhost:4000/api/v1/lists/:rl_id

```

### create a new reading list
``` bash
# Define the request body
BODY='{"name":"Manga Section","description":"List of current reading manga", "created_by":1}'

# Create the list
curl -X POST -d "$BODY" http://localhost:4000/api/v1/lists
```

### Update a reading list
``` bash
# Define the request body
BODY='{"name":"Manga Selection"}'

# Update the list (replace "1" with the list ID)
curl -X PUT -d "$BODY" http://localhost:4000/api/v1/lists/1
```


### Delete a reading List
```bash
curl -X DELETE localhost:4000/api/v1/lists/:rl_id
```

### Add book to reading list
```bash

#replace book_id and status field with valid id(int) and status(string) values
#status values can only be: ' currently reading' or 'completed'

BODY='{
  "book_id":BOOK_ID, 
  "status":STATUS
}'


#replace :rl_id with the id or list book is being added to
curl -X POST localhost:4000/api/v1/lists/:rl_id/books -d
```

## USER SECTION

### View User Profile 

```bash
# Replace ":uid" with the user ID
curl -i http://localhost:4000/api/v1/users/:uid
```

 ## BOOK SECTION

### Inser New Book

```bash 
#declare body
BODY='{
  "title": "Advanced Programming in Go",
  "author": ["John Doe", "Jane Smith", "Alice Brown"],
  "isbn": 9781234567890,
  "publication_date": "2022-05-15T00:00:00Z",
  "genre": ["Programming", "Technology", "Computer Science"],
  "description": "A comprehensive guide to advanced programming concepts and techniques in Go."
}'

curl -X POST -d "$BODY" localhost:4000/api/v1/books
```

## Fetch all Books with Pagination

```bash 
curl -i localhost:4000/api/v1/books
```

## Fetch Book Using ID
 ```bash

 curl -i localhost:4000/api/v1/books/:b_id

 ```

 ## Update Book Using ID
```bash
BODY='{
  "title": "Advanced Programming USING GO",
  "author": ["John Doe", "Jane Smith", "Alice Brown"],
  "isbn": 9781234567890,
  "publication_date": "2022-05-15T00:00:00Z",
  "genre": ["Programming", "Technology", "Computer Science"],
  "description": "A comprehensive guide to advanced programming concepts and techniques in Go."
}'

#replace 'b_id' with book value
 curl -X PUT localhost:4000/api/v1/books/b_id -d "BODY"

 ```

## Delete a Book
```bash
#replace 'b_id' with book value
 curl -X DELETE localhost:4000/api/v1/books/b_id
```

## Search Book By author/title/genre
```bash
# replace TITLE/GENRE/AUTHOR with query parameter
curl -i localhost:4000/api/v1/book/search?title=TITLE

curl -i localhost:4000/api/v1/book/search?author=AUTHOR

curl -i localhost:4000/api/v1/book/search?genre=GENRE
```