# Machine Learning Process

## 1. Problem Definition

- Clearly define the problem you want to solve.
- Identify the type of problem (classification, regression, clustering, etc.).
- Determine the success metrics (accuracy, precision, recall, F1 score, etc.).

### Tech Stack used:

- JIRA, Confluence for docuemtation
- Requirements gathering and establishing a backlog

## 2. Data Collection

- Gather data relevant to the problem. This can come from various sources like databases, APIs, web scraping, etc.
- Ensure the data is representative of the problem space.

### Tech Stack used:

- DataWarehouses
- Pre-trained proprietary trained data-set that is high in quality for the domain the model needs to be trained on.

## 3. Data Preprocessing

- Clean the data by handling missing values, outliers, and duplicates.
- Transform the data into a suitable format for analysis (normalization, encoding, etc.).
- Split the data into training, validation, and test sets.
- Feature engineering: Create new features or modify existing ones to improve model performance.

### Tech Stack used:

- Python Libraries like Pandas, Numpy, Scikit-learn, etc.
- Data Cleaning and Normalization Libraries like OpenCV, etc.
- Data Augmentation Libraries like Albumentations, etc.

## 4. Model Selection

- Choose appropriate algorithms for the problem (linear regression, decision trees, neural networks, etc.).
- Consider the trade-offs between model complexity and performance.
- Use techniques like cross-validation to evaluate model.

### Tech Stack used:

- Python Libraries like Scikit-learn, TensorFlow, Keras, etc.
- Model Selection Libraries like Optuna, Hyperopt, etc.
- Model Evaluation Libraries like Scikit-learn, TensorFlow, Keras, etc.
- Model Interpretation Libraries like SHAP, LIME, etc.
- Model Explainability Libraries like ELI5, etc.
- Model Visualization Libraries like Matplotlib, Seaborn etc.

## 5. Model Training

- Train the selected models on the training data. Use techniques like gradient descent, backpropagation, etc.
- Use techniques like early stopping, learning rate scheduling, etc. to prevent overfitting.
- Use techniques like batch normalization, dropout, etc. to improve model performance.
- Use techniques like regularization, weight decay, etc. to prevent overfitting.
- Use techniques like data augmentation, synthetic data generation, etc. to improve model performance.
- Use tecniques like Model Quantization, distillation and Optimization

### Tech Stack used:

- Python algorithms like Gradient Boosting, Random Forest, etc.
- Python Libraries like Scikit-learn, TensorFlow, Keras, etc.
- Model Training Libraries like PyTorch, etc.
- Model Optimization Libraries like Optuna, Hyperopt, etc.
- Model Quantization Libraries like TensorFlow Lite, etc.
- Model Distillation Libraries like DistilBERT, etc.
- Model Optimization Libraries like TensorFlow Model Optimization, etc.

## 6. Model Evaluation

- Evaluate the trained models on the validation data. Use techniques like confusion matrix, ROC curve, etc. to evaluate model performance.
- Use techniques like precision-recall curve, F1 score, etc. to evaluate model performance.
- Use techniques like AUC-ROC, etc. to evaluate model performance.
- Use techniques like cross-validation, etc. to evaluate model performance.

### Tech stack used:

- Tools like Jupyter Notebooks for experimentation and evaluation.
- Libraries like Scikit-learn, TensorFlow, Keras, etc. for model evaluation.
- Libraries like Matplotlib, Seaborn, etc. for visualization.
- Libraries like SHAP, LIME, etc. for model interpretation.
- Libraries like ELI5, etc. for model explainability.
- Libraries like Pandas, Numpy, etc. for data manipulation.
- Libraries like Optuna, Hyperopt, etc. for hyperparameter tuning.  


## 7. Model Deployment

- Deploy the best performing model to a production environment. Use techniques like containerization, orchestration, etc. to deploy model.
- Monitor the model's performance in production. Use techniques like A/B testing, etc. to monitor model performance.
- Update the model as needed based on new data or changing business requirements. Use techniques like retraining, fine-tuning, etc. to update model.
- Use techniques like model versioning, etc. to manage model.
- Use techniques like model explainability, etc. to understand model.
- Use techniques like model interpretability, etc. to interpret model.

### Tech Stack used:

- Python Libraries like Flask, Django, etc. for creating APIs.
- Docker, Kubernetes, etc. for containerization and orchestration.
- Cloud Hosting platforms like AWS, Azure or GCP that host Models on the cloud for Inference
- Model Serving platforms like TensorFlow Serving, TorchServe, etc. for serving models.
- Model Management platforms like MLflow, etc. for managing models.
- Model Monitoring platforms like Prometheus, Grafana, etc. for monitoring models.
- Model Explainability platforms like SHAP, LIME, etc. for explaining models.

## 8. Model Monitoring

- Monitor the model's performance in production. Use techniques like A/B testing, etc. to monitor model performance.
- Detect and address model drift. Use techniques like data drift detection, etc. to detect model drift.
- Use Monitoring and Observability like Prometheus / Grafana to compare the performance of the model in production vs how it was during training.
- If model drift is detected, retrain the model with new data. Use techniques like retraining, fine-tuning, etc. to update model.
- Use techniques like model versioning, etc. to manage model.
- Use techniques like model explainability, etc. to understand model and interpret model.
- Use techniques like model interpretability, etc. to interpret model and understand model.

### Tech Stack used:

- Prometheus, Grafana, etc. for monitoring.
- ELK Stack for logging.
- CloudWatch for monitoring.
- Datadog for monitoring.
- New Relic for monitoring.
- Splunk for monitoring.
- Use techniques like data drift detection, etc. to detect model drift.
- Use techniques like model versioning, etc. to manage model.
- Use techniques like model explainability, etc. to understand
- Use techniques like model interpretability, etc. to interpret model.
- Measure the model's performance in production using metrics like accuracy, precision, recall, F1 score, etc.

## 9. Model Retraining

- As new data evolves, the models needs to be retrained. Use techniques like retraining, fine-tuning, etc. to update model.
- Use techniques like model versioning, etc. to manage model.
- Use techniques like model explainability, etc. to understand model and interpret model.

### Tech Stack used:

- It's the same tools used before to train more
- Communication tools like Slack, Microsoft Teams, etc. for communicating with the team.

## 10. Model Interpretability

- Use techniques like SHAP, LIME, etc. to interpret model.
- Use techniques like model explainability, etc. to understand

## 11. Model Explainability

- For concerns related to transparency, the model is required to also mentioned as to how it arrived to the conclusion. Use techniques like SHAP, LIME, etc. to explain model

## 12. Model Optimization

- Use techniques like model quantization, distillation, pruning, etc. to optimize model.
- Use techniques like model compression, etc. to optimize model.

### Tech Stack used:

- Tools that perform Model Optimization like TensorFlow, PyTorch, etc.
- Tools that perform Model Compression like TensorFlow Lite, PyTorch Mobile, etc.
- Tools that perform Model Quantization like TensorFlow Lite, PyTorch Mobile, etc.
- Tools that perform Model Distillation like TensorFlow Lite, PyTorch Mobile, etc.
- Tools that perform Pruning would be the same as in some Python libaries.

## 13. Model Retirement

- After reaching the end of its lifecycle, the model is retired. Use techniques like model versioning, etc. to manage model.
- Use techniques like model explainability, etc. to understand model and interpret model.
- The data used for training the model is also archived.
- The pipelines contain logic for Canary deployments and Blue Green testing so that a new model which is more efficient
  would cost less and be more accurate can be deployed in production.
- The pipelines contain logic for A/B testing so that a new model which is more efficient
  would cost less and be more accurate can be deployed in production.

# Machine Learning Engineers

## Tech Stack they use

- Machine Learning Engineers use Architecture and Models for their work.
- They use Python, R, Julia, MATLAB, Scala, Java, C++, C#, Go, Rust, Swift, Kotlin, TypeScript, JavaScript, SQL, NoSQL, Hadoop, Spark, Kafka, Flink, Storm, Hive, Pig, Impala, Presto, Drill, Cassandra, MongoDB, Redis, Memcached, Elasticsearch.

### Architectures they use

- They use Convolutional Neural Networks (CNNs), Recurrent Neural Networks (RNNs), Long Short-Term Memory (LSTM), Gated Recurrent Units (GRUs), Transformers, Autoencoders, Generative Adversarial Networks (GANs), Reinforcement Learning, Transfer Learning, Ensemble Learning, Bagging, Boosting, Stacking, Voting, etc.

### Percentage of time they spend on each task

- 20% of their time is spent on Data Collection and Preprocessing.
- 30% of their time is spent on Model Selection and Training.
- 20% of their time is spent on Model Evaluation and Validation.
- 10% of their time is spent on Model Deployment.
- 10% of their time is spent on Model Monitoring and Maintenance.
- 10% of their time is spent on Model Documentation and Communication.

## Architecture and Model selection

- It is common for Machine Learning Engineers to choose an existing architecture or model and tweak it to fit the problem at hand.
- They do not typically design new architectures or models from scratch.
- Their responsibiltiies include fine tuning a model on the processed data-set that applies to the problem
- They also have to ensure that the model is optimized for performance and efficiency.
- They are also responsible for deploying the model to production.
- They are also responsible for monitoring the model's performance in production and retraining it as needed.
- They are also responsible for documenting the model's performance and making it available to other data scientists and engineers.
- They are also responsible for communicating the model's performance to stakeholders.
- They must make sure that the model is free of Bias and Hallucinations. Alkso, the model must be explainable.
- They should also ensure that the model does not give any false positives or false negatives.
- They should also ensure that the model is robust to adversarial attacks.
- They should also ensure that the model is free of any privacy leaks.
- They should make sure that the model does respond to the explicit and implicit requirements that were requested in the prompt, which is also called IF
- The model should be having a good truthfulness rating and high performance rating (accuracy, precision, recall, F1 score, etc.)
- Upon deployment, the model should have less latency and a high through-put.
- The model is observed and monitored on how it behaved during the training data and how it compares to the validation data which is unseen data which was not used during the training.
- The model should be tested on the test data which is also unseen data which was not used during the training.
- The model should also be tested on the production data which is the real world data.
- The model should also be tested on the adversarial data which is the data that is designed to fool the model.
- There are hence mechanisms such as RLHF, which is called Reinforcement Learning from Human Feedback, which is used to improve the model's performance.
- The model should also be tested on the data that is not in the training data but is similar to the training data. This is called the out of distribution data.

## Tasks

- Some of the tasks they perform is Model quantization, Model pruning, Model distillation, Model compression, Model optimization, Model deployment, Model monitoring, Model maintenance, Model documentation, Model communication, Model interpretability, Model fairness, Model robustness, Model security, Model explainability, Model transparency, Model accountability, Model auditing, Model compliance, Model governance

## Challenges

- The challenges during model training is overfittinbg and underfitting.
- Overfitting is when the model performs well on the training data but poorly on the validation data. This is because the model has learned the noise in the training data.
- Underfitting is when the model performs poorly on both the training data and the validation data. This is because the model is too simple to capture the underlying patterns in the data. This is because the model didn't go too much in-depth on the data-sets and hence, it didn't learn the patterns in the data.
- A correct balance is training the model somewhere between too deep and too shallow so that it learns the patterns but also generalizes well to unseen data and does not overfit ( train on the noise )
- The model should be able to generalize well to unseen data.
- The model should not be generalizing, characterizing and learning from the biases in the data. Also, it should not be stereotypical of training datasets as it should generate the same biased responses especially for fields like healthcare, finance that require a high degree of fairness and accuracy.

## Architectures and Algorithms

- The most common architectures and algorithms used in Machine Learning are:
  - Linear Regression
  - Logistic Regression
  - Decision Trees
  - Random Forests
  - Support Vector Machines
  - K-Nearest Neighbors
  - K-Means Clustering
  - Principal Component Analysis
  - Neural Networks
  - Convolutional Neural Networks
  - Recurrent Neural Networks
  - Long Short-Term Memory
  - Generative Adversarial Networks
  - Transformers
  - Reinforcement Learning
  - Q-Learning
  - Deep Q-Networks
  - Proximal Policy Optimization
  - Actor-Critic Methods

# MLOps Engineers

## Tech Stack they use

- After the model is trained, it is deployed to production. This is done by MLOps Engineers. They use the following tech stack:
  - Docker
  - Kubernetes
  - TensorFlow Serving
  - TorchServe
  - MLflow
  - Kubeflow
  - Seldon
  - TFX
  - Kubeflow Pipelines
  - Apache Beam
  - Apache Spark
  - Apache Flink
  - Apache Kafka
  - Apache Pulsar
  - Apache Airflow
  - Prefect
  - Dask
  - Ray
  - DVC
  - Weights & Biases
  - Comet.ml
  - Neptune.ai

## Duties and Responsibilities

- After the model is deployed, MLOps engineers monitor the following related to the model after it is deployed
  i.e. Inference time, throughput, latency, accuracy, precision, recall, F1 score, ROC-AUC, etc.
- They also monitor the data drift, concept drift, and model drift.
- They also monitor the system health, resource utilization, and system logs.
- They also monitor the following

  - model versioning
  - model deployment
  - model rollback
  - model retraining
  - model monitoring
  - model explanation
  - model interpretability
  - model fairness
  - model accountability
  - model transparency,
  - model compliance
  - model ethics

  - MLOps Engineers define ML pipelines that perform tasks such as deploying the model to Production using ML workflow related pipelines. The pipeline includes steps that include
    - Model is trained and validated.
    - Model is deployed to Production.
    - Model is monitored in Production.
    - Model is retrained and updated in Production when necessary ( especially when data drift, concept drift, or model drift is detected)
    - Perform Canary deployments
    - Perform A/B testing.
    - Configuring rollback strategies in the ML pipeline to rollback to a previous version of the model if necessary.
    - Monitoring the model in Production and alerting the team if any anomalies are detected.
    - Configuring the model to handle new data and new features.
    - Monitoring and Observability using Prometheus / Grafana to ensure that there is not alot of performance degradation of model in Production.
    - Comparing the Model's performance as it was during training vs the way it is in production, which could send the model back for re-training if it's not seen fit for real world data. This is called "Data Drift". Also, the scripts in the pipeline adhere to Canary deployments for this reason as well so that the model is not immediately affected by new data.
    - Also, a Model that does not perform well in Production leads to customer churn as well and hence immediately restoring the model to a previous version that was performing well is necessary.
    - Model Explainability and Interpretability.
    - Model Fairness and Accountability.
    - Model versioning and tracking.
    - Model compliance and ethics.
    - When a new version of the model is released, it is important to ensure that the model is compatible with the existing version of the model.
    - Also, when it is time to retire a model, the pipeline executes A/B tests to ensure that the new model performs better than the old model and only then does it retire the old model.
    - At some point in time, the model will need to be retrained and updated. The pipeline will need to be updated to reflect the new model.
    - The pipeline will also need to be updated to reflect new data and new features.
    - But, after many iterations, the data the old model has trained on becomes very noisy and newer models often perform better. So, the pipeline will need to be updated to reflect the new model.
    - The pipeline will also need to be updated to reflect new data and new features.

## Percentage of time they spend on each task

    - 20% of the time they spend on model deployment and monitoring.
    - 20% of the time they spend on model retraining and updating.
    - 20% of the time they spend on model fairness and accountability.
    - 20% of the time they spend on model compliance and ethics.
    - 20% of the time they spend on model versioning and tracking.

## Challenges

    - Model Drift
    - Data Drift
    - Concept Drift
    - Model Explainability and Interpretability
    - Model Fairness and Accountability
    - Model Compliance and Ethics
    - Mitigating any biases from the Model and ensuring that it does not generate harmful responses
    - Truthfulness of the Model and ensuring that it does not generate false responses.
    - Ensuring that the Model is not generating harmful responses.
    - Latency, Throughput and how well the Model scales in Production
