# Import necessary libraries
import yfinance as yf
import pandas as pd
import logging
from sklearn.metrics import mean_squared_error
from math import sqrt

from statsmodels.tsa.arima.model import ARIMA

def fetch_stock_data(ticker, start_date, end_date):
    """
    Fetches historical stock data from Yahoo Finance.

    Args:
        ticker (str): The stock ticker symbol (e.g., AAPL).
        start_date (str): The start date for fetching data (YYYY-MM-DD).
        end_date (str): The end date for fetching data (YYYY-MM-DD).

    Returns:
        pandas.DataFrame: A DataFrame containing the historical stock data, or None if an error occurs.
    """
    try:
        data = yf.download(ticker, start=start_date, end=end_date)
        return data
    except Exception as e:
        logger.error(f"Error fetching data for {ticker}: {e}")
        return None

def preprocess_data(data):
    """
    Preprocesses the stock data.

    Args:
        data (pandas.DataFrame): The DataFrame containing stock data.

    Returns:
        pandas.Series: A Series containing the adjusted closing prices, resampled to daily frequency,
                       with missing values filled using forward fill, or None if input is invalid.
    """
    if data is None or not isinstance(data, pd.DataFrame) or data.empty:
        logger.error("Invalid or empty data provided for preprocessing.")
        return None

    # Use only the 'Adj Close' column and resample to daily frequency
    prices = data['Adj Close'].resample('D').mean()

    # Fill missing values using forward fill
    prices = prices.fillna(method='ffill')

    return prices

def train_arima_model(prices, order=(5, 1, 0)):
    """
    Trains an ARIMA model on the preprocessed stock prices.

    Args:
        prices (pandas.Series): The Series containing the preprocessed stock prices.
        order (tuple): The order (p, d, q) of the ARIMA model. Default is (5, 1, 0).

    Returns:
        statsmodels.tsa.arima.model.ARIMAResults: The trained ARIMA model, or None if an error occurs.
    """
    try:
        model = ARIMA(prices, order=order)
        model_fit = model.fit()
        return model_fit, prices
    except Exception as e:
        logger.error(f"Error training ARIMA model: {e}")
        return None, None

def generate_future_predictions(model, periods):
    """
    Generates future predictions using the trained ARIMA model.

    Args:
        model (statsmodels.tsa.arima.model.ARIMAResults): The trained ARIMA model.
        periods (int): The number of periods to forecast.

    Returns:
        pandas.Series: A Series containing the future predictions.
    """
    
    future_predictions = model.forecast(steps=periods)
    return future_predictions

def format_predictions_with_dates(predictions):
    """Formats the predictions with their corresponding dates."""
    formatted_predictions = []
    for date, prediction in predictions.items():
        formatted_predictions.append((date.strftime('%Y-%m-%d'), prediction))
    return future_predictions

if __name__ == "__main__":
    # Configure logging
    logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
    logger = logging.getLogger(__name__)

    # Specify stock ticker and date range
    ticker = "AAPL"
    end_date = pd.Timestamp.today()
    start_date = end_date - pd.DateOffset(years=5)
    
    logger.info(f"Starting stock analysis for {ticker} from {start_date} to {end_date}")

    # Fetch data
    stock_data = fetch_stock_data(ticker, start_date, end_date)

    if stock_data is not None:
        # Split the data into training and testing sets before preprocessing
        split_index = int(len(stock_data) * 0.8)  # 80% for training
        train_data = stock_data[:split_index]
        test_data = stock_data[split_index:]

        # Preprocess data
        train_prices = preprocess_data(train_data)
        test_prices = preprocess_data(test_data)

        if train_prices is not None and test_prices is not None:
            # Train ARIMA model
            model, prices_serie = train_arima_model(train_prices)

            if model is not None:
                # Print model summary
                logger.info("ARIMA Model Summary:")
                logger.info(model.summary())
                
                # Evaluate the model
                test_predictions = model.predict(start=len(train_prices), end=len(train_prices)+len(test_prices)-1)
                
                if len(test_predictions) != len(test_prices):
                    logger.error(f"Prediction and test sets length differ")
                else:
                    rmse = sqrt(mean_squared_error(test_prices.values, test_predictions.values))
                    logger.info(f"ARIMA Model Evaluation - Root Mean Squared Error (RMSE): {rmse}")

                    logger.info(f"Test Predictions: {test_predictions.values}")
                    # Generate future predictions
                    future_predictions = generate_future_predictions(model, periods=30)
                    formatted_future_predictions = list(zip(future_predictions.index.strftime('%Y-%m-%d'), future_predictions.values))
                    logger.info(f"Future Predictions (next 30 days): {formatted_future_predictions}")


        else:
             logger.error("Error during data preprocess")