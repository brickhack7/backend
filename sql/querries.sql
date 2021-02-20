----select profile
SELECT * FROM profiles WHERE uid = $1;

----create profile
INSERT INTO profiles (uid, interests)
VALUES ($1, $2);

--update profiles
UPDATE profiles 
SET (interests, lat, long) = ($2, $3, $4) 
WHERE uid = $1;

--select locations
SELECT 
loc_id,
city, 
name,
(
   6371 *
   acos(cos(radians($1)) * 
   cos(radians(lat)) * 
   cos(radians(long) - 
   radians($2)) + 
   sin(radians($1)) * 
   sin(radians(lat )))
) AS distance 
FROM locations 
HAVING distance < $3 
ORDER BY distance LIMIT 1;

--create locations
INSERT INTO locations (loc_id, name, city, lat, long)
VALUES ($1, $2, $3, $4, $5);

--create location match
INSERT INTO loc_matches (uid, loc_id, status)
VALUES ($1, $2, 1);

--create user match
INSERT INTO user_matches (uid1, uid2, status)
VALUES ($1, $2, 1);

--select user matches
SELECT * FROM user_matches WHERE uid1 = $1 or uid2 = $1


