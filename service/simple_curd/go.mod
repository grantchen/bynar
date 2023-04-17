module git-codecommit.eu-central-1.amazonaws.com/v1/repos/simplecurd

go 1.19

require (
	github.com/go-sql-driver/mysql v1.7.0
	github.com/joho/godotenv v1.5.1
)

require git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs v0.0.0

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs => ../../pkg
