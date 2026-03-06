import pandas as pd
from pathlib import Path


BASE_DIR = Path(__file__).parent
INPUT_PATH = BASE_DIR / "processed" / "mediPath-data.csv"
OUTPUT_PATH = BASE_DIR / "processed" / "cleaned_data.csv"


def clean_dataset():

    print("Loading dataset...")

    df = pd.read_csv(INPUT_PATH)

    print("Initial rows:", len(df))

    # Fill missing values
    df = df.fillna("")

    # Combine useful text fields
    df["text_features"] = (
        df["description"] + " " +
        df["indication"] + " " +
        df["mechanism"] + " " +
        df["pharmacodynamics"] + " " +
        df["targets"] + " " +
        df["categories"] + " " +
        df["atc_codes"]
    )

    # Remove rows with too little information
    df = df[df["text_features"].str.len() > 50]

    df = df.reset_index(drop=True)

    print("Rows after cleaning:", len(df))

    OUTPUT_PATH.parent.mkdir(parents=True, exist_ok=True)
    df.to_csv(OUTPUT_PATH, index=False)

    print("Saved cleaned dataset →", OUTPUT_PATH)
    return  OUTPUT_PATH


if __name__ == "__main__":
    clean_dataset()