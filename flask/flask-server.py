from flask import Flask
from flask_restful import Resource, Api, reqparse
from datetime import datetime

app = Flask(__name__)
api = Api(app)

usersDB = {
	1: [1, "admin", "admin"],
	2: [2, "marc", "user"],
	3: [3, "jordi", "user"]
}

class User(Resource):
    def get(self):
        parser = reqparse.RequestParser()
        parser.add_argument('id')
        args = parser.parse_args()
        userID = int(args.get('id'))
        if userID in usersDB: 
            user = usersDB[userID]
            return {"ID": user[0], "Name": user[1], "Role": user[2]}
        else:
            return {"Error": "Invalid user id."}

class Greet(Resource):
    def get(self):
    	return {"Message": "Hello!"}

class CurrentTime(Resource):
    def get(self):
    	now = datetime.today()
    	date_time = now.strftime("%m/%d/%Y, %H:%M:%S")
    	return {'time': date_time}

api.add_resource(User, '/user')
api.add_resource(Greet, '/greet')
api.add_resource(CurrentTime, '/time')

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0')
