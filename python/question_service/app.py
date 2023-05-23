import os
import requests
from flask import Flask

app = Flask(__name__)

ANSWER_ENDPOINT = os.getenv("ANSWER_ENDPOINT", "http://localhost:5678") + "/answer_service"

@app.route("/question_service")
def question_and_answer():
    question = determine_question()
    answer = get_answer()
    return (question + "\n" + answer + "\n")

def determine_question():
    return "what is the answer to the ultimate question of life, the universe, and everything?"

def get_answer():
    r = requests.get(ANSWER_ENDPOINT)
    return r.text

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=1234)