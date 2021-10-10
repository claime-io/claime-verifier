import {
  DomainName,
  EndpointType,
  LambdaIntegration,
  MockIntegration,
  PassthroughBehavior,
  Resource,
  RestApi,
  SecurityPolicy,
} from '@aws-cdk/aws-apigateway'
import * as certificatemanager from '@aws-cdk/aws-certificatemanager'
import { Code, Function, Runtime, Tracing } from '@aws-cdk/aws-lambda'
import * as route53 from '@aws-cdk/aws-route53'
import * as alias from '@aws-cdk/aws-route53-targets'
import * as cdk from '@aws-cdk/core'
import { resolve } from 'path'
import { basicPolicytStatements } from './discord'
import * as environment from './env'
import { hostedZoneFromId } from './route53'

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
    })
    const func = discordFunction(
      this,
      this.region,
      this.account,
      'verify',
      target,
      api,
    )
    const customDomain = withCustomDomain(this, api, target)

    aRecord(this, target, customDomain)
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
    environment: environmentVariables(target),
    tracing: Tracing.ACTIVE,
  })
  basicPolicytStatements(region, account, target).forEach((s) =>
    func.addToRolePolicy(s),
  )
  const rs = api.root.addResource('verify')
  rs.addMethod('POST', new LambdaIntegration(func))
  addCorsOptions(rs, environment.valueOf(target).allowedOrigin)
  return func
}
export const environmentVariables = (target: environment.Environments) => {
  return {
    AllowedOrigin: environment.valueOf(target).allowedOrigin,
    EnvironmentId: target,
  }
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
const aRecord = (
  stack: RestApiStack,
  target: environment.Environments,
  customDomain: DomainName,
) => {
  new route53.ARecord(stack, 'RestApiARecord', {
    zone: hostedZoneFromId(stack, target),
    recordName: restApiDomainName(target),
    target: route53.RecordTarget.fromAlias(
      new alias.ApiGatewayDomain(customDomain),
    ),
  })
}
const restApiDomainName = (target: environment.Environments) => {
  return `discord.` + environment.valueOf(target).rootDomain
}

const withCustomDomain = (
  stack: cdk.Stack,
  api: RestApi,
  target: environment.Environments,
) => {
  const customDomain = api.addDomainName(
    environment.withEnvPrefix(target, 'domain'),
    customDomainProps(stack, target),
  )
  return customDomain
}
const customDomainProps = (
  stack: cdk.Stack,
  target: environment.Environments,
) => {
  return {
    domainName: restApiDomainName(target),
    certificate: certificatemanager.Certificate.fromCertificateArn(
      stack,
      'Cert',
      environment.valueOf(target).certificateArn,
    ),
    securityPolicy: SecurityPolicy.TLS_1_2,
    endpointType: EndpointType.REGIONAL,
  }
}
