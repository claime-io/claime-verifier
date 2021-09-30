import {
  ApiKeySourceType,
  LambdaIntegration,
  RequestValidator,
  RestApi,
} from '@aws-cdk/aws-apigateway'
import { Effect, PolicyStatement } from '@aws-cdk/aws-iam'
import { Code, Function, Runtime, Tracing } from '@aws-cdk/aws-lambda'
import { Secret } from '@aws-cdk/aws-secretsmanager'
import * as cdk from '@aws-cdk/core'
import { resolve } from 'path'
import { dataSourceReadWritePolicyStatement } from './datasource'
import * as environment from './env'

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
