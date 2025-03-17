## Big Query

- BigQuery is used for Data Analytics. It provides 1 TB of queries per month for free. It is a serverless, highly scalable, and cost-effective multi-cloud data warehouse that enables super-fast SQL. It is deemed to be the best for Analytics.

- It is ideal for querying large datasets, performing complex analyses, and generating insights. It is also used for machine learning, data warehousing, and business intelligence.

- It offers 10 GB of free storage per month and 1 TB of queries per month. It is also used for data integration, data transformation, and data migration.

- Optimize BigQuery usage by batch queries and use of caching.

## Context in Machine Learning

- We are going to integrate BigQuery ( datawarehouse) that specializes in real time anomaly detection.
- BigQuery ML is a managed service that enables you to create, train, and serve machine learning models.
- We are going to train the models using BigQuery ML, TensorFlow, PyTorch, and XGBoost.
- BigQuery ML is a managed service that enables you to create, train, and serve machine learning models using BigQuery SQL.
- TensorFlow, PyTorch and XGBoost are open-source libraries for machine learning and deep learning.

## BigQuery ML

## Workflow

- A transaction request is sent to the Go API's hosted on AWS EKS
- The API logs the transaction request in BigQuery.
- BigQuery ML model is triggered to predict the anomaly. It does so by sending the data to an AI system
  i.e. AWS SageMaker, Vertex AI for fraud prediction
- if fraud is detected, the Cloud function sends an alert i.e. (email, webhook, SMS, logs) and generates a report and an action to lock the account possibly.
