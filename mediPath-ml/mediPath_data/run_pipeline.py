import os
from mediPath_data.kaggle_downloader import download_drugbank
from mediPath_data.drugbank_parser import parse_drugbank
from mediPath_data.clean_data import clean_dataset

def main():
    # Step 1: Download dataset if not exists
    data_path = download_drugbank()
    print(f"Dataset ready at: {data_path}")

    # Step 2: Parse the raw dataset
    parsed_data_path = parse_drugbank()
    print(f"Parsed data saved at: {parsed_data_path}")

    # Step 3: Clean the parsed data
    cleaned_data_path = clean_dataset()
    print(f"Cleaned data saved at: {cleaned_data_path}")

    print("Pipeline completed successfully!")

if __name__ == "__main__":
    main()