import * as cdk from '@aws-cdk/core'
import { DatasourceStack } from '../lib/datasource'
import { DiscordStack } from '../lib/discord'
import * as environment from '../lib/env'
import { RestApiStack } from '../lib/restapi'

const app = new cdk.App()
const target: environment.Environments = app.node.tryGetContext(
  'target',
) as environment.Environments
if (!target || !environment.valueOf(target))
  throw new Error('Invalid target environment')

new DatasourceStack(
  app,
  environment.withEnvPrefix(target, 'datasource'),
  target,
  {},
)
new DiscordStack(app, environment.withEnvPrefix(target, 'discord'), target, {})
new RestApiStack(app, environment.withEnvPrefix(target, 'restapi'), target, {})
