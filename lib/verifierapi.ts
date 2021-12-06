import { LambdaIntegration, RestApi } from '@aws-cdk/aws-apigateway'
import { Function, Runtime, Tracing } from '@aws-cdk/aws-lambda'
import { IHostedZone } from '@aws-cdk/aws-route53'
import * as cdk from '@aws-cdk/core'
import {
  addCorsOptions,
  code,
  environmentVariables,
  withCustomDomain,
} from './api'
import { basicPolicytStatements } from './discord'
import * as environment from './env'

type VerifierApiStackProps = {
  hostedZone: IHostedZone
}
export class VerifierApiStack extends cdk.Stack {
  constructor(
    scope: cdk.Construct,
    id: string,
    target: environment.Environments,
    props: VerifierApiStackProps & cdk.StackProps,
  ) {
    super(scope, id, props)
    const { hostedZone } = props
    const api = new RestApi(this, 'RestApi', {
      restApiName: environment.withEnvPrefix(target, 'verifier'),
    })
    apifunction(this, 'verify', ['verify', '{eoa}'], target, api)
    apifunction(
      this,
      'testVerification',
      ['test', 'verify', '{eoa}'],
      target,
      api,
    )
    withCustomDomain(this, api, restApiDomainName(target), target, hostedZone)
  }
}

const apifunction = (
  stack: cdk.Stack,
  resource: string,
  paths: string[],
  target: environment.Environments,
  api: RestApi,
) => {
  const func = new Function(stack, resource, {
    functionName: `${environment.withEnvPrefix(target, resource)}`,
    code: code(resource),
    handler: 'bin/main',
    timeout: cdk.Duration.minutes(1),
    runtime: Runtime.GO_1_X,
    environment: environmentVariables(target),
    tracing: Tracing.ACTIVE,
  })
  basicPolicytStatements(stack.region, stack.account, target).forEach((s) =>
    func.addToRolePolicy(s),
  )
  const re = paths.reduce(
    (resource, path) => resource.addResource(path),
    api.root,
  )
  re.addMethod('GET', new LambdaIntegration(func))
  addCorsOptions(re, target)
  return func
}

const restApiDomainName = (target: environment.Environments) => {
  return `verifier.` + environment.valueOf(target).rootDomain
}
