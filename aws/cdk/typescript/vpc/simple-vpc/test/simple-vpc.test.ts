
import { App } from 'aws-cdk-lib';
import { SimpleVpcStack } from '../lib/simple-vpc-stack';
import { Template } from 'aws-cdk-lib/assertions';

const EIP = 'AWS::EC2::EIP';
const InternetGateway = 'AWS::EC2::InternetGateway';
const NatGateway = 'AWS::EC2::NatGateway';
const NetworkAcl = 'AWS::EC2::NetworkAcl';
const NetworkAclEntry = 'AWS::EC2::NetworkAclEntry';
const RouteTable = 'AWS::EC2::RouteTable';
const Subnet = 'AWS::EC2::Subnet';
const VPC = 'AWS::EC2::VPC';

test('SimpleVpc Created', () => {
    const app = new App();
    const stack = new SimpleVpcStack(app, 'SimpleVpcTest');
    const template = Template.fromStack(stack);

    // validate internet gateway
    template.resourceCountIs(InternetGateway, 1);

    // validate nat gateway
    template.resourceCountIs(EIP, 1);
    template.resourceCountIs(NatGateway, 1);

    // validate network acls
    template.resourceCountIs(NetworkAcl, 2);
    template.resourceCountIs(NetworkAclEntry, 2);

    // validate route tables
    template.resourceCountIs(RouteTable, 4);

    // validate subnets
    template.hasResourceProperties(Subnet, {
        CidrBlock: '10.0.0.0/26'
    });
    template.hasResourceProperties(Subnet, {
        CidrBlock: '10.0.0.64/26'
    });
    template.hasResourceProperties(Subnet, {
        CidrBlock: '10.0.0.128/26'
    });
    template.hasResourceProperties(Subnet, {
        CidrBlock: '10.0.0.192/26'
    });

    // validate vpc
    template.hasResourceProperties(VPC, {
        CidrBlock: '10.0.0.0/24'
    });
});
