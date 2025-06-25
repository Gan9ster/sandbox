# Sandbox Server

This repository provides a simple sandbox server that can execute
agent-generated Python code along with accompanying tests. The server
is implemented using only the Python standard library, making it
suitable for offline environments.

Execution happens in a subprocess with basic CPU and memory limits
applied (about 2 CPU seconds and 100 MB of memory). The process runs
with a minimal environment where only the temporary directory is added
to `PYTHONPATH` so external packages are not available.

## Usage

1. Run the server:

   ```bash
   python3 server.py
   ```

   The server listens on `localhost:8000` by default.

2. Send code and tests to the `/run` endpoint. An example client is
   provided:

   ```bash
   python3 example_client.py
   ```

   The client sends a small code snippet and unit test to the server
   and prints the resulting JSON response.

The server executes the tests inside a temporary directory and returns
both the test output and exit code.
