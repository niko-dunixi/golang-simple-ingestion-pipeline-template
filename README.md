# Simple Ingestion Pipeline Template

A staple of distributed systems: this is a relatively simple
implementation of an API fronting and ingestion pipeline that
performs asyncronous work for any given request.

It can be run either on
  * AWS deployed with the [AWS CDK](https://aws.amazon.com/cdk/)
  * Locally via the [Go Cloud SDK](https://gocloud.dev/) and Docker

## Cut to the chase, is it good?
🦄 Hell Yeah! 🚀
This is the template *I* want every time this problem comes up. It addresses:
* Infrastructure
* Dependency Injection
  * Wiring to either real resources or local ones

## Problem

You need to perform a non-trivial amount of work upon user request.
Often I find that I want something:
  * Ephemeral and Lambda-like
    * YET
  * Work for a long period of time without interruption

This generally boils-down to Lambda dumping a payload into SQS, which a
AWS Fargate Task can pick up and perform the work without being interrupted.

```mermaid
flowchart LR
  client[Client Application]
  subgraph AWS
    lambda[AWS Lambda]
    sqs[AWS SQS]
    lambda -- "Pushes event payloads" --> sqs
    subgraph Your VPC With Private Subnets
      worker[AWS Fargate - Long Running Task]
    end
    sqs -- "Pulls event payloads" --> worker
  end
  client -- "Performs HTTP Request via Function URL" --> lambda
```

The largest barrior to implementing such a solution is often:
 * Implementing it in a timely manor
   * Every time I've had to reinvent this at a new organisation,
     I have been unable to reference previous details and this results
     in restarting a lot of the detail work from scratch
 * Lack of completeness in off-the-shelf options/samples
   * Many examples online leave out details and don't get past the
     most basic 'Hello, World!' implementation which don't address
     larger questions.
   * How do I tie in mature dependency injection into the solution?
      * How do I actually refactor common types and functions between the
        ingestion Lambda and the consuming Fargate worker?
      * How do I expand this or incorporate it into my current architecture?
        How does this get from my machine into a target environment?

## Dependencies
You will need to install the following for your machine

### Running Justfile commands
- [Just](https://github.com/casey/just)
  - A language agnostic command runner


### Building and Testing
- [Docker](https://docs.docker.com/desktop/)
  - You will need to follow the Docker Desktop installation instructions
    for your specific OS.
  - When deploying to AWS, the AWS CDK will build your containers locally first
  - When running locally, docker compose will build your stack and run them
- [GoLang](https://go.dev/doc/install)
- [Wire](https://github.com/google/wire)
  - Compile time dependency injection. Used to switch AWS and local implementations

### Deploying to AWS via [AWS CDK](https://aws.amazon.com/cdk/)

> The AWS CDK is analagous to Terraform except being declarative,
> your Infrastructure-as-Code is done with a supported language of
> your choosing. In this case, it is GoLang to keep the project
> homogenous and lessen the cognative-load needed to understand
> this repository.

- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
  - You must configure your machine for development with AWS if you want to deploy this to AWS
- NodeJS
  - I personally find [the fnm project](https://github.com/Schniz/fnm) to be reliable and portable node version manager
  - Alternatively [nvm](https://github.com/nvm-sh/nvm) is popular

You will need a VPC with Private (not just Isolated) Subnets. If you do not have
one, you can deploy a sample I've written:
 - [simple-private-vpc-template](https://github.com/niko-dunixi/simple-private-vpc-template)

## Running, Testing, Deploying

The fastest way to get things up and running:

## Running locally

`$ just local`
 * This will build the stack with the dependencies from [`docker-compose.yaml`](./docker-compose.yaml)
 * It is faster to interate changes locally, this is the recommended way to get started

## Running in AWS

`$ just bootstrap`
 * This is necessary the first time you use the CDK with your AWS account
 * It is only needed once, you don't need to run it ever again

`$ just deploy`
 * Will synthesize the CloudFormation template and deploy it to your account

`$ just destroy`
 * Will tear down and delete all the resources created when you deployed
 * Be sure to do this when you no longer need your VPC, the VPC Endpoints will incur costs


## House keeping

`$ just tidy`
 * This will run `go mod tidy` upon all sub directories

`$ just fmt`
 * Will run `go fmt` upon all sub directories

## FAQ

### Why is the `go.work` file commited?

While it has been discouraged at times, it is not
cannonically or authoritatively incorrect to do.
It makes working with multi-module monorepo projects,
such as this one, vastly simpler to do. The discussion
is still open and if it proves incorrect iterative
course correct will occur.
- [proposal: ref/mod: mention whether go.work files should be checked into VCS #53502](https://github.com/golang/go/issues/53502#issuecomment-1204134618)