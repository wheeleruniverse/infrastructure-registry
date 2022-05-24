#!/bin/bash

# https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ec2/delete-vpc.html
# you can't delete the main route table
# you can't delete the default network acl
# you can't delete the default security group
  
# get regions
for region in $(aws ec2 describe-regions \
--region us-east-1 \
| jq -r .Regions[].RegionName); do

  # get default vpc
  vpc=$(aws ec2 describe-vpcs \
  --filter Name=isDefault,Values=true \
  --region ${region} \
  | jq -r .Vpcs[0].VpcId)
  if [[ ${vpc} == 'null' ]]; then
    continue
  fi

  # get internet gateway
  igw=$(aws ec2 describe-internet-gateways \
  --filter Name=attachment.vpc-id,Values=${vpc} \
  --region ${region} \
  | jq -r .InternetGateways[0].InternetGatewayId)
  if [[ ${igw} != 'null' ]]; then
    echo "aws ec2 detach-internet-gateway --internet-gateway-id ${igw} --region ${region} --vpc-id ${vpc}"
    aws ec2 detach-internet-gateway --internet-gateway-id ${igw} --region ${region} --vpc-id ${vpc}
    
    echo "aws ec2 delete-internet-gateway --internet-gateway-id ${igw} --region ${region}"
    aws ec2 delete-internet-gateway --internet-gateway-id ${igw} --region ${region}
  fi

  # get subnets
  subnets=$(aws ec2 describe-subnets \
  --filters Name=vpc-id,Values=${vpc} \
  --region ${region} \
  | jq -r .Subnets[].SubnetId)
  if [[ ${subnets} != 'null' ]]; then
    for subnet in ${subnets}; do
      echo "aws ec2 delete-subnet --region ${region} --subnet-id ${subnet}"
      aws ec2 delete-subnet --region ${region} --subnet-id ${subnet}
    done
  fi

  echo "aws ec2 delete-vpc --region ${region} --vpc-id ${vpc}"
  aws ec2 delete-vpc --region ${region} --vpc-id ${vpc}

done