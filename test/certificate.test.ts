import { SynthUtils } from '@aws-cdk/assert'
import { App } from '@aws-cdk/core'
import { DatasourceStack } from '../lib/datasource'
import { Environments } from '../lib/env'

test('resource created', () => {
  const app = new App()
  const target = Environments.TEST
  const stack = new DatasourceStack(app, 'test', target)
  expect(SynthUtils.toCloudFormation(stack)).toMatchSnapshot()
})
