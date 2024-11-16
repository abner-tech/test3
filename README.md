# Books-Api

#BEGISTER USER 

//declaring the body with user beign signed
BODY='{"username": "John Doe", "email": "john@example.com", "password": "mangotree"}'

//actual sign in
curl -d "$BODY" localhost:4000/api/v1/register/user

//activating the user....NOTE: the TOKEN VALUE is sent through email
curl -X PUT -d '{"token": "TOKEN VALUE"}' localhost:4000/api/v1/users/activated




//AUTHENTICATE USER

//declaring the body for our request NOTE: replace values with credentials
BODY='{ "email": "john@example.com", "password": "mangotree"}'

//triger authentication endpoint: this will return a token, and the expiry date for its use
curl -d "$BODY" localhost:4000/api/v1/tokens/authentication

//send a request along with the bearer token that is returned
curl -i -H "Authorization: Bearer BEARER_TOKEN" localhost:4000/api/v1/healthcheck
//current authorization token for user1: LNHATDTSKLKFVVNTKC4YQ5IFIY

//get user information printes, note that password wont be printed 
 curl -i  http://localhost:4000/api/v1/users/1



//READING LIST SECTION

//show all lists
curl -i localhost:4000/api/v1/lists

//show a specific list
curl -i localhost:4000/api/v1/lists/:rl_id

//making a list for the existing user with id 1
BODY='{"name":"Manga Section","description":"list of current reading manga", "created_by":1}'
curl -X POST -d "$BODY" localhost:4000/api/v1/lists


