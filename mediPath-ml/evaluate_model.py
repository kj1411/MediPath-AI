import joblib
from pathlib import Path
from sklearn.metrics import classification_report, accuracy_score
import numpy as np

from feature_engineering import prepare_features


BASE_DIR = Path(__file__).resolve().parent

MODEL_DIR = BASE_DIR / "models"


def evaluate():

    (
        X_train,
        X_test,
        y_train,
        y_test,
        vectorizer,
        label_encoder
    ) = prepare_features()

    model = joblib.load(MODEL_DIR / "disease_model.pkl")

    predictions = model.predict(X_test)

    print("Accuracy:", accuracy_score(y_test, predictions))

    print("\nClassification Report:\n")

    labels = np.unique(y_test)

    print(
        classification_report(
            y_test,
            predictions,
            labels=labels,
            target_names=label_encoder.inverse_transform(labels)
        )
    )


if __name__ == "__main__":
    evaluate()