Friday, 2023/09/08
- Complete environment creation when user created. (done).

Thursday, 2023/09/07
- Use a custom verification process when registering. The hmac-sha1 algorithm is used (done).
- Complete signup, confirm-email, verify-card, create-user handlers (done).
- Write the sendgrid configuration to the configuration file (done).
- Complete signin-email, signin handlers (done).

Wednesday, 2023/09/06
- Refactor the gip package and add unit tests (done).
- Accounts service handler add gip and checkout providers for verify email and card.
- Refactor the checkout.com generate GenerateAuthToken and ValidateCard api and test successfully.
- When registering, use the identity platform to send an email containing a link, verification email (oobcode verification fails).

Tuesday, 2023/09/05
- implement some interfaces of google identity platform.
- implment some sign up process code
- add checkout GenerateAccessToken api code
- add send mail using sendgrid code

Wednesday, 2023/04/28
- Build a common handlers (http and lambda) for paths: /data /page /cell /upload.
- Refactor transfer and payment apps using new handlers.


Wednesday, 2023/04/26
- Move treegrid model from payments (/upload api) and transfer (/data and /page api), ready for make a solid structure with handler layer is reuseable.
- Write handlers layer in pkgs (in-progress)


Sunday, 2023/04/23
- Fix bug for refactored transfer:
  - null component
  - loop recursive

Thursday, 2023/04/19
- Refactor transfer service (done). This service now has a better structure:
	- 3 layers
	- Some function move to pkgs for reusing in the future.

Wenesday, 2023/04/19
- Refactor transfer service (in progress)

Monday, 2023/04/17
- Refactor cell_url and simple_curd done (3 layer and use components db from pkgs)
-
Thurday, 2023/04/13
- Integrating and refactoring cell_url app.
	+ using db component from pkgs app
- Integrating and refactoring (in progress) simple_curd app.
	+ reorganise folder structure
	+ reorganise code structure (3 layer: handler, service and repository)

