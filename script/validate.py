import sys
import os
import datetime
from subprocess import call
from pymongo import MongoClient

type_object = sys.argv[1]
file_loc = sys.argv[2]
[username, hexid, filename] = file_loc.split('+++')
send_notif = True

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
print("[+] file set to visible")
collection.update_one({'hexid': hexid}, { '$set': {'visible': True}})

if type_object == "solution":
        db.crackme.update_one({'_id': db_object["crackmeid"]}, {'$inc': {"nbsolutions": 1}})

call(["mv", file_loc, filename])
print("[+] mv " + file_loc + " " + filename)
call(["zip", "-j", "--password", "crackmes.one" , "/home/crackmesone/go/src/github.com/5tanislas/crackmes.one/static/" + type_object + "/" + hexid, filename])
print("[+] zip -j --password crackmes.one /home/crackmesone/go/src/github.com/5tanislas/crackmes.one/static/" + type_object + "/" + hexid + " " + filename)
call(["rm", filename])
print("[+] rm " + filename)

if send_notif:
    print("[+] Sending " + type_object + " approval notification!")
    notif_coll = db.notifications
    author_name = db_object["author"]
    if type_object == "solution":
        crackme_obj = db.crackme.find_one({'_id': db_object["crackmeid"]})
        ins_id = notif_coll.insert_one({"user": author_name, "time": datetime.datetime.utcnow(), "seen": False, \
                "text": "Your solution for '" + crackme_obj["name"] + "' has been accepted!"}).inserted_id
        # Set HexId here too for this case
        notif_coll.find_one_and_update({'_id': ins_id}, {'$set': {'hexid': str(ins_id)}})
        ins_id = notif_coll.insert_one({"user": crackme_obj["author"], "time": datetime.datetime.utcnow(), "seen": False, \
                "text": "A new solution for your crackme '" + crackme_obj["name"] \
                + "' has been submitted by: " + author_name}).inserted_id
    elif type_object == "crackme":
        ins_id = notif_coll.insert_one({"user": author_name, "time": datetime.datetime.utcnow(), "seen": False, \
                "text": "Your crackme '" + db_object["name"] + "' has been accepted!"}).inserted_id
    # Set HexId here
    notif_coll.find_one_and_update({'_id': ins_id}, {'$set': {'hexid': str(ins_id)}})
