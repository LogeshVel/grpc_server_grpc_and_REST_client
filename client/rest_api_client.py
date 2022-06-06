import requests
print("Get before POST")
res = requests.get("http://localhost:51151/v1/emp/1")
print(res)
a = {
  "id": "1",
  "first_name": "string",
  "last_name": "string",
  "role": "string",
  "contact": {
    "home_addr": "string",
    "mob_num": "string",
    "mail_id": "string"
  }
}
print("POST")
res = requests.post("http://localhost:51151/v1/emp",json=a)
print(res)
print("GET after POST")
res = requests.get("http://localhost:51151/v1/emp/1")
print(res.json())
print("PATCH")
aa = {
  "emp": {
    "id": "1",
    "first_name": "patch string",
    "contact": {
      "home_addr": "patch string"
    }
  },
  "update_mask": 
      "emp.firstName,emp.contact.homeAddr"
}
pat = requests.patch("http://localhost:51151/v1/emp/1",json=aa)
print(pat)
print(pat.json())
print("GET after PATCH")
res = requests.get("http://localhost:51151/v1/emp/1")
print(res.json())