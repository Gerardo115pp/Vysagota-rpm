from requests_toolbelt.multipart import decoder as multipart_decoder
from keras.preprocessing.image import img_to_array
from keras.models import load_model
from socket import socket
from PIL import Image
from io import BytesIO
import http.server as http
import tensorflow as tf
import numpy as np
import requests
import json
import os

tf.compat.v1.logging.set_verbosity(tf.compat.v1.logging.ERROR)

JD_ADDRESS = os.getenv("JD_ADDRESS") or "localhost:8080"
cnn = load_model('./model/model.h5')
cnn.load_weights('./model/weights.h5')

os.system("clear")


class IAServer(http.BaseHTTPRequestHandler):
    @classmethod
    def run(cls, port=8000, host='localhost'):
        server = http.HTTPServer((host, port), cls)
        cls.boot(port, host)
        server.serve_forever()
    
    @classmethod
    def boot(cls, port, host) -> bool:
        print(f"Booting server at address: {host}:{port}")
        response = requests.post(
            f'http://{JD_ADDRESS}/new-service', 
            files=(
                ("name", (None, "ia_server")),
                ("host", (None, host)),
                ("port", (None, port)),
                ("dns", (None, "ia_server")),
                ("status", (None, "online")),
                ("protocols", (None, json.dumps(["http"]))),
                ("actions", (None, json.dumps({"diagnose": {"name": "diagnose", "methods": ["POST"]}})))
            )
        )
        if response.status_code > 299:
            print(f"Error booting server: {response.text}")
            return False
        
        print(f"Server booted successfully")
        return True
        
    def __init__(self, request: socket, client_address: tuple[str, int], server: http.HTTPServer) -> None:
        super().__init__(request, client_address, server)
    
    def diagnose(self):
        image:np.ndarray = self.loadImage()
        prediction:int = int(self.predict(image))
        print(f"Prediction: {prediction}")
        response_body:str = json.dumps({"diagnosis": prediction})
        self.send_response(200)
        self.setCorsPolicy()
        self.end_headers()
        self.wfile.write(response_body.encode())
    
    
    def do_GET(self):
        self.send_response(405)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        self.wfile.write(b'Hello World')
        
    def do_POST(self):
        if self.path == '/diagnose':
            self.diagnose()
    
    def loadImage(self) -> np.ndarray:
        body_data = self.rfile.read(int(self.headers['Content-Length']))
        form = multipart_decoder.MultipartDecoder(body_data, self.headers['Content-Type'])
        image = Image.open(BytesIO(form.parts[0].content))
        if image.mode != 'RGB':
            image = image.convert('RGB')
        image = image.resize((200, 200))
        image_array = img_to_array(image)
        image_array = np.expand_dims(image_array, axis=0) # image_array.shape = (1, width, height, 3)
        return image_array

    def predict(self, image:np.ndarray) -> int:
        global cnn
        result = cnn.predict(image)
        return result[0][0]

    def setCorsPolicy(self) -> None:
        self.send_header("Access-Control-Allow-Origin", "*")
        self.send_header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
        self.send_header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        
        
    
if __name__ == '__main__':
    IAServer.run()