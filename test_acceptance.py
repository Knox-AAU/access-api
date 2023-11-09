import requests
import unittest


class TestAcceptance(unittest.TestCase):
    url = "http://localhost:8080/test"

    def test_request(self):
        headers = {"Access-Authorization": "internal_key"}
        try:
            response = requests.get(self.url, headers=headers)
            self.assertEqual(response.status_code, 200, "Expected status code 200")
        except requests.RequestException as e:
            self.fail(f"Failed to make the HTTP request: {e}")


if __name__ == "__main__":
    unittest.main()
