import { SynthUtils } from '@aws-cdk/assert'
import { App } from '@aws-cdk/core'
import { Environments } from '../lib/env'
import { VerifierApiStack } from '../lib/verifierapi'

test('resource created', () => {
  const app = new App()
  const target = Environments.TEST
  const stack = new VerifierApiStack(app, 'test', target)
  expect(SynthUtils.toCloudFormation(stack)).toMatchSnapshot()
})
