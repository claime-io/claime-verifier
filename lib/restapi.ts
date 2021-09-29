import {
  ApiKeySourceType,
  LambdaIntegration,
  RestApi,
} from "@aws-cdk/aws-apigateway";
import { Effect, PolicyStatement } from "@aws-cdk/aws-iam";
import { Code, Function, Runtime, Tracing } from "@aws-cdk/aws-lambda";
import * as cdk from "@aws-cdk/core";
import { resolve } from "path";
import * as environment from "../lib/env";

export class RestApiStack extends cdk.Stack {
  constructor(
    scope: cdk.Construct,
    id: string,
    target: environment.Environments,
    props?: cdk.StackProps
  ) {
    super(scope, id, props);
    const api = new RestApi(this, "RestApi", {
      restApiName: environment.withEnvPrefix(target, "restapi"),
      apiKeySourceType: ApiKeySourceType.HEADER,
    });
    functions(
      this,
      this.region,
      this.account,
      ["startsubscription"],
      target,
      api
    );
  }
}

const functions = (
  scope: cdk.Construct,
  region: string,
  account: string,
  resources: string[],
  target: environment.Environments,
  api: RestApi
) => {
  return resources.map((r) => {
    const resource = api.root.addResource(r);
    const func = new Function(scope, r, {
      functionName: `${environment.withEnvPrefix(target, r)}`,
      code: code(r),
      handler: "bin/main",
      timeout: cdk.Duration.minutes(1),
      runtime: Runtime.GO_1_X,
      tracing: Tracing.ACTIVE,
    });
    basicPolicytStatements(region, account).forEach((s) =>
      func.addToRolePolicy(s)
    );
    resource.addMethod("POST", new LambdaIntegration(func));
    return resource;
  });
};

export const basicPolicytStatements = (region: string, account: string) => {
  return [
    new PolicyStatement({
      effect: Effect.ALLOW,
      actions: [
        "dynamodb:Put*",
        "dynamodb:Get*",
        "dynamodb:Scan*",
        "dynamodb:Delete*",
        "dynamodb:Batch*",
      ],
      resources: [
        `arn:aws:dynamodb:${region}:${account}:table/claime-verifier-main*`,
      ],
    }),
    new PolicyStatement({
      effect: Effect.ALLOW,
      actions: ["ssm:Get*"],
      resources: [
        `arn:aws:ssm:${region}:${account}:parameter/claime-verifier*`,
      ],
    }),
  ];
};

const code = (dirname: string) => {
  return Code.fromAsset(
    resolve(`${__dirname}/../`, "lib", "functions", dirname, "bin", "main.zip")
  );
};
