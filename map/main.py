from bottle import route, run


@route('/')
def index():
    with open('./client.html', 'r') as file:
        return file.read()


if __name__ == '__main__':
    run(host='0.0.0.0', port=8080)
