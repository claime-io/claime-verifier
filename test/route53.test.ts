import { SynthUtils } from '@aws-cdk/assert'
import { App } from '@aws-cdk/core'
import { Environments } from '../lib/env'
import { Route53Stack } from '../lib/route53'

test('resource created', () => {
  const app = new App()
  const target = Environments.TEST
  const stack = new Route53Stack(app, 'test', target)
  expect(SynthUtils.toCloudFormation(stack)).toMatchSnapshot()
})
