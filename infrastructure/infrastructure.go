//go:build !generate

package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecrassets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type InfrastructureStackProps struct {
	awscdk.StackProps
}

func SimpleIngestionPipelineStack(scope constructs.Construct, id string, props *InfrastructureStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	// This is pulled in from the cdk.context.json file automatically.
	// You will need to run 'go generate ./...' to create it
	mainVPC := awsec2.Vpc_FromLookup(stack, jsii.String("MainVpc"), &awsec2.VpcLookupOptions{
		VpcId: jsii.String(stack.Node().GetContext(jsii.String("vpc-id")).(string)),
	})

	epoch := fmt.Sprintf("%d", time.Now().Unix())

	queue := awssqs.NewQueue(stack, jsii.String("IngestionSQS"), &awssqs.QueueProps{
		QueueName:         jsii.String("ingestion-sqs"),
		VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number[float64](300)),
	})

	stack.ExportValue(queue.QueueUrl(), &awscdk.ExportValueOptions{
		Name: jsii.String("QueueURL"),
	})

	// Work Supplying Function
	lambdaPrincipal := awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), &awsiam.ServicePrincipalOpts{})
	workSupplierRole := awsiam.NewRole(stack, jsii.String("WorkSupplierRole"), &awsiam.RoleProps{
		AssumedBy: lambdaPrincipal,
	})
	workSupplierRole.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: jsii.Strings(
			"logs:CreateLogGroup",
			"logs:CreateLogStream",
			"logs:PutLogEvents",
			"logs:DescribeLogStreams",
		),
		Resources: jsii.Strings(
			"arn:aws:logs:*:*:*",
		),
	}))
	queue.GrantSendMessages(workSupplierRole)
	// policyStatement := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
	// 	Actions:   jsii.Strings("sqs:SendMessage"),
	// 	Resources: jsii.Strings(*queue.QueueArn()),
	// })
	// awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
	// 	Actions: jsii.Strings("lambda:InvokeFunctionUrl"),
	// 	Resources: ,
	// })
	workSupplierDockerImage := awslambda.DockerImageCode_FromImageAsset(
		jsii.String(path.Join("..")),
		&awslambda.AssetImageCodeProps{
			AssetName: jsii.String("WorkSupplierContainerImage"),
			Target:    jsii.String("main-lambda"),
			Platform:  awsecrassets.Platform_LINUX_AMD64(),
			ExtraHash: jsii.String(epoch),
			Invalidation: &awsecrassets.DockerImageAssetInvalidationOptions{
				ExtraHash: jsii.Bool(true),
			},
			BuildArgs: &map[string]*string{
				"TARGET_PACKAGE": jsii.String("work-supplier"),
				"WIRE_TAGS":      jsii.String("aws"),
			},
		})
	workSupplierFunction := awslambda.NewDockerImageFunction(stack, jsii.String("WorkSupplierDockerImageFunction"), &awslambda.DockerImageFunctionProps{
		FunctionName: jsii.String("WorkSupplier"),
		Environment: &map[string]*string{
			"QUEUE_URL": queue.QueueUrl(),
		},
		Code:         workSupplierDockerImage,
		Role:         workSupplierRole,
		LogRetention: awslogs.RetentionDays_FIVE_DAYS,
	})
	workSupplierFunctionURL := workSupplierFunction.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_NONE,
		// Cors: &awslambda.FunctionUrlCorsOptions{
		// 	AllowedMethods: &[]awslambda.HttpMethod{awslambda.HttpMethod_POST},
		// },
	})
	workSupplierFunctionURL.GrantInvokeUrl(awsiam.NewAnyPrincipal())
	stack.ExportValue(workSupplierFunctionURL.Url(), &awscdk.ExportValueOptions{
		Name: jsii.String("IngestionURL"),
	})
	// Work Consuming Fargate Task
	workConsumerDockerImage := awsecs.AssetImage_FromDockerImageAsset(awsecrassets.NewDockerImageAsset(stack,
		jsii.String("WorkConsumerDockerImageAsset"),
		&awsecrassets.DockerImageAssetProps{
			Directory: jsii.String(".."),
			AssetName: jsii.String("WorkConsumerContainerImage"),
			Target:    jsii.String("main-vanilla"),
			Platform:  awsecrassets.Platform_LINUX_AMD64(),
			ExtraHash: jsii.String(epoch),
			Invalidation: &awsecrassets.DockerImageAssetInvalidationOptions{
				ExtraHash: jsii.Bool(true),
			},
			BuildArgs: &map[string]*string{
				"TARGET_PACKAGE": jsii.String("work-consumer"),
				"WIRE_TAGS":      jsii.String("aws"),
			},
		}))
	// ecsPrincipal := awsiam.NewServicePrincipal(jsii.String("ecs.amazonaws.com"), &awsiam.ServicePrincipalOpts{})
	ecsTaskPrincipal := awsiam.NewServicePrincipal(jsii.String("ecs-tasks.amazonaws.com"), &awsiam.ServicePrincipalOpts{})
	workConsumerExecutionRole := awsiam.NewRole(stack, jsii.String("WorkConsumerExecutionRole"), &awsiam.RoleProps{
		AssumedBy: ecsTaskPrincipal,
	})
	workConsumerTaskRole := awsiam.NewRole(stack, jsii.String("WorkConsumerTaskRole"), &awsiam.RoleProps{
		AssumedBy: ecsTaskPrincipal,
	})
	queue.GrantConsumeMessages(workConsumerTaskRole)
	fargateCluster := awsecs.NewCluster(stack, jsii.String("SimpleIngestionPipelineCluster"), &awsecs.ClusterProps{
		ClusterName: jsii.String("SimpleIngestionPipelineCluster"),
		Vpc:         mainVPC,
	})
	taskDefinition := awsecs.NewTaskDefinition(stack, jsii.String("WorkConsumerTaskDefinition"), &awsecs.TaskDefinitionProps{
		ExecutionRole: workConsumerExecutionRole,
		TaskRole:      workConsumerTaskRole,
		Cpu:           jsii.String("256"),
		MemoryMiB:     jsii.String("512"),
		Compatibility: awsecs.Compatibility("FARGATE"),
	})
	taskDefinition.AddContainer(jsii.String("WorkConsumerTaskContainer"), &awsecs.ContainerDefinitionOptions{
		Image: workConsumerDockerImage,
		Environment: &map[string]*string{
			"QUEUE_URL": queue.QueueUrl(),
		},
		Logging: awsecs.NewAwsLogDriver(&awsecs.AwsLogDriverProps{
			StreamPrefix: jsii.String("TaskContainerInstance"),
			Mode:         awsecs.AwsLogDriverMode_NON_BLOCKING,
			LogRetention: awslogs.RetentionDays_FIVE_DAYS,
		}),
	})
	awsecs.NewFargateService(stack, jsii.String("WorkConsumerService"), &awsecs.FargateServiceProps{
		Cluster: fargateCluster,
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
		},
		TaskDefinition: taskDefinition,
		DesiredCount:   jsii.Number(1),
	})
	// consumerService.AutoScaleTaskCount(&awsapplicationautoscaling.EnableScalingProps{
	// 	MinCapacity: jsii.Number(0),
	// 	MaxCapacity: jsii.Number(1),
	// }).ScaleOnMetric(jsii.String("SqsQueueVisibilityCount"), &awsapplicationautoscaling.BasicStepScalingPolicyProps{
	// 	Metric: queue.MetricApproximateNumberOfMessagesVisible(&awscloudwatch.MetricOptions{
	// 		// Period of one minute
	// 		Period: awscdk.Duration_Millis(jsii.Number(60 * 1000 * 1)),
	// 	}),
	// 	AdjustmentType: awsapplicationautoscaling.AdjustmentType_EXACT_CAPACITY,
	// 	// Cooldown after 2 minutes
	// 	Cooldown: awscdk.Duration_Millis(jsii.Number(60 * 1000 * 2)),
	// 	ScalingSteps: &[]*awsapplicationautoscaling.ScalingInterval{
	// 		{
	// 			// If there are no messages in queue, scale to zero
	// 			Lower:  jsii.Number(0),
	// 			Upper:  jsii.Number(1),
	// 			Change: jsii.Number(0),
	// 		},
	// 		{
	// 			// Any messages, scale to one
	// 			Lower:  jsii.Number(1),
	// 			Change: jsii.Number(1),
	// 		},
	// 	},
	// 	EvaluationPeriods: jsii.Number(1),
	// })

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	SimpleIngestionPipelineStack(app, "SimpleIngestionPipelineStack", &InfrastructureStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	// return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
