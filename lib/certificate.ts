import {
  DnsValidatedCertificate,
  ICertificate,
  ValidationMethod,
} from '@aws-cdk/aws-certificatemanager'
import { IHostedZone } from '@aws-cdk/aws-route53'
import { Construct, Stack, StackProps } from '@aws-cdk/core'
import * as environment from './env'

type CertificateStackProps = {
  hostedZone: IHostedZone
}
export class CertificateStack extends Stack {
  public readonly certificate: ICertificate

  constructor(
    scope: Construct,
    id: string,
    target: environment.Environments,
    props: CertificateStackProps & StackProps,
  ) {
    super(scope, id)
    const { hostedZone } = props
    this.certificate = certificate(this, hostedZone, target)
  }
}

const certificate = (
  scope: Construct,
  hostedZone: IHostedZone,
  target: environment.Environments,
): ICertificate => {
  const { rootDomain } = environment.valueOf(target)
  return new DnsValidatedCertificate(scope, 'Certificate', {
    domainName: `${rootDomain}`,
    subjectAlternativeNames: [`*.${rootDomain}`],
    hostedZone,
    validationMethod: ValidationMethod.DNS,
    region: 'us-east-1',
  })
}
