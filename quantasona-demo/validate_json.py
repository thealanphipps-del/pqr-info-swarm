import json
import logging
import sys

logging.basicConfig(level=logging.INFO, format='%(levelname)s: %(message)s')

def validate():
    file_path = 'anonymized_real_voice_prints.json'
    
    try:
        with open(file_path, 'r') as f:
            data = json.load(f)
    except Exception as e:
        logging.error(f"Failed to load JSON file: {e}")
        sys.exit(1)
        
    required_fields = {
        "sigmaId": str,
        "semanticWeight": float,
        "confidence": float,
        "provenance": str,
        "deltaType": str,
        "relationType": str,
        "sourceMic": str
    }
    
    total_records = len(data)
    errors = 0
    
    logging.info(f"Starting validation for {file_path}")
    logging.info(f"Total records to process: {total_records}")
    
    for idx, record in enumerate(data):
        for field, expected_type in required_fields.items():
            if field not in record:
                logging.error(f"Record {idx}: Missing required field '{field}'")
                errors += 1
                continue
            
            value = record[field]
            
            if not isinstance(value, expected_type):
                # handle python's bool being a subclass of int, but we just want float/int
                if expected_type == float and isinstance(value, (int, float)) and not isinstance(value, bool):
                    pass
                else:
                    logging.error(f"Record {idx}: Type mismatch for '{field}'. Expected {expected_type.__name__}, got {type(value).__name__}")
                    errors += 1
                    
            if field == 'confidence':
                if isinstance(value, (int, float)):
                    if not (0 <= value <= 1):
                        logging.error(f"Record {idx}: Field 'confidence' out of bounds [0, 1]. Got {value}")
                        errors += 1
                    
    if errors == 0:
        logging.info("Validation successful. Zero instances of type mismatches or missing critical fields.")
        logging.info("All records adhere to the schema.")
    else:
        logging.error(f"Validation failed with {errors} errors.")
        sys.exit(1)
        
if __name__ == '__main__':
    validate()
