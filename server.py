#!/usr/bin/env python3
"""Simple sandbox server for running agent-generated code with tests."""

from http.server import BaseHTTPRequestHandler, HTTPServer
import json
import tempfile
import subprocess
import os
from pathlib import Path
import resource


def _limit_resources():
    """Apply basic resource limits to the sandbox process."""
    # Limit CPU time to 2 seconds.
    resource.setrlimit(resource.RLIMIT_CPU, (2, 2))
    # Limit address space to 100 MB.
    mem = 100 * 1024 * 1024
    resource.setrlimit(resource.RLIMIT_AS, (mem, mem))


class SandboxHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        if self.path != '/run':
            self.send_response(404)
            self.end_headers()
            self.wfile.write(b'Not Found')
            return
        length = int(self.headers.get('Content-Length', 0))
        data = self.rfile.read(length)
        try:
            payload = json.loads(data)
            code = payload.get('code', '')
            tests = payload.get('tests', '')
        except json.JSONDecodeError:
            self.send_response(400)
            self.end_headers()
            self.wfile.write(b'Invalid JSON')
            return

        with tempfile.TemporaryDirectory() as tmpdir:
            code_path = Path(tmpdir) / 'solution.py'
            test_path = Path(tmpdir) / 'test_solution.py'
            code_path.write_text(code)
            test_path.write_text(tests)

            try:
                proc = subprocess.run(
                ['python3', '-m', 'unittest', test_path.stem],
                cwd=tmpdir,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
                text=True,
                preexec_fn=_limit_resources,
                timeout=5,
                env={'PYTHONPATH': tmpdir}
                )
                response = {
                    'returncode': proc.returncode,
                    'output': proc.stdout
                }
            except subprocess.TimeoutExpired as exc:
                response = {
                    'returncode': -1,
                    'output': f'Timeout after {exc.timeout} seconds'
                }

        body = json.dumps(response).encode()
        self.send_response(200)
        self.send_header('Content-Type', 'application/json')
        self.send_header('Content-Length', str(len(body)))
        self.end_headers()
        self.wfile.write(body)


def run(server_class=HTTPServer, handler_class=SandboxHandler, port=8000):
    server_address = ('', port)
    httpd = server_class(server_address, handler_class)
    print(f"Starting sandbox server on port {port}...")
    httpd.serve_forever()


if __name__ == '__main__':
    run()
