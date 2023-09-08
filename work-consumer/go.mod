module github.com/niko-dunixi/golang-simple-ingestion-pipeline-template/work-consumer

go 1.21.0

require (
	github.com/google/wire v0.5.0
	github.com/mb-14/gomarkov v0.0.0-20210216094942-a5b484cc0243
	github.com/niko-dunixi/golang-simple-ingestion-pipeline-template/lib v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.30.0
	github.com/stretchr/testify v1.8.4
	gocloud.dev v0.34.0
	gocloud.dev/pubsub/rabbitpubsub v0.34.0
	golang.org/x/sync v0.3.0
)

require (
	github.com/aws/aws-sdk-go v1.44.314 // indirect
	github.com/aws/aws-sdk-go-v2 v1.20.0 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.18.32 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.13.31 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.7 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.37 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.31 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.38 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.31 // indirect
	github.com/aws/aws-sdk-go-v2/service/sns v1.21.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sqs v1.24.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.13.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.15.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.21.1 // indirect
	github.com/aws/smithy-go v1.14.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rabbitmq/amqp091-go v1.8.1 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.134.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230731193218-e0aa005b6bdf // indirect
	google.golang.org/grpc v1.57.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/niko-dunixi/golang-simple-ingestion-pipeline-template/lib => ../lib
