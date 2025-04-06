FROM python:3.9-slim

WORKDIR /app

# Copy requirements file
COPY requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Create directories for uploads and models
RUN mkdir -p uploads models

# Set environment variables
ENV PYTHONPATH=/app \
    HOST=0.0.0.0 \
    PORT=8000 \
    DEBUG=False \
    API_KEY=your-api-key \
    MODEL_PATH=/app/models \
    UPLOAD_DIR=/app/uploads

# Expose port
EXPOSE 8000

# Run the application
CMD ["uvicorn", "src.main:app", "--host", "0.0.0.0", "--port", "8000"]
