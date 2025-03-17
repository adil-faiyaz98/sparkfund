# Data Pipelines

## AWS Glue

- We are going to use AWS Glue for the data pipelines

## What is AWS Glue ?

- It is a fully managed ETL (Extract, Transform, Load) service. It allows you to prepare and load data for analytics. It is a serverless service. You don't have to provision any servers. You just need to define the ETL jobs and Glue will take care of the rest.

## What is ETL ?

- ETL stands for Extract, Transform, Load. It is a process in database usage and especially in data warehousing that involves: 1. Extracting data from homogeneous or heterogeneous data sources. 2. Transforming it into an appropriate form for analysis. 3. Loading it into a data warehouse, database, or another type of information system.

## Why use AWS Glue ?

- It is a fully managed service, so you don't have to worry about provisioning or managing servers. 2. It is serverless, so you only pay for what you use. 3. It is integrated with other AWS services like S3, Redshift, Athena, etc. 4. It is easy to use, you just need to define the ETL jobs and Glue will take care of the rest. 5. It is scalable, you can scale up or down depending on your needs.

## What are we achieving from it towards Machine Learning API integration with our Microservices ?

- The data pipelines will be used to prepare the data for the machine learning models. The data will be extracted from the data sources, transformed into the appropriate form for the models, and then loaded into the data lake. The machine learning models will then be trained on the data in the data lake. The trained models will be deployed as APIs and integrated with our microservices.

## What are the steps involved in the data pipelines ?

- The steps involved in the data pipelines are as follows: 1. Extract data from the data sources. 2. Transform the data into the appropriate form for the machine learning models. 3. Load the data into the data lake. 4. Train the machine learning models on the data in the data lake. 5. Deploy the trained models as APIs. 6. Integrate the APIs with our microservices.

## What are the data sources ?

- The data sources are the databases and other data stores that contain the data that we need to train the machine learning models. For example, we may have data about user transactions, user demographics, and user behavior in our databases. This data will be used to train the machine learning models.
