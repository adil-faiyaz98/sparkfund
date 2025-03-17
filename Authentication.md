# Authentication

- When we integrate with AWS SageMaker, we have 2 options.
- [1] Use the Cloud AI client library '
- If chosing the Cloud AI client library, we need to provide the necessary credentials and configuration to authenticate with AWS SageMaker. This can be done using AWS Identity and Access Management (IAM) roles or access keys. The client library typically handles the authentication process, including obtaining and refreshing tokens as needed.

- [2] Make direct REST API calls.
- For direct REST API calls, we need to sign the requests using AWS Signature Version 4. This involves calculating a signature using the access key, secret key, and other parameters such as the request method, headers, and body. The signature is then included in the request as a header, allowing AWS to authenticate the request.
