import sys
import os
from subprocess import call
from pymongo import MongoClient

type_object = sys.argv[1]
file_loc = sys.argv[2]

[username, hexid, filename] = file_loc.split('+++')

client = MongoClient('127.0.0.1')
db = client.crackmesone

if type_object == "crackme":
	file_loc = "/home/crackmesone/go/src/github.com/5tanislas/crackmes.one/tmp/crackme/" + file_loc
	collection = db.crackme
elif type_object == "solution":
	file_loc = "/home/crackmesone/go/src/github.com/5tanislas/crackmes.one/tmp/solution/" + file_loc
	collection = db.solution
else:
	print("[-] I don't understand the type")
	sys.exit()

db_object = collection.find_one({'hexid': hexid})

if db_object is None:
	print("not found in db")
	os._exit(0)

print("[+] found in database !")
print(db_object)
print("[+] file deleted in db")
collection.delete_one({'hexid': hexid})

call(["rm", file_loc])
print("[+] rm " + file_loc)
