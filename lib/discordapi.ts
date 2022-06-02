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

export class DiscordApiStack extends cdk.Stack {
  constructor(
    scope: cdk.Construct,
    id: string,
    target: environment.Environments,
    props?: cdk.StackProps,
  ) {
    super(scope, id, props)
    const api = new RestApi(this, 'RestApi', {
      restApiName: environment.withEnvPrefix(target, 'discord-restapi'),
    })
    const func = discordFunction(
      this,
      this.region,
      this.account,
      'verifydiscord',
      target,
      api,
    )
    withCustomDomain(this, api, restApiDomainName(target), target)
  }
}
// GET /${eoa}?type=domain
// -> いったんbool,at,actual
// -> (ゆくゆくは検証NGだった場合に理由出したい)
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
  addCorsOptions(rs, { methods: ['POST'] })
  return func
}

const restApiDomainName = (target: environment.Environments) => {
  return `discord.` + environment.valueOf(target).rootDomain
}
