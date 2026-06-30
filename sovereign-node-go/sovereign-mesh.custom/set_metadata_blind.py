import sys
import os
from google.oauth2 import service_account
from google.cloud import compute_v1


def set_metadata_blind(project_id, zone, instance_name, key_file, ssh_key_string):
    credentials = service_account.Credentials.from_service_account_file(key_file)
    client = compute_v1.InstancesClient(credentials=credentials)

    print(f"Blindly setting metadata for {instance_name}...")

    # Try to set without fingerprint
    metadata = compute_v1.types.Metadata(
        items=[compute_v1.types.Items(key="ssh-keys", value=ssh_key_string)]
    )

    try:
        operation = client.set_metadata(
            project=project_id,
            zone=zone,
            instance=instance_name,
            metadata_resource=metadata,
        )
        print(f"Operation {operation.name} started...")
        operation.result()
        print("Metadata set successfully.")
    except Exception as e:
        print(f"Set Metadata Error: {e}")


if __name__ == "__main__":
    project = "model-loader-495607-m2"
    zone = "us-central1-c"
    instance = "instance-20260507-075600"
    key_path = "gcp-key.json"
    new_ssh_key = "billing:ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOab0/cxSWK+wR0iy/ckEIjDSYpAYx0Ivsh8jMUjcV7I root@yogapqrinfo"

    set_metadata_blind(project, zone, instance, key_path, new_ssh_key)
