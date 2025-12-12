--META
{
    "id": "test-id-override-1",
    "name": "Test 1",
    "description": "Test description 1"
}
--end

--SQL:GetData1
SELECT '515bbf3c-93c5-476a-8dbc-4a6db4fe3c0c' AS id, 'Igor' AS name, 'en' AS language, 'igor@example.com' AS email, ARRAY['token1','token2'] AS tokens;
--end

--SQL:GetData2
-- Comment to be ignored
SELECT 'ef84af8f-bb55-4f74-9d7c-3db30e740d20' AS id, 'Alexey' AS name, 'en' AS language, 'alex@example.com' AS email, '{}'::varchar[] as tokens;
--end

--SQL:GetData3
SELECT 'e192f9e5-5e5c-4bba-b13e-0f9de32ec6bd' AS id, 'Denis' AS name, 'en' AS language, 'denis@example.com' AS email, ARRAY['token3','token4'] AS tokens;
--end
