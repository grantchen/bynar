module git-codecommit.eu-central-1.amazonaws.com/v1/repos/main

go 1.19

require (
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs v0.0.0
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/usergroups v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.4.0
)

require (
	cloud.google.com/go v0.110.7 // indirect
	cloud.google.com/go/compute v1.23.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/firestore v1.12.0 // indirect
	cloud.google.com/go/iam v1.1.2 // indirect
	cloud.google.com/go/longrunning v0.5.1 // indirect
	cloud.google.com/go/storage v1.32.0 // indirect
	firebase.google.com/go/v4 v4.12.0 // indirect
	github.com/MicahParks/keyfunc v1.9.0 // indirect
	github.com/aws/aws-lambda-go v1.41.0 // indirect
	github.com/aws/aws-sdk-go v1.44.262 // indirect
	github.com/aws/aws-secretsmanager-caching-go v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.5 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/oauth2 v0.11.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.138.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/appengine/v2 v2.0.5 // indirect
	google.golang.org/genproto v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/grpc v1.57.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs => ../../pkg

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments => ../payments

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers => ../transfers

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales => ../sales

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements => ../procurements

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations => ../organizations

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/usergroups => ../usergroups

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup => ../general_posting_setup

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts => ../accounts
