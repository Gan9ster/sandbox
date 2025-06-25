#!/usr/bin/env python3
"""Example client demonstrating how to call the sandbox server."""

import json
import http.client


CODE = """def add(a, b):
    return a + b
"""

TESTS = """import unittest
from solution import add

class AddTests(unittest.TestCase):
    def test_add(self):
        self.assertEqual(add(2, 3), 5)

if __name__ == '__main__':
    unittest.main()
"""


def main():
    conn = http.client.HTTPConnection('localhost', 8000)
    payload = json.dumps({'code': CODE, 'tests': TESTS})
    headers = {'Content-Type': 'application/json'}
    conn.request('POST', '/run', body=payload, headers=headers)
    resp = conn.getresponse()
    data = resp.read()
    print('Status:', resp.status)
    print('Response:', data.decode())


if __name__ == '__main__':
    main()
