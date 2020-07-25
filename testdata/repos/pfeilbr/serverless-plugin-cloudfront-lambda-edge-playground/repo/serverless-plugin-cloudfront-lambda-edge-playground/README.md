# serverless-plugin-cloudfront-lambda-edge-playground

learn [silvermine/serverless-plugin-cloudfront-lambda-edge](https://github.com/silvermine/serverless-plugin-cloudfront-lambda-edge)

visit <https://d3cztrjc4xcpde.cloudfront.net> and login with user01/password01

**Basic Auth Challenge**

```
username: user01
password: password01
```

users are stored in aws secrets manager and sourced from [`users.json`](users.json)

## Usage

```sh
# deploy all (provision infra, build static site, copy to s3, invalidate cache).
npm run deploy

# separate component deploys
npm run deploy-infrastructure
npm run build-static-site # gatsby
npm run publish-static-assets-to-bucket
npm run cloudfront:invalidate

# bucket
# not accessible because of Origin Access Identity applied
# e.g. http://s3-cf-private-static-site-01-dev.s3-website-us-east-1.amazonaws.com
open "http://$(node scripts/get-stack-property.js WebsiteBucketName).s3-website-us-east-1.amazonaws.com/index.html"

# cloudfront url. e.g. https://d3cztrjc4xcpde.cloudfront.net
open "https://$(node scripts/get-stack-property.js CloudFrontDistributionDomainName)"
```

## Removing Auth

Update CloudFront Behavior to remove "Viewer Request" "Lambda Function Association", then invalidate cache on all "*".

## Viewing CloudWatch Logs for Lambda@Edge Functions

logging takes place in the region of the edge (PoP) location.  This will vary based on the client location.

Visit CloudFront | Monitoring | Lambda@Edge Functions | YOUR FUNCTION, then click `[View Function Metrics]` button

see [Determining the Lambda@Edge Region](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/lambda-edge-testing-debugging.html#lambda-edge-testing-debugging-determine-region) for more details

![](https://www.evernote.com/l/AAH3ZzTj929MNqiRZIcn1zj1G0WzeI5ZSGQB/image.png)

## Lambda@Edge Request

```json
{
    "event": {
        "Records": [
            {
                "cf": {
                    "config": {
                        "distributionDomainName": "d13ydba49ilc9v.cloudfront.net",
                        "distributionId": "E20V9SS2N0VT6P",
                        "eventType": "viewer-request",
                        "requestId": "NcxNfOh2NaptKOLBlSimVF7AAlRp10OM-F1-CMZsJOGSMEkHUVd25A=="
                    },
                    "request": {
                        "clientIp": "100.11.96.70",
                        "headers": {
                            "host": [
                                {
                                    "key": "Host",
                                    "value": "d13ydba49ilc9v.cloudfront.net"
                                }
                            ],
                            "user-agent": [
                                {
                                    "key": "User-Agent",
                                    "value": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36"
                                }
                            ],
                            "authorization": [
                                {
                                    "key": "authorization",
                                    "value": "Basic dXNlcjpwYXNz"
                                }
                            ],
                            "upgrade-insecure-requests": [
                                {
                                    "key": "upgrade-insecure-requests",
                                    "value": "1"
                                }
                            ],
                            "sec-fetch-mode": [
                                {
                                    "key": "sec-fetch-mode",
                                    "value": "navigate"
                                }
                            ],
                            "sec-fetch-user": [
                                {
                                    "key": "sec-fetch-user",
                                    "value": "?1"
                                }
                            ],
                            "accept": [
                                {
                                    "key": "accept",
                                    "value": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3"
                                }
                            ],
                            "sec-fetch-site": [
                                {
                                    "key": "sec-fetch-site",
                                    "value": "none"
                                }
                            ],
                            "accept-encoding": [
                                {
                                    "key": "accept-encoding",
                                    "value": "gzip, deflate, br"
                                }
                            ],
                            "accept-language": [
                                {
                                    "key": "accept-language",
                                    "value": "en-US,en;q=0.9,nb;q=0.8,fr;q=0.7"
                                }
                            ],
                            "if-none-match": [
                                {
                                    "key": "if-none-match",
                                    "value": "\"2cad58ac06c32be3c6384050881bc507\""
                                }
                            ],
                            "if-modified-since": [
                                {
                                    "key": "if-modified-since",
                                    "value": "Tue, 10 Sep 2019 21:32:57 GMT"
                                }
                            ]
                        },
                        "method": "GET",
                        "querystring": "",
                        "uri": "/"
                    }
                }
            }
        ]
    }
}
```

## Resources

* add origin access identity
    * see the following for cfn markup <https://github.com/lroguet/amzn-cloudformation/blob/master/storage-content-delivery/static-website-with-cloudfront.yml>
    * [Restricting Access to Amazon S3 Content by Using an Origin Access Identity](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/private-content-restricting-access-to-s3.html)
* S3 Bucket | [Granting Permission to an Amazon CloudFront Origin Identity](https://docs.aws.amazon.com/AmazonS3/latest/dev/example-bucket-policies.html#example-bucket-policies-use-case-6)

## TODO

* DNS (route 53)
    * [Routing Traffic to an Amazon CloudFront Web Distribution by Using Your Domain Name](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/routing-to-cloudfront-distribution.html).  Covers [root|sub]domainsG
    * see https://www.brautaset.org/articles/2017/route-53-cloudformation.html.  Contains ApexRecordSet (example.com) and WwwRecordSet (www.example.com)
    * [Alias Resource Record Set for a CloudFront Distribution](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/quickref-route53.html#w2ab1c17c23c81c11)
* update `serverless.yml:iamRoleStatements` with dynamic region, account id, and secretsmanager name
    * via [serverless-pseudo-parameters](https://github.com/svdgraaf/serverless-pseudo-parameters#readme) plugin

## Scratch

```
$(node ./scripts/get-stack-property.js )
```