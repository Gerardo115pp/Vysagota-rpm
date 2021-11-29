from subprocess import DEVNULL
from time import sleep
import subprocess
import threading
import os, sys

class ServerRunner:
    def __init__(self) -> None:
        self.react_server = None
        self.flask_server = None
        self.react_process = None
        self.flask_process = None
        
    def startServer(self) -> None:
        self.react_server = threading.Thread(target=self.reactServer)
        self.flask_server = threading.Thread(target=self.startFlaskServer)
        self.flask_server.setDaemon(True)
        self.react_server.setDaemon(True)
        self.react_server.start()
        self.flask_server.start()
        os.system("start msedge.exe http://localhost:3000")
        self.react_process.communicate()[0]
        self.flask_process.communicate()[0]
        
    def reactServer(self,build_path="./alz-client/build/") -> None:
        print("Starting azh-client...", end="\n")
        self.react_process = subprocess.Popen(["serve","-l","3000", "-s", build_path], shell=True, stderr=DEVNULL)
        print("done")
        sleep(0.3)
        
    def startFlaskServer(self) -> None:
        print("Starting diagnostic server...")
        self.flask_process = subprocess.Popen(["python", "alzhaimer-server.py" ], shell=True)
        
    
if __name__ == "__main__":
    server = ServerRunner()
    server.startServer()
