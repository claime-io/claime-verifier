import { SynthUtils } from '@aws-cdk/assert'
import { App } from '@aws-cdk/core'
import { Environments } from '../lib/env'
import { Route53Stack } from '../lib/route53'
import { VerifierApiStack } from '../lib/verifierapi'

test('resource created', () => {
  const app = new App()
  const target = Environments.TEST
  const r53stack = new Route53Stack(app, 'test-hostedZone', target)
  const stack = new VerifierApiStack(app, 'test', target, {
    hostedZone: r53stack.hostedZone,
  })
  expect(SynthUtils.toCloudFormation(stack)).toMatchSnapshot()
})
