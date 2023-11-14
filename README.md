Bynar backend golang application
==================================================

## Directory Description
### pkg
Common functions for service use.

| Name | Description |
| :--- | :--- |
| checkout | Checkout payment replate api, include GenerateAuthToken generate payment api access token and credit card verify|
| config | Global config, DB config|
| db | Database connection util|
| email | Sendgrid email relate util|
| errors | Error const variable and error handle relate functions|
| gcs | Google cloud storage relate util|
| gip | Google identify platform relate util|
| handler | Http treegrid handler|
| i18n | I18n handler relate util|
| logger | Log hanler util |
| middleware | Token auth middleware util |
| render | Http render relate util |
| treegrid | Treegrid relate handler util|
| utils | Collection of commonly used functions|

### service
System service http api collections, it call pkg common functions in pkg directory.

| Name | Description |
| :--- | :--- |
| accounts | Account relate service code, include sign in, sign up, upload profile picture, treegrid api code|
| general_posting_setup | General posting setup service treegrid relate api code|
| invoices | Invoices service treegrid relate api code|
| main | Main run code|
| organizations | Organizations service treegrid relate api code|
| usergroups | Usergroups service treegrid relate api code|
| warehouses | Warehouses service treegrid relate api code|
| sales | Sales service treegrid relate api code|
| payments | Payments service treegrid relate api code|
| procurements | Procurements service treegrid relate api code|

## How to run
Before run should create .env file in service/main directory, you can copy .env.template as base template.

```shell
 $ cd service/main/
 $ go run main.go
```

## ENV keys

#### Check out api keys
Can get key values from Checkout.com.

| Key | Description |
| :--- | :--- |
| CHECKOUT_SANDBOX | Is checkout sandbox enviroment |
| CHECKOUT_CLIENT_ID | Checkout client id |
| CHECKOUT_CLIENT_SECRET | Checkout client sercret |
| CHECKOUT_PROCESS_CHANNEL_ID | Checkout process channel id |

#### google identify platform keys

| Key | Description |
| :--- | :--- |
| GOOGLE_APPLICATION_CREDENTIALS_JSON | FireBase Service accounts private json key value. https://console.firebase.google.com/project/xxx/settings/serviceaccounts/adminsdk?hl=en |
| GOOGLE_API_KEY | Google Api key|
| SIGNUP_REDIRECT_URL | Sign up redirect url. Ex: "http://localhost:3000/signup"|
| SIGNUP_CUSTOM_VERIFICATION_KEY | Sign up custom verfication key|
| SIGNIN_REDIRECT_URL | Sign in redirect url. Ex: "http://localhost:3000/signin"|

#### database config keys
| Key | Description |
| :--- | :--- |
| DB_CONN_KEY | Bynar database connction string. Ex: "root:123456@tcp(localhost:3306)/bynar"|
| DB_ACCOUNT_CONN_KEY | Bynar accounts manager database connction string.  Ex: "root:123456@tcp(localhost:3306)/accounts_manager"|

#### sendgrid config keys
| Key | Description |
| :--- | :--- |
| SENDGRID_API_KEY | Send grid api key|
| SENDGRID_FROM_NAME | Send email from user name|
| SENDGRID_FROM_ADDRESS | Send email address|
| SENDGRID_TO_NAME | Send to name|
| SENDGRID_REDIRECT_URL | Sendgrid redirect url|

#### google cloud storage keys
| Key | Description |
| :--- | :--- |
| GOOGLE_CLOUD_STORAGE_BUCKET | Google cloud storage bucket name|

#### Code runtime environment
| Key | Description |
| :--- | :--- |
| ENV | Code runtime environment settings. Default: development|

## How to deploy

When push code to master will trigger google cloud build, google cloud build will run by cloudbuild.yaml step by step

## What services and Api uses

### Signup API
| Url | Description |
| :--- | :--- |
| /signup  | User sign up api|
| /confirm-email   | Confirm user email api when sign up|
| /verify-card  | Verify credit card api when sign up|
| /create-user  | Create user info to database api|

### Signin API
| Url | Description |
| :--- | :--- |
| /signin-email  | Send email api when sign in|
| /signin  | User sign in api|

### User API
| Url | Description |
| :--- | :--- |
| /user  | Get user info|
| /user/:id  | Get user info by id|
| /upload  | Upload user photo api|
| /profile-image  | Delete user photo api|
| /update-user-language-preference | Update user language preference api|
| /update-user-theme-preference | Update user theme preference api|
| /update-profile| Update user profile api|
| /organization-account| Get organization account api|
| /update-organization-account| update organization account api|
| /delete-organization-account| delete organization account api|

### Cards API
| Url | Description |
| :--- | :--- |
| /cards/list  | Get credit card list|
| /cards/add  | add credit card|
| /cards/update  | update credit card|
| /cards/delete  | delete credit card|

### TreeGrid API
| Url | Description |
| :--- | :--- |
| /:service/data  | Get list data api|
| /:service/page  | Get page data api|
| /:service/upload  | Hanlde upload api|
| /:service/cell  | Update cell value api|

Current avaliable sevices are: 
 * organizations
 * user_list
 * general_posting_setup
 * user_groups
 * warehouses
 * sales
 * payments
 * procurements
 * languages

## How to config domain
* DNS A address set

A address set to ip 34.160.127.53

* Upload https Certificate

Network services -> Load balancing -> bynar-backend -> EDIT -> Frontend configuration -> Add Frontend IP and protocol -> Upload https certificate.