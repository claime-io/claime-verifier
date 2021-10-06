import {
  ApiKeySourceType,
  LambdaIntegration,
  MockIntegration,
  PassthroughBehavior,
  Resource,
  RestApi,
} from '@aws-cdk/aws-apigateway'
import { Code, Function, Runtime, Tracing } from '@aws-cdk/aws-lambda'
import * as cdk from '@aws-cdk/core'
import { resolve } from 'path'
import { basicPolicytStatements } from './discord'
import * as environment from './env'

export class RestApiStack extends cdk.Stack {
  constructor(
    scope: cdk.Construct,
    id: string,
    target: environment.Environments,
    props?: cdk.StackProps,
  ) {
    super(scope, id, props)
    const api = new RestApi(this, 'RestApi', {
      restApiName: environment.withEnvPrefix(target, 'verification-restapi'),
      apiKeySourceType: ApiKeySourceType.HEADER,
    })
    const func = discordFunction(
      this,
      this.region,
      this.account,
      'verify',
      target,
      api,
    )
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
  const func = new Function(scope, resource, {
    functionName: `${environment.withEnvPrefix(target, resource)}`,
    code: code(resource),
    handler: 'bin/main',
    timeout: cdk.Duration.minutes(1),
    runtime: Runtime.GO_1_X,
    environment: {
      AllowedOrigin: environment.valueOf(target).allowedOrigin,
    },
    tracing: Tracing.ACTIVE,
  })
  basicPolicytStatements(region, account, target).forEach((s) =>
    func.addToRolePolicy(s),
  )
  const rs = api.root.addResource('discord').addResource('verity')
  rs.addMethod('POST', new LambdaIntegration(func))
  addCorsOptions(rs, environment.valueOf(target).allowedOrigin)
  return func
}

const code = (dirname: string) => {
  return Code.fromAsset(
    resolve(`${__dirname}/../`, 'lib', 'functions', dirname, 'bin', 'main.zip'),
  )
}

function addCorsOptions(apiResource: Resource, allowedOrigin: string) {
  apiResource.addMethod(
    'OPTIONS',
    new MockIntegration({
      integrationResponses: [
        {
          statusCode: '200',
          responseParameters: {
            'method.response.header.Access-Control-Allow-Headers':
              "'Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,X-Amz-User-Agent'",
            'method.response.header.Access-Control-Allow-Origin': `'${allowedOrigin}'`,
            'method.response.header.Access-Control-Allow-Credentials':
              "'false'",
            'method.response.header.Access-Control-Allow-Methods':
              "'OPTIONS,GET,PUT,POST,DELETE'",
          },
        },
      ],
      passthroughBehavior: PassthroughBehavior.NEVER,
      requestTemplates: {
        'application/json': '{"statusCode": 200}',
      },
    }),
    {
      methodResponses: [
        {
          statusCode: '200',
          responseParameters: {
            'method.response.header.Access-Control-Allow-Headers': true,
            'method.response.header.Access-Control-Allow-Methods': true,
            'method.response.header.Access-Control-Allow-Credentials': true,
            'method.response.header.Access-Control-Allow-Origin': true,
          },
        },
      ],
    },
  )
}