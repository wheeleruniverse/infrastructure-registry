import { Stack, StackProps, Tags } from 'aws-cdk-lib';
import {
  AclCidr,
  AclTraffic,
  Action,
  NetworkAcl,
  NetworkAclEntry,
  SubnetType,
  TrafficDirection,
  Vpc
} from 'aws-cdk-lib/aws-ec2';
import { Construct } from 'constructs';

export class SimpleVpcStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);
    const simpleVpc = new Vpc(this, 'SimpleVpc', {
      cidr: '10.0.0.0/24',
      natGateways: 1,
    });
    this.createNACL(SubnetType.PRIVATE_WITH_NAT, simpleVpc);
    this.createNACL(SubnetType.PUBLIC, simpleVpc);
  }

  private createNACL(subnetType: SubnetType, vpc: Vpc): void {
    const subnetSelection = vpc.selectSubnets({ subnetType });
    const networkAclName = SubnetType.PUBLIC === subnetType ? 'PublicNacl' : 'PrivateNacl';
    const networkAcl = new NetworkAcl(this, networkAclName, { networkAclName, subnetSelection, vpc });
    Tags.of(networkAcl).add('Name', networkAclName);

    const ephemeralPortsEntryName = `${networkAclName}EphemeralPorts`;
    new NetworkAclEntry(this, ephemeralPortsEntryName, {
      cidr: AclCidr.anyIpv4(),
      direction: TrafficDirection.EGRESS,
      ruleAction: Action.ALLOW,
      ruleNumber: 9999,
      networkAcl,
      networkAclEntryName: ephemeralPortsEntryName,
      traffic: AclTraffic.tcpPortRange(1024, 65535),
    });
  }
}
