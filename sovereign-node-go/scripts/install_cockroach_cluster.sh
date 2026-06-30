#!/bin/bash
set -e

# Download CockroachDB Linux x86_64
if [ ! -f /tmp/cockroach-v23.1.10.linux-amd64.tgz ]; then
  echo 'Downloading CockroachDB...'
  curl -o /tmp/cockroach-v23.1.10.linux-amd64.tgz https://binaries.cockroachdb.com/cockroach-v23.1.10.linux-amd64.tgz
fi

# Extract and install binary to /usr/local/bin
echo 'Extracting CockroachDB...'
tar -C /tmp -xzf /tmp/cockroach-v23.1.10.linux-amd64.tgz
sudo cp /tmp/cockroach-v23.1.10.linux-amd64/cockroach /usr/local/bin/
sudo chmod +x /usr/local/bin/cockroach

# Clean up
rm -rf /tmp/cockroach-v23.1.10.linux-amd64*

echo 'CockroachDB installed successfully.'
