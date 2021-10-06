const projectName = 'claime-verifier'

export enum Environments {
  PROD = 'prod',
  DEV = 'dev',
  TEST = 'test',
}
export interface EnvironmentVariables {
  allowedOrigin: string
}

const EnvironmentVariablesSetting: {
  [key in Environments]: EnvironmentVariables
} = {
  [Environments.PROD]: {
    allowedOrigin: '',
  },
  [Environments.DEV]: {
    allowedOrigin: 'https://claime-webfront-k6p1srx99-squard.vercel.app',
  },
  [Environments.TEST]: {
    allowedOrigin: '',
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
