import os
from flask import Flask

app = Flask(__name__)

@app.route("/answer_service")
def provide_answer():
    return "42"

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5678)