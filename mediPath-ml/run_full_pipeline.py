import os
import shutil
from mediPath_data import run_pipeline
from train_model import train
from evaluate_model import evaluate

def main():
    try:
        # -----------------------------
        # Step 1: Run data pipeline
        # -----------------------------
        print("Running data pipeline (download, parse, clean)...")
        cleaned_data_path = run_pipeline.main()
        print(f"Data pipeline completed. Cleaned data at: {cleaned_data_path}")

        # -----------------------------
        # Step 2: Train model
        # -----------------------------
        print("Training model...")
        model_path = train()
        print(f"Model trained and saved at: {model_path}")

        # -----------------------------
        # Step 3: Evaluate model
        # -----------------------------
        print("Evaluating model...")
        evaluation_path = evaluate()
        print(f"Evaluation completed. Results at: {evaluation_path}")

    finally:
        # -----------------------------
        # Step 4: Delete raw dataset to save space
        # -----------------------------
        raw_folder = os.path.join("mediPath_data", "raw")
        if os.path.exists(raw_folder):
            shutil.rmtree(raw_folder)
            print(f"Deleted raw dataset folder: {raw_folder}")

    print("Full pipeline completed successfully!")

if __name__ == "__main__":
    main()