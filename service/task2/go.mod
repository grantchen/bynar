module git-codecommit.eu-central-1.amazonaws.com/v1/repos/task2

go 1.19

require git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs v0.0.0

require (
	github.com/aws/aws-lambda-go v1.40.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
)

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs => ../../pkg
