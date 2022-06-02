import {
  DomainName,
  EndpointType,
  IResource,
  MockIntegration,
  PassthroughBehavior,
  RestApi,
  SecurityPolicy,
} from '@aws-cdk/aws-apigateway'
import * as certificatemanager from '@aws-cdk/aws-certificatemanager'
import { Code } from '@aws-cdk/aws-lambda'
import * as route53 from '@aws-cdk/aws-route53'
import * as alias from '@aws-cdk/aws-route53-targets'
import * as cdk from '@aws-cdk/core'
import { Stack } from '@aws-cdk/core'
import { resolve } from 'path'
import * as environment from './env'
import { hostedZoneFromId } from './route53'

type CordOptionsConfig = {
  methods: ('GET' | 'POST' | 'PUT' | 'DELETE')[]
  withCredentials?: true
  headers?: string[]
  origin?: string
}
export function addCorsOptions(
  apiResource: IResource,
  config: CordOptionsConfig,
) {
  const responseParameters = integrationResponseParameters(config)
  const methodResponseParameters = Object.keys(responseParameters).reduce(
    (res, key) => Object.assign(res, { [key]: true }),
    {},
  )
  apiResource.addMethod(
    'OPTIONS',
    new MockIntegration({
      integrationResponses: [{ statusCode: '200', responseParameters }],
      passthroughBehavior: PassthroughBehavior.NEVER,
      requestTemplates: { 'application/json': '{"statusCode": 200}' },
    }),
    {
      methodResponses: [
        { statusCode: '200', responseParameters: methodResponseParameters },
      ],
    },
  )
}

const integrationResponseParameters = (config: CordOptionsConfig) => {
  const { methods, withCredentials, headers, origin } = config
  const params = {
    'method.response.header.Access-Control-Allow-Headers': `'${
      headers?.join(',') || '*'
    }'`,
    'method.response.header.Access-Control-Allow-Origin': `'${origin || '*'}'`,
    'method.response.header.Access-Control-Allow-Methods': `'OPTIONS,${methods.join(
      ',',
    )}'`,
  }
  if (withCredentials) {
    Object.assign(params, {
      'method.response.header.Access-Control-Allow-Credentials': "'true'",
    })
  }
  return params
}

export const environmentVariables = (target: environment.Environments) => {
  return {
    AllowedOrigin: environment.valueOf(target).allowedOrigin,
    EnvironmentId: target,
    SubgraphEndpoint: environment.valueOf(target).subgraphEndpoint,
  }
}
const aRecord = (
  stack: Stack,
  target: environment.Environments,
  domainName: string,
  customDomain: DomainName,
) => {
  new route53.ARecord(stack, 'RestApiARecord', {
    zone: hostedZoneFromId(stack, target),
    recordName: domainName,
    target: route53.RecordTarget.fromAlias(
      new alias.ApiGatewayDomain(customDomain),
    ),
  })
}

export const code = (dirname: string) => {
  return Code.fromAsset(
    resolve(`${__dirname}/../`, 'lib', 'functions', dirname, 'bin', 'main.zip'),
  )
}
export const withCustomDomain = (
  stack: cdk.Stack,
  api: RestApi,
  domain: string,
  target: environment.Environments,
) => {
  const customDomain = api.addDomainName(
    environment.withEnvPrefix(target, 'domain'),
    customDomainProps(stack, domain, target),
  )
  aRecord(stack, target, domain, customDomain)
  return customDomain
}

const customDomainProps = (
  stack: cdk.Stack,
  domain: string,
  target: environment.Environments,
) => {
  return {
    domainName: domain,
    certificate: certificatemanager.Certificate.fromCertificateArn(
      stack,
      'Cert',
      environment.valueOf(target).certificateArn,
    ),
    securityPolicy: SecurityPolicy.TLS_1_2,
    endpointType: EndpointType.REGIONAL,
  }
}
