URL="http://localhost:8002/memaccess"
BODY_TEMPLATE='{
    "type": "Out",
    "numpage": 0,
    "Content": 20,
    "pid": 1,
    "offset": OFFSET_VALUE
}'


make_put_request() {
    local offset=$1
    local body=$(echo "$BODY_TEMPLATE" | sed "s/OFFSET_VALUE/$offset/")

    curl -X PUT "$URL" -H "Content-Type: application/json" -d "$body"
    
}


for offset in {0..15}
do
    make_put_request $offset
done