
# Core Services

## Resources

### Backend

* Cloud Run
* Cloud Run IAM Policy `allUsers`
* Container Registry

### CI/CD

* Cloud Build Trigger for app
* Cloud Build Trigger for web

### Database

* App Engine Application `Firestore`

### DNS

* Cloud DNS Record Set for api
* Cloud DNS Record Set for root
* Cloud DNS Record Set for www
* Cloud Run Domain Mapping `api.${var.fqdn}`
* HTTPS Load Balancer Domain Mapping [`var.fqdn`, `www.${var.fqdn}`]

### Frontend

* Cloud Storage Bucket
* Cloud Storage Bucket IAM Policy `allUsers`
* HTTPS Load Balancer
* HTTPS Load Balancer Cloud Storage Backend

## Variables

### image

**Description**

container registry container image

**Example**
```
image = "gcr.io/cloudguruchallenge-2108/wheelersadvice:09d42903517ad853459c51924b791081af748f13"
```

### service_account_json

**Description**

service account json key

**Example**
```
service_account_json = <<-EOT
{
  "type": "service_account",
  "project_id": "cloudguruchallenge-2108",
  "private_key_id": **redacted**,
  "private_key": **redacted**,
  "client_email": "wheelersadvice-svc@cloudguruchallenge-2108.iam.gserviceaccount.com",
  "client_id": **redacted**,
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/wheelersadvice-svc%40cloudguruchallenge-2108.iam.gserviceaccount.com"
}
EOT
```
> :warning: service_account_json contains sensitive information that must be protected!
> * never commit the value to version control
> * never share the value

