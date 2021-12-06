import { App } from '@aws-cdk/core'
import { Environments } from '../lib/env'

test('resource created', () => {
  const app = new App()
  const target = Environments.TEST
  // TODO: fix too long test
  // const stack = new DiscordStack(app, 'test', target)
  // expect(SynthUtils.toCloudFormation(stack)).toMatchSnapshot()
})
