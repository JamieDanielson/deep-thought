import os
import requests
from flask import Flask

app = Flask(__name__)

QUESTION_ENDPOINT = os.getenv("QUESTION_ENDPOINT", "http://localhost:1234") + "/question_service"

@app.route("/gateway_service")
def do_the_thing():
    return get_question()

def get_question():
    r = requests.get(QUESTION_ENDPOINT)
    return r.text

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=4242)