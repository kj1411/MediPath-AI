import sys
import json
import joblib
from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent
MODEL_DIR = BASE_DIR / "models"

try:
    model = joblib.load(MODEL_DIR / "disease_model.pkl")
    vectorizer = joblib.load(MODEL_DIR / "vectorizer.pkl")
    encoder = joblib.load(MODEL_DIR / "label_encoder.pkl")

    drugs = sys.argv[1:]

    text = " ".join(drugs)

    X = vectorizer.transform([text])

    pred = model.predict(X)[0]
    prob = model.predict_proba(X).max()

    disease = encoder.inverse_transform([pred])[0]

    print(json.dumps({
        "disease": disease,
        "confidence": float(prob)
    }))

except Exception as e:

    print(json.dumps({
        "disease": "Unknown",
        "confidence": 0,
        "error": str(e)
    }))