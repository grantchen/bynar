Welcome to Bynar backend application
==================================================

## Directory Description
### pkg
 - checkout - checkout payment replate api, include GenerateAuthToken generate payment api access token and credit card verify
 - config - global config, DB config
 - db - database connection util
 - email - email relate util
 - errors - error const variable and error handle relate functions
 - gcs - google cloud storage relate util
 - gip - google identify platform relate util
 - handler - http treegrid handler
 - i18n - i18n handler relate util
 - logger - log hanler util
 - middleware - token auth middleware util
 - render - http render relate util
 - treegrid - treegrid relate handler util
 - utils - Collection of commonly used functions

### service
 - accounts - account relate service code, include sign in, sign up, upload profile picture, treegrid api code
 - general_posting_setup - general posting setup service treegrid relate api code
 - invoices - invoices service treegrid relate api code
 - main - main run code
 - organizations - organizations service treegrid relate api code
 - usergroups - usergroups service treegrid relate api code

## How to run

```shell
 $ cd service/main/
 $ go run main.go
```

## How to deploy

when push code to master will trigger google cloud build, google cloud build will run by cloudbuild.yaml step by step


## What services and Api uses

### Signup API
- /signup 
user sign up api
- /confirm-email 
confirm user email api when sign up
- /verify-card
verify credit card api when sign up
- /create-user
create user info to database api

### Signin API
- /signin-email
send email api when sign in
- /signin
user sign in api

### User API
- /user
get user info
- /user/:id
get user info by id
- /upload
upload user photo api
- /profile-image
delete user photo api
- /update-user-language-preference
update user language preference api
- /update-user-theme-preference
update user theme preference api
- /update-profile
update user profile api

### TreeGrid API
- /:service/data
get list data api
- /:service/page
get page data api
- /:service/upload
hanlde upload api
- /:service/cell
update cell value api

current support sevices are: organizations, user_list, general_posting_setup, user_groups