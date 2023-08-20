export DOCKER_DEFAULT_PLATFORM := "linux/amd64"

fmt: tidy
  just --unstable lib/fmt
  just --unstable work-supplier/fmt
  just --unstable work-consumer/fmt
  just --unstable infrastructure/fmt

tidy:
  just --unstable lib/tidy
  just --unstable work-supplier/tidy
  just --unstable work-consumer/tidy
  just --unstable infrastructure/tidy

generate:
  just --unstable lib/generate
  just --unstable work-supplier/generate
  just --unstable work-consumer/generate
  just --unstable infrastructure/generate

test: generate
  just --unstable lib/test
  just --unstable work-supplier/test
  just --unstable work-consumer/test
  just --unstable infrastructure/test

update:
  just --unstable lib/update
  just --unstable work-supplier/update
  just --unstable work-consumer/update
  just --unstable infrastructure/update

vulnerability-check:
  just --unstable lib/vulnerability-check
  just --unstable work-supplier/vulnerability-check
  just --unstable work-consumer/vulnerability-check
  just --unstable infrastructure/vulnerability-check

clean:
  git clean -Xdf
  just --unstable lib/clean-go
  just --unstable work-supplier/clean-go
  just --unstable work-consumer/clean-go
  just --unstable infrastructure/clean-go
  docker-compose down
  docker volume prune --all --force

synth:
  just --unstable infrastructure/synth

deploy: test
  just --unstable infrastructure/deploy

local:
  docker compose up --build -d
  docker compose logs --follow

check-ingress-reachable:
  curl $(aws cloudformation describe-stacks --stack SimpleIngestionPipelineStack | jq -r '.Stacks[].Outputs[] | select(.ExportName == "IngestionURL") | .OutputValue')

push-ingestion-event:
  curl -XPOST $(aws cloudformation describe-stacks --stack SimpleIngestionPipelineStack | jq -r '.Stacks[].Outputs[] | select(.OutputKey = "ExportIngestionURL") | .OutputValue')task/foobar

diff:
  just --unstable infrastructure/diff

destroy:
  just --unstable infrastructure/destroy

bootstrap:
  just --unstable infrastructure/bootstrap