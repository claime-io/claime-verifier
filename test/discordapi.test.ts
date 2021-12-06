import { SynthUtils } from '@aws-cdk/assert'
import { App } from '@aws-cdk/core'
import { DiscordApiStack } from '../lib/discordapi'
import { Environments } from '../lib/env'
import { Route53Stack } from '../lib/route53'

test('resource created', () => {
  const app = new App()
  const target = Environments.TEST
  const r53stack = new Route53Stack(app, 'test-hostedZone', target)
  const stack = new DiscordApiStack(app, 'test', target, {
    hostedZone: r53stack.hostedZone,
  })
  expect(SynthUtils.toCloudFormation(stack)).toMatchSnapshot()
})
