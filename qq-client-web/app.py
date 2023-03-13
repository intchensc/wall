import requests
from flask import Flask, request, url_for, redirect, flash

app = Flask(__name__)


@app.route('/', methods=['GET', 'POST'])
def hello_world():  # put application's code here
    if request.method == "GET":
        code = request.args.get("code")
        response2 = requests.get(url='https://api.q.qq.com/sns/jscode2session',
                                 params={"appid": "1112206385", "secret": "zj1zTHoYqtyhaUHh",
                                         "js_code": code, "grant_type": "authorization_code"})
        print(response2.json())
        return response2.json()


if __name__ == '__main__':
    app.run()
