
# Services

## Resources

### Enable GCP Services

* Cloud Build
* Cloud Compute
* Cloud DNS
* Cloud Domains
* Cloud Run
* Cloud Source Repositories
* Container Registry

### IAM

* Role Mappings `build`
* Role Mappings `core`
* Service Account
* Service Account Key

## Outputs

service account json key
```
terraform output -raw service_account_json | base64 -d -
```
