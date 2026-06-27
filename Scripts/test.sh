```bash
#!/bin/bash



echo "== Health Check =="
curl http://localhost:8080/health
echo -e "\n"

echo "== Create/Update Key =="
curl -X PUT http://localhost:8080/kv/name \
  -H "Content-Type: application/json" \
  -d '{"value":"Ben"}'
echo -e "\n"

echo "== Get Key =="
curl http://localhost:8080/kv/name
echo -e "\n"

echo "== Delete Key =="
curl -X DELETE http://localhost:8080/kv/name
echo -e "\n"

echo "== Get Deleted Key =="
curl http://localhost:8080/kv/name
echo -e "\n"
```
