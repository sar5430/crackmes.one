import sys
import os
import datetime
from subprocess import call
from pymongo import MongoClient

type_object = sys.argv[1]
file_loc = sys.argv[2]

[username, hexid, filename] = file_loc.split('+++')
send_notif = True
rej_reason = None
if send_notif and len(sys.argv) >= 4:
    rej_reason = sys.argv[3]

client = MongoClient('127.0.0.1')
db = client.crackmesone

if type_object == "crackme":
	file_loc = "/home/crackmesone/go/src/github.com/5tanislas/crackmes.one/tmp/crackme/" + file_loc
	collection = db.crackme
	rating_diff = db.rating_difficulty
	rating_qual = db.rating_quality
	
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

collection.delete_one({'hexid': hexid})
print("[+] file deleted in db")

if type_object == "crackme":
	rating_diff.delete_many({"crackmehexid": hexid})
	rating_qual.delete_many({"crackmehexid": hexid})

call(["rm", file_loc])
print("[+] rm " + file_loc)

if send_notif:
    print("[+] Sending " + type_object + " rejection notification!")
    notif_coll = db.notifications
    author_name = db_object["author"]
    if type_object == "solution":
        crackme_obj = db.crackme.find_one({'_id': db_object["crackmeid"]})
        notif_text = "Your solution for '" + crackme_obj["name"] + "' has been rejected!"
        if rej_reason is not None:
            notif_text += " Reason: " + rej_reason
        ins_id = notif_coll.insert_one({"user": author_name, "time": datetime.datetime.utcnow(), "seen": False, "text": notif_text}).inserted_id
    elif type_object == "crackme":
        notif_text = "Your crackme '" + db_object["name"] + "' has been rejected!"
        if rej_reason is not None:
            notif_text += " Reason: " + rej_reason
        ins_id = notif_coll.insert_one({"user": author_name, "time": datetime.datetime.utcnow(), "seen": False, "text": notif_text}).inserted_id
    # Set HexId here
    notif_coll.find_one_and_update({'_id': ins_id}, {'$set': {'hexid': str(ins_id)}})
