from http.server import BaseHTTPRequestHandler
import os
import requests


class Handler(BaseHTTPRequestHandler):

    def do_POST(self):
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()

        content_len = int(self.headers.get('content-length', 0))
        post_body = self.rfile.read(content_len)
        print(post_body.decode("utf-8"))

        # send to openai
        url = "https://api.openai.com/v1/chat/completions"
        headers = {
            "Content-Type": "application/json",
            'Authorization': 'Bearer ' + os.getenv("OPENAI_API_KEY")
        }

#         requests.post(url, stream=True, headers=headers, data=post_body)
#         self.wfile.write(b'{"hello":"world"}')

        with requests.post(url, stream=True, headers=headers, data=post_body) as response:
            for chunk in response.iter_content(chunk_size=512):
                self.wfile.write(chunk)
                break
