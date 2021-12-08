#!/usr/bin/env node
import 'source-map-support/register';
import { App } from 'aws-cdk-lib';
import { SimpleVpcStack } from '../lib/simple-vpc-stack';

const app = new App();
new SimpleVpcStack(app, 'SimpleVpcStack', {});
