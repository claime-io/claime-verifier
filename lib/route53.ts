import {
  ARecord,
  ARecordProps,
  CnameRecord,
  HostedZone,
  IHostedZone,
  PublicHostedZone,
  RecordTarget,
} from '@aws-cdk/aws-route53'
import { Construct, Stack, StackProps } from '@aws-cdk/core'
import * as environment from './env'

export class Route53Stack extends Stack {
  public readonly hostedZone: IHostedZone
  constructor(
    scope: Construct,
    id: string,
    target: environment.Environments,
    props?: StackProps,
  ) {
    super(scope, id, props)
    const { rootDomain } = environment.valueOf(target)
    if (environment.isProd(target)) {
      return
    }
    this.hostedZone = hostedZone(this, target)
    // Vercel Domain Verification
    const vercelARecordProps: ARecordProps = {
      zone: this.hostedZone,
      recordName: `${rootDomain}`,
      target: RecordTarget.fromIpAddresses('76.76.21.21'),
    }
    new ARecord(this, 'VercelARecord', vercelARecordProps)

    new CnameRecord(this, 'VercelCNAMERecord', {
      zone: this.hostedZone,
      recordName: `www.${rootDomain}`,
      domainName: 'cname.vercel-dns.com',
    })
  }
}

const hostedZone = (
  scope: Construct,
  target: environment.Environments,
): IHostedZone => {
  return new PublicHostedZone(scope, 'HostedZone', {
    zoneName: hostedZoneName(target),
    comment: 'created by cdk',
  })
}

export const hostedZoneName = (target: environment.Environments) => {
  return environment.valueOf(target).rootDomain
}

export const hostedZoneFromId = (
  scope: Construct,
  target: environment.Environments,
) => {
  return HostedZone.fromHostedZoneAttributes(scope, `HostedZone`, {
    zoneName: hostedZoneName(target),
    hostedZoneId: environment.valueOf(target).hostedZoneId,
  })
}
