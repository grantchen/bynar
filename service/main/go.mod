module git-codecommit.eu-central-1.amazonaws.com/v1/repos/main

go 1.19

require (
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments v0.0.0
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs v0.0.0
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements v0.0.0
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales v0.0.0
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers v0.0.0
)

require (
	github.com/aws/aws-lambda-go v1.40.0 // indirect
	github.com/aws/aws-sdk-go v1.44.237 // indirect
	github.com/aws/aws-secretsmanager-caching-go v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs => ../../pkg

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments => ../payments

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers => ../transfers

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales => ../sales

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements => ../procurements
