
TEMP_ACCESS_TOKEN=$(curl -X POST -H "Content-Type: application/json" \
                        -d '{"email": "admin@example.com", "password": "d1r3ctu5"}' \
                        $DIRECTUS_URL/auth/login \
                        | jq .data.access_token | cut -d '"' -f2)

USER_ID=$(curl -X GET -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TEMP_ACCESS_TOKEN" \
    $DIRECTUS_URL/users/me | jq .data.id | cut -d '"' -f2)

curl -X PATCH -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TEMP_ACCESS_TOKEN" \
    -d "{\"token\": \"$ADMIN_ACCESS_TOKEN\"}" \
    $DIRECTUS_URL/users/$USER_ID

# modbot_settings table
curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"collection":"modbot_settings","fields":[{"field":"chat_id","type":"integer","meta":{"hidden":true,"interface":"input","readonly":true},"schema":{"is_primary_key":true,"has_auto_increment":true}},{"field":"date_created","type":"timestamp","meta":{"special":["date-created"],"interface":"datetime","readonly":true,"hidden":true,"width":"half","display":"datetime","display_options":{"relative":true}},"schema":{}},{"field":"date_updated","type":"timestamp","meta":{"special":["date-updated"],"interface":"datetime","readonly":true,"hidden":true,"width":"half","display":"datetime","display_options":{"relative":true}},"schema":{}}],"schema":{},"meta":{"singleton":false}}' \
    $DIRECTUS_URL/collections

# modbot_settings fields
curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"type":"integer","meta":{"interface":"input","special":null,"required":true,"options":{"min":10},"validation":{"_and":[{"expiry_time":{"_gte":"10"}}]}},"field":"expiry_time"}' \
    $DIRECTUS_URL/fields/modbot_settings \

curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"type":"integer","meta":{"interface":"input","special":null,"required":true,"options":{"min":1},"validation":{"_and":[{"threshold":{"_gte":"1"}}]}},"field":"threshold"}' \
    $DIRECTUS_URL/fields/modbot_settings \

# modbot_polls table
curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"collection":"modbot_polls","fields":[{"field":"poll_id","type":"string","meta":{"interface":"input","readonly":false,"hidden":false},"schema":{"is_primary_key":true,"length":255,"has_auto_increment":false}},{"field":"date_created","type":"timestamp","meta":{"special":["date-created"],"interface":"datetime","readonly":true,"hidden":true,"width":"half","display":"datetime","display_options":{"relative":true}},"schema":{}},{"field":"date_updated","type":"timestamp","meta":{"special":["date-updated"],"interface":"datetime","readonly":true,"hidden":true,"width":"half","display":"datetime","display_options":{"relative":true}},"schema":{}}],"schema":{},"meta":{"singleton":false}}' \
    $DIRECTUS_URL/collections

# modbot_polls fields
curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"type":"integer","meta":{"interface":"input","special":null,"required":true},"field":"message_id"}' \
    $DIRECTUS_URL/fields/modbot_polls \

curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"type":"integer","meta":{"interface":"select-dropdown-m2o","special":["m2o"],"required":true},"field":"chat_id"}' \
    $DIRECTUS_URL/fields/modbot_polls \

# modbot_polls relations
curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"collection":"modbot_polls","field":"chat_id","related_collection":"modbot_settings","meta":{"sort_field":null},"schema":{"on_delete":"SET NULL"}}' \
    $DIRECTUS_URL/relations \
