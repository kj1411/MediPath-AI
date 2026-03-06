import os
from kaggle.api.kaggle_api_extended import KaggleApi

# Get the folder where this script lives
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
DEFAULT_DOWNLOAD_PATH = os.path.join(BASE_DIR, "raw")

def download_drugbank(dataset="sergeguillemart/drugbank",
                       download_path=DEFAULT_DOWNLOAD_PATH):
    """
    Downloads the Kaggle dataset if it doesn't exist locally.
    Unzips automatically.
    """
    expected_file = os.path.join(download_path, "drugbank.xml")

    if os.path.exists(expected_file):
        print(f"{expected_file} already exists. Skipping download.")
        return expected_file

    os.makedirs(download_path, exist_ok=True)

    api = KaggleApi()
    api.authenticate()

    print(f"Downloading Kaggle dataset {dataset} to {download_path}...")
    api.dataset_download_files(dataset, path=download_path, unzip=True, quiet=False)

    if os.path.exists(expected_file):
        print("Download complete!")
        return expected_file
    else:
        raise FileNotFoundError(f"Expected file not found after download: {expected_file}")


if __name__ == "__main__":
    download_drugbank()