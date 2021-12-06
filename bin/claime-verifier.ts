import * as cdk from '@aws-cdk/core'
import { CertificateStack } from '../lib/certificate'
import { DatasourceStack } from '../lib/datasource'
import { DiscordStack } from '../lib/discord'
import { DiscordApiStack } from '../lib/discordapi'
import * as environment from '../lib/env'
import { Route53Stack } from '../lib/route53'
import { VerifierApiStack } from '../lib/verifierapi'

const app = new cdk.App()
const target: environment.Environments = app.node.tryGetContext(
  'target',
) as environment.Environments
if (!target || !environment.valueOf(target))
  throw new Error('Invalid target environment')
const r53Stack = new Route53Stack(
  app,
  environment.withEnvPrefix(target, 'route53'),
  target,
  {},
)
new CertificateStack(
  app,
  environment.withEnvPrefix(target, 'certificate'),
  target,
  { hostedZone: r53Stack.hostedZone },
)
new VerifierApiStack(
  app,
  environment.withEnvPrefix(target, 'verifierapi'),
  target,
  { hostedZone: r53Stack.hostedZone },
)
new DiscordApiStack(app, environment.withEnvPrefix(target, 'restapi'), target, {
  hostedZone: r53Stack.hostedZone,
})
new DatasourceStack(
  app,
  environment.withEnvPrefix(target, 'datasource'),
  target,
  {},
)
new DiscordStack(app, environment.withEnvPrefix(target, 'discord'), target, {})
