import {
  DnsValidatedCertificate,
  ICertificate,
  ValidationMethod,
} from '@aws-cdk/aws-certificatemanager'
import { Construct, Stack } from '@aws-cdk/core'
import * as environment from './env'
import { hostedZoneFromId } from './route53'

export class CertificateStack extends Stack {
  public readonly certificate: ICertificate

  constructor(scope: Construct, id: string, target: environment.Environments) {
    super(scope, id)
    const { hostedZoneId } = environment.valueOf(target)
    if (environment.isProd(target)) return
    if (!hostedZoneId) throw new Error('env.hostedZoneId is requied')
    this.certificate = certificate(this, target)
  }
}

const certificate = (
  scope: Construct,
  target: environment.Environments,
): ICertificate => {
  const { rootDomain } = environment.valueOf(target)

  return new DnsValidatedCertificate(scope, 'Certificate', {
    domainName: `${rootDomain}`,
    subjectAlternativeNames: [`*.${rootDomain}`],
    hostedZone: hostedZoneFromId(scope, target),
    validationMethod: ValidationMethod.DNS,
    region: 'us-east-1',
  })
}
