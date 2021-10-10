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
    allowedOrigin: 'https://claime.io',
    rootDomain: 'claime.io',
    hostedZoneId: 'Z08305602GK0LP28IOTQ3',
    certificateArn:
      'arn:aws:acm:us-east-1:495476032358:certificate/4da06504-10a6-4231-8200-5581568a907c',
  },
  [Environments.DEV]: {
    allowedOrigin:
      'https://claime-webfront-git-feature-discord-squard.vercel.app',
    rootDomain: 'claime-dev.tk',
    hostedZoneId: 'Z08620181ARYV5PENJUEI',
    certificateArn:
      'arn:aws:acm:us-east-1:495476032358:certificate/7ba1d525-652e-4388-9e09-c06b86f7f29a',
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
