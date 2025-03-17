# AWS SageMaker

## What is AWS SageMaker?

- We are going to train and deploy our models for inference on AWS SageMaker.
- SageMaker will be exposing the AI models via exposing the ML models via an API allowing for inference requests to be made.
- SageMaker will be handling the scaling of the inference requests.
- SageMakers consists of pre-trained models as well which are best suited for the task

- AWS SageMaker is a fully managed service that provides every developer and data scientist with the ability to build, train, and deploy machine learning (ML) models quickly. SageMaker reduces the time it takes to experiment with different algorithms and hyperparameters to achieve the best results.
- We are going to use AWS SageMaker to host and deploy our machine learning models. SageMaker is a fully managed service that provides every developer and data scientist with the ability to build, train, and deploy machine learning (ML) models quickly.
- SageMaker reduces the time it takes to experiment with different algorithms and hyperparameters to achieve the best results.
- SageMaker is integrated with Jupyter notebooks, so you can easily build, train, and deploy models without leaving your notebook.
- SageMaker also provides built-in algorithms for common machine learning tasks, such as linear regression, logistic regression, and neural networks.
- SageMaker also provides tools for data preparation, model evaluation, and model deployment.
- SageMaker also provides tools for model monitoring and model retraining.
- SageMaker also provides tools for model explainability and model interpretability.
- SageMaker also provides tools for model fairness and model bias.
- SageMaker also provides tools for model privacy and model security.

## How to use AWS SageMaker?

### Step 1: Create a SageMaker Notebook Instance

- Go to the SageMaker console and click on "Notebook instances" in the left-hand menu.
- Click on "Create notebook instance".
- Choose a name for your notebook instance.
- Choose an instance type. For this example, we'll use a ml.t2.medium instance.
- Choose a IAM role. For this example, we'll use the default execution role.
- Click on "Create notebook instance".
- Wait for the notebook instance to be created. This may take a few minutes.
- Once the notebook instance is created, click on "Open Jupyter" to open the Jupyter notebook interface.

# Billing

- AWS Sagemaker lets you work on it without charge for 12 months.
- Use can either use pre-trained models or fine-tune them for your needs.
- You can use the SageMaker Studio to build, train, and deploy models.
- You can use the SageMaker Ground Truth to label data.
- You can use the SageMaker Neo to optimize models.
- You can use the SageMaker Clarify to detect bias.
- You can use the SageMaker Debugger to debug models.
- You can use the SageMaker Model Monitor to monitor models.
- You can use the SageMaker Feature Store to store features.
- Storage, logging and networking could lead to higher costs if not managed properly.
- You can use the SageMaker Autopilot to automate the model building process.
- Running high load inference on AWS SageMaker is going to cost money, since you are going to be using the GPU instances.

## Considerations

- Configure auto-scaling limits to avoid unnecessary charges. Prevent the instances from automatically scaling up.
- Setup rate limiting to restrict the number of API requests per minute/hour.
- Use spot instances to reduce costs.
- Use SageMaker Studio to reduce costs.
- Use SageMaker Ground Truth to reduce costs.
- Use SageMaker Neo to reduce costs.
- Consider setting a budget cap in Billing.
- Use IAM policies to restrict access to SageMaker.
- Use AWS CloudWatch to monitor costs.
- Use AWS Trusted Advisor to monitor costs.
- Use AWS Cost Explorer to monitor costs.
- Use AWS Budgets to monitor costs.

## More

- Also, it is possible to use AWS SageMaker to train and deploy models on AWS Lambda.
- AWS SageMaker can be used to train and deploy models on AWS Fargate.
- AWS SageMaker can be used to train and deploy models on AWS EKS.
- AWS SageMaker can be used to train and deploy models on AWS ECS.
- AWS SageMaker can be used to train and deploy models on AWS Batch.
- AWS SageMaker can be used to train and deploy models on AWS Glue.
- AWS SageMaker can be used to train and deploy models on AWS Step Functions.
- AWS SageMaker can be used to train and deploy models on AWS Data Pipeline.
- AWS SageMaker can be used to train and deploy models on AWS Data Wrangler.
- AWS SageMaker can be used to train and deploy models on AWS Data Exchange.
- AWS SageMaker can be used to train and deploy models on AWS Data Lake.

## Considerations for different AI/ML options

- Usually fine-tuning or training a model incurs additional costs.
- Fine-tuning a model is usually faster than training a model from scratch.
- Fine-tuning a model is usually cheaper than training a model from scratch.
- Custom ML model training is more expensive and hence tuning a pre-trained model is almost always preference.

## Other AI/ML Options

- When deploying API'S, we need to monitor number of requests per month. For example, Cloud Run (owned by Google) has a limit of 2 million requests per month.
- GKE ( Google Kubernetes Service ), which is a more powerful service by Google and runs on GPUs, has 1 free Autopilot cluster per month (small-scale) and you must make sure to not exceed the limit.
- Vertex AI offers free requests for Gemini API, which it does not for other models.
- Gemini is a pre-trained model with Vertex AI and is a product of Google.
- Notebooks have a free tier available with Vertex AI.

## AI Models

- We are going to use TensorFlow, Scikit-learn, Supervised Learning, Unsupervised Learning, Reinforcement Learning, Natural Language Processing, Time Series Analysis possibly
- RL Agents can be trained on AWS SageMaker.
- RL Agents can be trained on AWS SageMaker using the RL Toolkit. Expose RL-powered decision-making through Cloud Run APIs

## Example of how we would call the Vertex AI API in Go. AWS SageMaker is similar.

```bash

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
)

const vertexAIEndpoint = "https://us-central1-aiplatform.googleapis.com/v1/projects/YOUR_PROJECT_ID/locations/us-central1/publishers/google/models/gemini-pro:predict"

type RequestPayload struct {
	Instances []map[string]interface{} `json:"instances"`
}

type ResponsePayload struct {
	Predictions []map[string]interface{} `json:"predictions"`
}

func getAccessToken() (string, error) {
	ctx := context.Background()
	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return "", err
	}

	token, err := creds.TokenSource.Token()
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

func predict(inputText string) (string, error) {
	// Create request payload
	payload := RequestPayload{
		Instances: []map[string]interface{}{
			{"prompt": inputText},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Get OAuth token
	token, err := getAccessToken()
	if err != nil {
		return "", err
	}

	// Send HTTP request
	req, err := http.NewRequest("POST", vertexAIEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return


```
