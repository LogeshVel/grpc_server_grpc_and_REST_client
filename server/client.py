import requests

# res = requests.get("http://localhost:51151/v1/emp/1")
# print(res)
# a = {
#   "id": "1",
#   "first_name": "string",
#   "last_name": "string",
#   "role": "string",
#   "contact": {
#     "home_addr": "string",
#     "mob_num": "string",
#     "mail_id": "string"
#   }
# }
# res = requests.post("http://localhost:51151/v1/emp",json=a)
# print(res)

res = requests.get("http://localhost:51151/v1/emp/1")
print(res.json())
aa = {
  "emp": {
    "id": "1",
    "first_name": "patch string",
    "contact": {
      "home_addr": "patch string"
    }
  },
  "update_mask": {
    "paths": [
      "emp.first_name",
      "emp.contact.home_addr"
    ]
  }
}
pat = requests.patch("http://localhost:51151/v1/emp/1",json=aa)
print(pat)
print(pat.json())

res = requests.get("http://localhost:51151/v1/emp/1")
print(res.json())