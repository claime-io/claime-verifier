const projectName = 'claime-verifier'

export enum Environments {
  PROD = 'prod',
  DEV = 'dev',
  TEST = 'test',
}
export interface EnvironmentVariables {
  allowedOrigin: string
  rootDomain: string
  hostedZoneId: string
  certificateArn: string
  subgraphEndpoint: string
  ownerEOA: string
}

const EnvironmentVariablesSetting: {
  [key in Environments]: EnvironmentVariables
} = {
  [Environments.PROD]: {
    allowedOrigin: 'https://claime.io',
    rootDomain: 'claime.io',
    hostedZoneId: 'Z08305602GK0LP28IOTQ3',
    certificateArn:
      'arn:aws:acm:us-east-1:495476032358:certificate/4da06504-10a6-4231-8200-5581568a907c',
    subgraphEndpoint: 'TBD',
    ownerEOA: '0x81A2863ED122811A1197dB2D9b90a720d73ac81c',
  },
  [Environments.DEV]: {
    allowedOrigin:
      'https://claime-webfront-git-feature-discord-squard.vercel.app',
    rootDomain: 'claime-dev.tk',
    hostedZoneId: 'Z08620181ARYV5PENJUEI',
    certificateArn:
      'arn:aws:acm:us-east-1:495476032358:certificate/7ba1d525-652e-4388-9e09-c06b86f7f29a',
    subgraphEndpoint:
      'https://api.studio.thegraph.com/query/8417/claime-rinkeby/v0.0.1',
    ownerEOA: '0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb',
  },
  [Environments.TEST]: {
    allowedOrigin: '',
    rootDomain: 'test',
    hostedZoneId: 'test',
    certificateArn: '',
    subgraphEndpoint: 'endpoint-of-subgraph',
    ownerEOA: '0x0000000000000000000000000000000000000000',
  },
}

export function valueOf(env: Environments): EnvironmentVariables {
  return EnvironmentVariablesSetting[env]
}
export const withEnvPrefix = (target: Environments, str: string) =>
  `${projectName}-${str}-${target}`

export const isProd = (target: Environments) => {
  return target === Environments.PROD
}
export const isDev = (target: Environments) => {
  return target === Environments.DEV
}
