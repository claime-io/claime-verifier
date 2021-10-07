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
}

const EnvironmentVariablesSetting: {
  [key in Environments]: EnvironmentVariables
} = {
  [Environments.PROD]: {
    allowedOrigin: '',
    rootDomain: '',
    hostedZoneId: '',
    certificateArn: '',
  },
  [Environments.DEV]: {
    allowedOrigin: 'https://claime-webfront-k6p1srx99-squard.vercel.app',
    rootDomain: 'claime-dev.tk',
    hostedZoneId: 'Z08620181ARYV5PENJUEI',
    certificateArn:
      'arn:aws:acm:us-east-1:495476032358:certificate/92761d01-6e38-4ca4-8580-a945a5c379fe',
  },
  [Environments.TEST]: {
    allowedOrigin: '',
    rootDomain: 'test',
    hostedZoneId: 'test',
    certificateArn: '',
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
