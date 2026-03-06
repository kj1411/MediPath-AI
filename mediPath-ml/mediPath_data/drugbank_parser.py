import xml.etree.ElementTree as ET
import pandas as pd
from pathlib import Path


BASE_DIR = Path(__file__).parent

# Path to the raw XML file
RAW_PATH = BASE_DIR / "raw" / "drugbank.xml"
OUTPUT_PATH = BASE_DIR / "processed" / "mediPath-data.csv"



def parse_drugbank():

    print("Loading DrugBank XML...")

    tree = ET.parse(RAW_PATH)
    root = tree.getroot()

    ns = {"db": "http://www.drugbank.ca"}

    drugs = []

    for drug in root.findall("db:drug", ns):

        name = drug.findtext("db:name", default="", namespaces=ns)

        description = drug.findtext(
            "db:description", default="", namespaces=ns
        )

        indication = drug.findtext(
            "db:indication", default="", namespaces=ns
        )

        mechanism = drug.findtext(
            "db:mechanism-of-action", default="", namespaces=ns
        )

        pharmacodynamics = drug.findtext(
            "db:pharmacodynamics", default="", namespaces=ns
        )

        # categories
        categories = []
        for cat in drug.findall("db:categories/db:category/db:category", ns):
            if cat.text:
                categories.append(cat.text)

        # targets
        targets = []
        for target in drug.findall("db:targets/db:target/db:name", ns):
            if target.text:
                targets.append(target.text)

        # ATC codes
        atc_codes = []
        for atc in drug.findall("db:atc-codes/db:atc-code", ns):
            code = atc.attrib.get("code")
            if code:
                atc_codes.append(code)

        drugs.append({
            "drug_name": name.lower(),
            "description": description,
            "indication": indication,
            "mechanism": mechanism,
            "pharmacodynamics": pharmacodynamics,
            "targets": ", ".join(targets),
            "categories": ", ".join(categories),
            "atc_codes": ", ".join(atc_codes)
        })

    df = pd.DataFrame(drugs)

    print("Total drugs parsed:", len(df))

    OUTPUT_PATH.parent.mkdir(parents=True, exist_ok=True)
    df.to_csv(OUTPUT_PATH, index=False)

    print("Saved dataset →", OUTPUT_PATH)
    return OUTPUT_PATH


if __name__ == "__main__":
    parse_drugbank()