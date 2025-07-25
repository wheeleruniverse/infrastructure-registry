# Infrastructure Registry

A comprehensive collection of Infrastructure as Code (IaC) templates and tools for AWS, Azure, and multi-cloud deployments. This repository serves as a reference library for common infrastructure patterns, cloud automation, and enterprise-grade solutions.

## Overview

This repository contains practical, production-ready infrastructure templates organized by cloud provider and tool type. It includes everything from basic service configurations to complex multi-tier applications and enterprise account management systems.

### Key Features

- **Multi-Cloud Support**: AWS, Azure, and HashiCorp Terraform configurations
- **Production-Ready Templates**: Enterprise-grade solutions with security best practices
- **Account Management**: Automated AWS account vending machine with baseline configurations
- **Cloud Governance**: Cost management, compliance, and monitoring templates
- **Challenge Solutions**: Real-world implementations from cloud certification challenges

## Repository Structure

```
├── aws/                          # AWS-specific resources
│   ├── cdk/                      # AWS CDK projects (TypeScript)
│   ├── cli/                      # AWS CLI automation scripts
│   ├── cloudformation/           # CloudFormation templates
│   │   ├── project/             # Complete project templates
│   │   └── service/             # Individual service templates
│   ├── lambda/                   # Lambda functions
│   │   └── go/                  # Go-based Lambda functions
│   └── sam/                      # SAM application templates
├── azure/                        # Azure-specific resources
│   └── arm/                      # ARM templates
├── hashicorp/                    # HashiCorp tool configurations
│   └── terraform/                # Terraform modules and projects
```

## Key Components

### AWS Account Vending Machine
Located in `aws/lambda/go/account-vending-machine/`, this is a sophisticated Go-based Lambda function that automates:
- AWS Organizations account creation
- Baseline security configuration deployment
- Cross-account role assumption
- VPC cleanup and standardization

### Baseline Infrastructure Templates
The `aws/cloudformation/project/wheelerswebservices/` directory contains enterprise baseline templates including:
- **Budget Management**: Multi-threshold budget alerts and notifications
- **Cost & Usage Reporting**: Automated CUR setup with S3 lifecycle policies  
- **CloudTrail**: Organization-wide audit logging
- **Security**: Encryption, access controls, and compliance configurations

### Service-Specific Templates
Modular CloudFormation templates for individual AWS services:
- **Compute**: VPC, EC2, and networking configurations
- **Storage**: S3 buckets with security best practices
- **Monitoring**: CloudTrail and budget management
- **Data**: DynamoDB tables with proper scaling
- **CDN**: CloudFront distributions for static websites

### Cloud Challenge Solutions
Real-world implementations from cloud certification challenges:
- **COVID ETL Pipeline** (`cgc-aws-covid-etl`): Automated data processing with Lambda, DynamoDB, and EventBridge
- **Multi-tier Applications**: VPC, application, and database tiers
- **Azure VM Deployments**: ARM templates for compute resources

## Getting Started

### Prerequisites
- AWS CLI configured with appropriate permissions
- Terraform >= 1.0 (for Terraform templates)
- Node.js and npm (for CDK projects)
- Go 1.19+ (for Lambda functions)

### Quick Start
1. Choose the appropriate directory for your cloud provider and tool
2. Review the README in each subdirectory for specific instructions
3. Customize parameters and variables for your environment
4. Deploy using the appropriate tool (CloudFormation, Terraform, etc.)

### Example: Deploy a Simple VPC
```bash
# Using CloudFormation
aws cloudformation create-stack \
  --stack-name my-vpc \
  --template-body file://aws/cloudformation/service/vpc/yml/single-subnet-vpc.yml

# Using CDK
cd aws/cdk/typescript/vpc/simple-vpc
npm install
cdk deploy
```

## Security & Best Practices

All templates in this repository follow cloud security best practices:
- **Encryption**: Data encrypted at rest and in transit
- **Access Control**: Principle of least privilege IAM policies
- **Network Security**: Private subnets, security groups, and NACLs
- **Monitoring**: CloudTrail, budget alerts, and logging enabled
- **Compliance**: Enterprise-grade tagging and governance

## Contributing

This repository serves as a reference collection. When adding new templates:
1. Follow existing directory structure conventions
2. Include comprehensive documentation
3. Implement security best practices
4. Add appropriate tags and metadata
5. Test thoroughly before committing

## License

This project contains infrastructure templates for educational and reference purposes.
