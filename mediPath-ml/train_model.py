import joblib
from pathlib import Path
from sklearn.ensemble import RandomForestClassifier

from feature_engineering import prepare_features


BASE_DIR = Path(__file__).resolve().parent

MODEL_DIR = BASE_DIR / "models"


def train():

    (
        X_train,
        X_test,
        y_train,
        y_test,
        vectorizer,
        label_encoder
    ) = prepare_features()

    print("Training model...")

    model = RandomForestClassifier(
        n_estimators=200,
        random_state=42
    )

    model.fit(X_train, y_train)

    print("Training complete. Saving model...")

    MODEL_DIR.mkdir(parents=True, exist_ok=True)

    joblib.dump(model, MODEL_DIR / "disease_model.pkl")
    joblib.dump(vectorizer, MODEL_DIR / "vectorizer.pkl")
    joblib.dump(label_encoder, MODEL_DIR / "label_encoder.pkl")

    print("Model saved to:", MODEL_DIR)


if __name__ == "__main__":
    train()