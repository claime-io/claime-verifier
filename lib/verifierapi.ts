import { LambdaIntegration, RestApi } from '@aws-cdk/aws-apigateway'
import { Function, Runtime, Tracing } from '@aws-cdk/aws-lambda'
import * as cdk from '@aws-cdk/core'
import {
  addCorsOptions,
  code,
  environmentVariables,
  withCustomDomain,
} from './api'
import { basicPolicytStatements } from './discord'
import * as environment from './env'

export class VerifierApiStack extends cdk.Stack {
  constructor(
    scope: cdk.Construct,
    id: string,
    target: environment.Environments,
    props?: cdk.StackProps,
  ) {
    super(scope, id, props)
    const api = new RestApi(this, 'RestApi', {
      restApiName: environment.withEnvPrefix(target, 'verifier'),
    })
    apifunction(this, this.region, this.account, 'verify', target, api)
    withCustomDomain(this, api, restApiDomainName(target), target)
  }
}

// GET /${eoa}?type=domain
// -> いったんbool,at,actual
// -> (ゆくゆくは検証NGだった場合に理由出したい)
const apifunction = (
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
  api.root.addMethod('GET', new LambdaIntegration(func))
  addCorsOptions(api.root, target)
  return func
}

const restApiDomainName = (target: environment.Environments) => {
  return `verifier.` + environment.valueOf(target).rootDomain
}
