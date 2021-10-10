import {
  ApiKeySourceType,
  LambdaIntegration,
  RequestValidator,
  RestApi,
} from '@aws-cdk/aws-apigateway'
import { SubnetType, Vpc } from '@aws-cdk/aws-ec2'
import {
  Cluster,
  ContainerImage,
  FargateService,
  FargateTaskDefinition,
  LogDriver,
} from '@aws-cdk/aws-ecs'
import {
  Effect,
  ManagedPolicy,
  PolicyStatement,
  Role,
  ServicePrincipal,
} from '@aws-cdk/aws-iam'
import { Code, Function, Runtime, Tracing } from '@aws-cdk/aws-lambda'
import { LogGroup } from '@aws-cdk/aws-logs'
import { Secret } from '@aws-cdk/aws-secretsmanager'
import * as cdk from '@aws-cdk/core'
import { RemovalPolicy } from '@aws-cdk/core'
import { resolve } from 'path'
import { dataSourceReadWritePolicyStatement } from './datasource'
import * as environment from './env'
import { environmentVariables } from './restapi'

export class DiscordStack extends cdk.Stack {
  constructor(
    scope: cdk.Construct,
    id: string,
    target: environment.Environments,
    props?: cdk.StackProps,
  ) {
    super(scope, id, props)
    const api = new RestApi(this, 'RestApi', {
      restApiName: environment.withEnvPrefix(target, 'restapi'),
      apiKeySourceType: ApiKeySourceType.HEADER,
    })
    const discordAPISecrets = new Secret(
      this,
      environment.withEnvPrefix(target, 'discord-api-key'),
    )
    const func = discordFunction(
      this,
      this.region,
      this.account,
      'discord',
      target,
      api,
    )
    discordAPISecrets.grantRead(func)
    const vpc = new Vpc(this, 'vpc', {
      cidr: '192.156.0.0/16',
      natGateways: 1,
      subnetConfiguration: [
        {
          cidrMask: 24,
          name: 'Public1',
          subnetType: SubnetType.PUBLIC,
        },
        {
          cidrMask: 24,
          name: 'Public2',
          subnetType: SubnetType.PUBLIC,
        },
        {
          cidrMask: 24,
          name: 'Private1',
          subnetType: SubnetType.PRIVATE,
        },
        {
          cidrMask: 24,
          name: 'Private2',
          subnetType: SubnetType.PRIVATE,
        },
      ],
    })
    const executionRole = new Role(this, 'EcsTaskExecutionRole', {
      roleName: environment.withEnvPrefix(target, 'ecs-execution-role'),
      assumedBy: new ServicePrincipal('ecs-tasks.amazonaws.com'),
      managedPolicies: [
        ManagedPolicy.fromAwsManagedPolicyName(
          'service-role/AmazonECSTaskExecutionRolePolicy',
        ),
      ],
    })
    const serviceTaskRole = new Role(this, 'EcsServiceTaskRole', {
      roleName: environment.withEnvPrefix(target, 'ecs-service-task-role'),
      assumedBy: new ServicePrincipal('ecs-tasks.amazonaws.com'),
    })
    basicPolicytStatements(this.region, this.account, target).forEach((s) =>
      serviceTaskRole.addToPolicy(s),
    )
    const logGroup = new LogGroup(this, 'ServiceLogGroup', {
      logGroupName: '/aws/ecs/claime-cluster' + target,
      removalPolicy: RemovalPolicy.DESTROY,
    })
    const image = ContainerImage.fromAsset('.')
    const cpu = 256
    const mem = 1024
    const taskDef = new FargateTaskDefinition(this, 'ServiceTaskDefinition', {
      cpu: cpu,
      memoryLimitMiB: mem,
      executionRole: executionRole,
      taskRole: serviceTaskRole,
    })
    taskDef.addContainer('ContainerDef', {
      image,
      cpu: cpu,
      memoryLimitMiB: mem,
      logging: LogDriver.awsLogs({
        streamPrefix: 'claime',
        logGroup,
      }),
      environment: environmentVariables(target),
    })
    const cluster = new Cluster(this, 'claime-cluster', {
      clusterName: environment.withEnvPrefix(target, 'cluster'),
      containerInsights: true,
      vpc: vpc,
      enableFargateCapacityProviders: true,
    })
    const fargateService = new FargateService(this, 'FargateService', {
      cluster,
      vpcSubnets: vpc.selectSubnets({ subnetType: SubnetType.PRIVATE }),
      taskDefinition: taskDef,
      desiredCount: 1,
    })
  }
}

const discordFunction = (
  scope: cdk.Construct,
  region: string,
  account: string,
  resource: string,
  target: environment.Environments,
  api: RestApi,
) => {
  const requestValidator = fromDiscordRequestValidator(scope, 'validator', api)
  const func = new Function(scope, resource, {
    functionName: `${environment.withEnvPrefix(target, resource)}`,
    code: code(resource),
    handler: 'bin/main',
    timeout: cdk.Duration.minutes(1),
    runtime: Runtime.GO_1_X,
    tracing: Tracing.ACTIVE,
    environment: environmentVariables(target),
  })
  basicPolicytStatements(region, account, target).forEach((s) =>
    func.addToRolePolicy(s),
  )
  api.root.addMethod(
    'POST',
    new LambdaIntegration(func, fromDiscordIntegrationOptions()),
    discordMethodIntegration(requestValidator),
  )
  return func
}

const discordMethodIntegration = (validator: RequestValidator) => {
  return {
    requestValidator: validator,
    apiKeyRequired: false,
    methodResponses: [
      {
        statusCode: '200',
      },
      {
        statusCode: '401',
      },
    ],
  }
}

const fromDiscordRequestValidator = (
  scope: cdk.Construct,
  id: string,
  api: RestApi,
) => {
  return new RequestValidator(scope, id, {
    restApi: api,
    validateRequestBody: true,
    validateRequestParameters: true,
  })
}

const fromDiscordIntegrationOptions = () => {
  return {
    proxy: false,
    requestTemplates: {
      'application/json':
        '{\r\n\
          "timestamp": "$input.params(\'x-signature-timestamp\')",\r\n\
          "signature": "$input.params(\'x-signature-ed25519\')",\r\n\
          "jsonBody" : $input.json(\'$\')\r\n\
        }',
    },
    integrationResponses: [
      {
        statusCode: '200',
      },
      {
        statusCode: '401',
        selectionPattern: '.*[UNAUTHORIZED].*',
        responseTemplates: {
          'application/json': 'invalid request signature',
        },
      },
    ],
  }
}

export const basicPolicytStatements = (
  region: string,
  account: string,
  target: environment.Environments,
) => {
  return [
    dataSourceReadWritePolicyStatement(region, account, target),
    new PolicyStatement({
      effect: Effect.ALLOW,
      actions: ['ssm:Get*'],
      resources: [
        `arn:aws:ssm:${region}:${account}:parameter/claime-verifier*`,
      ],
    }),
  ]
}

const code = (dirname: string) => {
  return Code.fromAsset(
    resolve(`${__dirname}/../`, 'lib', 'functions', dirname, 'bin', 'main.zip'),
  )
}
