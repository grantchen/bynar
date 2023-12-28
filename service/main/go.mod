module git-codecommit.eu-central-1.amazonaws.com/v1/repos/main

go 1.21.4

require (
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/cards v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/invoices v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/languages v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs v0.0.0
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/sites v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/usergroups v0.0.0-00010101000000-000000000000
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/warehouses v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.4.0
)

require (
	cloud.google.com/go v0.111.0 // indirect
	cloud.google.com/go/compute v1.23.3 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/firestore v1.14.0 // indirect
	cloud.google.com/go/iam v1.1.5 // indirect
	cloud.google.com/go/longrunning v0.5.4 // indirect
	cloud.google.com/go/storage v1.36.0 // indirect
	firebase.google.com/go/v4 v4.12.0 // indirect
	github.com/BurntSushi/toml v1.0.0 // indirect
	github.com/MicahParks/keyfunc v1.9.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/nicksnyder/go-i18n/v2 v2.2.2 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.46.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.46.1 // indirect
	go.opentelemetry.io/otel v1.21.0 // indirect
	go.opentelemetry.io/otel/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.21.0 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/oauth2 v0.15.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/api v0.154.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/appengine/v2 v2.0.5 // indirect
	google.golang.org/genproto v0.0.0-20231212172506-995d672761c0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20231212172506-995d672761c0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231212172506-995d672761c0 // indirect
	google.golang.org/grpc v1.60.1 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs => ../../pkg

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments => ../payments

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers => ../transfers

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales => ../sales

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements => ../procurements

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations => ../organizations

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/sites => ../sites

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/usergroups => ../usergroups

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup => ../general_posting_setup

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts => ../accounts

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/invoices => ../invoices

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/cards => ../cards

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/warehouses => ../warehouses

replace git-codecommit.eu-central-1.amazonaws.com/v1/repos/languages => ../languages
