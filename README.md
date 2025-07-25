# Infrastructure Registry

A curated collection of reusable AWS Infrastructure as Code (IaC) templates and automation tools. This repository serves as a reference library for common AWS infrastructure components, enterprise baselines, and architectural patterns.

## Overview

This repository focuses on **reusable AWS infrastructure templates and patterns** rather than complete project implementations. It contains modular, production-ready AWS components including CloudFormation templates, CDK constructs, and Go-based automation functions.

For complete project implementations with application code and deployment automation, see the [Project-Specific Infrastructure](#project-specific-infrastructure) section below.

### Key Features

- **AWS-Focused**: Curated collection of AWS infrastructure templates and automation tools
- **Enterprise Baselines**: Production-ready account setup and governance patterns
- **Security Best Practices**: Templates following AWS security frameworks and compliance standards
- **Service Patterns**: Modular templates for common AWS services and architectural patterns
- **Automation Tools**: Go-based Lambda functions for advanced AWS automation scenarios

## Repository Structure

```
├── aws/
│   ├── cdk/typescript/vpc/       # AWS CDK VPC template (TypeScript)
│   ├── cli/                      # AWS CLI automation scripts
│   ├── cloudformation/
│   │   ├── project/wheelerswebservices/  # Enterprise baseline templates
│   │   └── service/              # Individual service templates
│   │       ├── budgets/          # Budget management
│   │       ├── cloudtrail/       # Audit logging
│   │       ├── dynamodb/         # NoSQL database
│   │       ├── s3/               # Object storage
│   │       └── vpc/              # Virtual private cloud
│   └── lambda/go/                # Go-based automation functions
│       ├── account-vending-machine/   # AWS Organizations automation
│       └── cloudformation-cleanup/    # Resource cleanup automation
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


## Project-Specific Infrastructure

For complete, project-specific infrastructure implementations, see these dedicated repositories:

### Multi-Cloud Resume Infrastructure
- **Repository**: [multicloud-resume](https://github.com/wheeleruniverse/multicloud-resume)
- **Infrastructure**: [`iac/terraform/`](https://github.com/wheeleruniverse/multicloud-resume/tree/main/iac/terraform)
- **Description**: Complete multi-cloud resume hosting solution with Terraform configurations for AWS, Azure, and GCP

### Cloud Guru Challenge Projects
- **[cgc-aws-covid-etl](https://github.com/wheeleruniverse/cgc-aws-covid-etl/tree/main/cloudformation)**: Automated COVID-19 data ETL pipeline with Lambda, DynamoDB, and EventBridge scheduling
- **[cgc-aws-app-performance](https://github.com/wheeleruniverse/cgc-aws-app-performance/tree/main/terraform)**: AWS application performance monitoring and optimization solutions with ElastiCache
- **[cgc-aws-ml-recommendation-engine](https://github.com/wheeleruniverse/cgc-aws-ml-recommendation-engine)**: Machine learning-powered recommendation system on AWS with SageMaker and Athena
- **[cgc-azure-cicd](https://github.com/wheeleruniverse/cgc-azure-cicd/tree/main/arm)** : Complete Azure CI/CD pipeline infrastructure and automation with ARM templates
- **[cgc-gcp-resume-env](https://github.com/wheeleruniverse/cgc-gcp-resume-env/tree/main/core)** : Google Cloud resume hosting environment with Cloud Run and CDN
- **[cgc-multicloud-madness](https://github.com/wheeleruniverse/cgc-multicloud-madness/tree/main/sam)** : Multi-cloud chaos engineering and resilience testing framework with SAM

> **Note**: The infrastructure-registry focuses on reusable templates and patterns, while project-specific repositories contain complete implementations with application code, deployment scripts, and detailed project documentation.

## Getting Started

### Prerequisites
- AWS CLI configured with appropriate permissions
- Node.js and npm (for CDK templates)
- Go 1.19+ (for Lambda function development)

### Quick Start
1. Browse the `aws/` directory for the appropriate template type
2. Review the README in each subdirectory for specific instructions
3. Customize parameters and variables for your environment
4. Deploy using CloudFormation or CDK

### Examples

**Deploy a VPC with CloudFormation:**
```bash
aws cloudformation create-stack \
  --stack-name my-vpc \
  --template-body file://aws/cloudformation/service/vpc/yml/single-subnet-vpc.yml
```

**Deploy a VPC with CDK:**
```bash
cd aws/cdk/typescript/vpc/simple-vpc
npm install
cdk deploy
```

**Enterprise Account Baseline:**
```bash
aws cloudformation create-stack \
  --stack-name baseline \
  --template-body file://aws/cloudformation/project/wheelerswebservices/baseline.yml \
  --parameters ParameterKey=pAccount,ParameterValue=mycompany
```

## Security & Best Practices

All templates in this repository follow cloud security best practices:
- **Encryption**: Data encrypted at rest and in transit
- **Access Control**: Principle of least privilege IAM policies
- **Network Security**: Private subnets, security groups, and NACLs
- **Monitoring**: CloudTrail, budget alerts, and logging enabled
- **Compliance**: Enterprise-grade tagging and governance

## Contributing

This repository serves as a reference collection of **reusable infrastructure templates**. When adding new templates:
1. Ensure templates are **modular and reusable** across different projects
2. Follow existing directory structure conventions
3. Include comprehensive documentation with usage examples
4. Implement security best practices and compliance standards
5. Add appropriate tags and metadata for discoverability
6. Test thoroughly in multiple environments before committing

**Note**: Complete project implementations with application code should be maintained in separate project-specific repositories.

## License

This project contains infrastructure templates for educational and reference purposes.
