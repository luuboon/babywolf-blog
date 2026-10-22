# Simula los dos endpoints de Supabase Storage que usa el backend.
import json, sys
from http.server import BaseHTTPRequestHandler, HTTPServer
class H(BaseHTTPRequestHandler):
    def do_POST(self):
        n = int(self.headers.get('Content-Length') or 0); self.rfile.read(n)
        if self.headers.get('apikey') != 'service-key':
            self.send_response(401); self.end_headers(); return
        p = self.path
        if p.startswith('/storage/v1/object/upload/sign/'):
            body = {"url": p[len('/storage/v1'):] + "?token=subida123"}
        elif p.startswith('/storage/v1/object/sign/'):
            body = {"signedURL": p[len('/storage/v1'):] + "?token=lectura456"}
        else:
            self.send_response(404); self.end_headers(); return
        b = json.dumps(body).encode()
        self.send_response(200); self.send_header('Content-Type','application/json'); self.end_headers(); self.wfile.write(b)
    def log_message(self, *a): pass
HTTPServer(('127.0.0.1', int(sys.argv[1])), H).serve_forever()
