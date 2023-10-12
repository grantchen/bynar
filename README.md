Bynar backend golang application
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
Before run should create .env file in service/main directory, you can copy .env.template as base template.

```shell
 $ cd service/main/
 $ go run main.go
```

## ENV keys

#### Check out api keys
Can get key values from Checkout.com.
* CHECKOUT_SANDBOX
Is checkout sandbox enviroment
* CHECKOUT_CLIENT_ID
Checkout client id
* CHECKOUT_CLIENT_SECRET
Checkout client sercret
* CHECKOUT_PROCESS_CHANNEL_ID
Checkout process channel id

#### google identify platform keys
* GOOGLE_APPLICATION_CREDENTIALS_JSON
FireBase Service accounts private json key value
https://console.firebase.google.com/project/xxx/settings/serviceaccounts/adminsdk?hl=en
* GOOGLE_API_KEY
Google Api key
* SIGNUP_REDIRECT_URL
Sign up redirect url. Ex: "http://localhost:3000/signup"
* SIGNUP_CUSTOM_VERIFICATION_KEY
Sign up custom verfication key
* SIGNIN_REDIRECT_URL
Sign in redirect url. Ex: "http://localhost:3000/signin"

#### database config keys
* DB_CONN_KEY
bynar database connction string. Ex: "root:123456@tcp(localhost:3306)/bynar"
* DB_ACCOUNT_CONN_KEY
bynar accounts manager database connction string.  Ex: "root:123456@tcp(localhost:3306)/accounts_manager"

#### sendgrid config keys
* SENDGRID_API_KEY
Send grid api key
* SENDGRID_FROM_NAME
Send email from user name
* SENDGRID_FROM_ADDRESS
Send email address
* SENDGRID_TO_NAME
Send to name
* SENDGRID_REDIRECT_URL
Sendgrid redirect url

#### google cloud storage keys
* GOOGLE_CLOUD_STORAGE_BUCKET
Google cloud storage bucket name

#### Code runtime environment
* ENV
Code runtime environment settings. Default: development

## How to deploy

when push code to master will trigger google cloud build, google cloud build will run by cloudbuild.yaml step by step

## What services and Api uses

### Signup API
* /signup 
user sign up api
* /confirm-email 
confirm user email api when sign up
* /verify-card
verify credit card api when sign up
* /create-user
create user info to database api

### Signin API
* /signin-email
send email api when sign in
* /signin
user sign in api

### User API
* /user
get user info
* /user/:id
get user info by id
* /upload
upload user photo api
* /profile-image
delete user photo api
* /update-user-language-preference
update user language preference api
* /update-user-theme-preference
update user theme preference api
* /update-profile
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