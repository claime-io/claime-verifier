import {
  AttributeType,
  BillingMode,
  Table,
  TableProps,
} from '@aws-cdk/aws-dynamodb'
import { Effect, PolicyStatement } from '@aws-cdk/aws-iam'
import * as cdk from '@aws-cdk/core'
import * as environment from './env'
export class DatasourceStack extends cdk.Stack {
  constructor(
    scope: cdk.Construct,
    id: string,
    target: environment.Environments,
    props?: cdk.StackProps,
  ) {
    super(scope, id, props)
    const table = new Table(
      this,
      'claime-verifier-table',
      ddbProps(tableName(target)),
    )
    table.addGlobalSecondaryIndex({
      indexName: 'GSI-1',
      partitionKey: {
        name: 'PK',
        type: AttributeType.STRING,
      },
      sortKey: {
        name: 'Timestamp',
        type: AttributeType.STRING,
      },
    })
  }
}

const tableName = (target: environment.Environments) => {
  return environment.withEnvPrefix(target, 'main')
}

const ddbProps = (tableName: string): TableProps => {
  return {
    tableName: tableName,
    partitionKey: {
      name: `PK`,
      type: AttributeType.STRING,
    },
    sortKey: {
      name: `SK`,
      type: AttributeType.STRING,
    },
    billingMode: BillingMode.PAY_PER_REQUEST,
    removalPolicy: cdk.RemovalPolicy.DESTROY,
    pointInTimeRecovery: true,
  }
}

export const dataSourceReadWritePolicyStatement = (
  region: string,
  account: string,
  target: environment.Environments,
) => {
  return new PolicyStatement({
    effect: Effect.ALLOW,
    actions: [
      'dynamodb:Put*',
      'dynamodb:Get*',
      'dynamodb:Scan*',
      'dynamodb:Query*',
      'dynamodb:Delete*',
      'dynamodb:Batch*',
    ],
    resources: [
      `arn:aws:dynamodb:${region}:${account}:table/${tableName(target)}*`,
    ],
  })
}
