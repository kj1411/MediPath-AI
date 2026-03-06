import pandas as pd
from pathlib import Path
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import LabelEncoder


BASE_DIR = Path(__file__).resolve().parent

DATA_PATH = BASE_DIR / "mediPath_data" / "processed" / "mediPath-clean.csv"


def extract_atc_label(atc_codes):

    if not isinstance(atc_codes, str) or atc_codes.strip() == "":
        return "Unknown"

    # first ATC code
    first = atc_codes.split(",")[0].strip()

    # first letter = therapeutic group
    return first[0]


def prepare_features():

    df = pd.read_csv(DATA_PATH)

    df = df.fillna("")

    # Create label from ATC code
    df["label"] = df["atc_codes"].apply(extract_atc_label)

    # Remove unknown labels
    df = df[df["label"] != "Unknown"]

    print("Dataset size:", len(df))

    X_text = df["text_features"]
    y = df["label"]

    # Encode labels
    label_encoder = LabelEncoder()
    y_encoded = label_encoder.fit_transform(y)

    vectorizer = TfidfVectorizer(
        max_features=5000,
        stop_words="english"
    )

    X = vectorizer.fit_transform(X_text)

    X_train, X_test, y_train, y_test = train_test_split(
        X,
        y_encoded,
        test_size=0.2,
        random_state=42
    )

    return (
        X_train,
        X_test,
        y_train,
        y_test,
        vectorizer,
        label_encoder
    )