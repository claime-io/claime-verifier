import { SynthUtils } from '@aws-cdk/assert'
import { App } from '@aws-cdk/core'
import { Environments } from '../lib/env'
import { DiscordStack } from './../lib/discord'

test('resource created', () => {
  const app = new App()
  const target = Environments.TEST
  const stack = new DiscordStack(app, 'test', target)
  expect(SynthUtils.toCloudFormation(stack)).toMatchSnapshot()
})
