#!/usr/bin/env python3
import os
import posixpath
import ssl
from http import HTTPStatus
from http.server import HTTPServer, SimpleHTTPRequestHandler
from socketserver import ThreadingMixIn
from urllib.parse import unquote, urlsplit


ROOT = os.environ.get("AI_STUDY_FILE_ROOT", "/data/file")
HOST = os.environ.get("AI_STUDY_FILE_HOST", "0.0.0.0")
PORT = int(os.environ.get("AI_STUDY_FILE_PORT", "443"))
CERT = os.environ.get("AI_STUDY_FILE_CERT", "/data/file/certs/collect-ui.top.crt")
KEY = os.environ.get("AI_STUDY_FILE_KEY", "/data/file/certs/collect-ui.top.key")


class StaticHandler(SimpleHTTPRequestHandler):
  extensions_map = {
    **SimpleHTTPRequestHandler.extensions_map,
    ".js": "application/javascript",
    ".json": "application/json",
    ".mjs": "application/javascript",
    ".svg": "image/svg+xml",
    ".webp": "image/webp"
  }

  def end_headers(self):
    self.send_header("Access-Control-Allow-Origin", "*")
    self.send_header("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
    self.send_header("Access-Control-Allow-Headers", "Content-Type, Range")
    if self.path.startswith("/ai-study/assets/"):
      self.send_header("Cache-Control", "public, max-age=31536000, immutable")
    else:
      self.send_header("Cache-Control", "public, max-age=300")
    super().end_headers()

  def list_directory(self, path):
    self.send_error(HTTPStatus.NOT_FOUND, "File not found")
    return None

  def send_head(self):
    request_path = posixpath.normpath(unquote(urlsplit(self.path).path))
    blocked = (
      request_path == "/bin" or
      request_path.startswith("/bin/") or
      request_path == "/certs" or
      request_path.startswith("/certs/") or
      "/." in request_path
    )
    if blocked:
      self.send_error(HTTPStatus.NOT_FOUND, "File not found")
      return None
    return super().send_head()

  def do_OPTIONS(self):
    self.send_response(HTTPStatus.NO_CONTENT)
    self.end_headers()


class ThreadingHTTPSServer(ThreadingMixIn, HTTPServer):
  allow_reuse_address = True
  daemon_threads = True
  request_queue_size = 128

  def __init__(self, server_address, handler_class, context):
    self.ssl_context = context
    super().__init__(server_address, handler_class)

  def get_request(self):
    socket, address = self.socket.accept()
    socket.settimeout(10)
    secure_socket = self.ssl_context.wrap_socket(
      socket,
      server_side=True,
      do_handshake_on_connect=False
    )
    secure_socket.settimeout(10)
    return secure_socket, address


def main():
  os.chdir(ROOT)
  context = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
  context.load_cert_chain(certfile=CERT, keyfile=KEY)
  server = ThreadingHTTPSServer((HOST, PORT), StaticHandler, context)
  print(f"serving {ROOT} on https://{HOST}:{PORT}", flush=True)
  server.serve_forever()


if __name__ == "__main__":
  main()
