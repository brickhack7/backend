INSERT INTO profiles (uid, lat, long, interests) VALUES ('1', 42, 50, ARRAY['anime', 'hiking']);
INSERT INTO profiles (uid, lat, long, interests) VALUES ('2', 42.0000001, 50.0000001, ARRAY['biking', 'hiking']);
INSERT INTO profiles (uid, lat, long, interests) VALUES ('3', 41.9999999, 49.9999999, ARRAY['anime', 'biking']);
INSERT INTO profiles (uid, lat, long, interests) VALUES ('4', 41.9999, 49.999, ARRAY['running', 'hiking']);
INSERT INTO profiles (uid, lat, long, interests) VALUES ('5', 41.9999, 49.999, ARRAY['anime', 'hiking']);


INSERT INTO locations (loc_id, name, city, lat, long)
VALUES ('1', 'park', 'Toronto', '42', '50');
INSERT INTO locations (loc_id, name, city, lat, long)
VALUES ('2', 'hike', 'Toronto', '42', '50');
INSERT INTO locations (loc_id, name, city, lat, long)
VALUES ('3', 'resturant', 'Toronto', '42', '50');


INSERT INTO loc_matches (uid, loc_id, status)
VALUES ('1', '1', 1);
INSERT INTO loc_matches (uid, loc_id, status)
VALUES ('1', '2', 1);
INSERT INTO loc_matches (uid, loc_id, status)
VALUES ('1', '3', 1);

INSERT INTO loc_matches (uid, loc_id, status)
VALUES ('2', '1', 1);
INSERT INTO loc_matches (uid, loc_id, status)
VALUES ('2', '2', 1);

INSERT INTO loc_matches (uid, loc_id, status)
VALUES ('3', '1', 1);
INSERT INTO loc_matches (uid, loc_id, status)
VALUES ('3', '3', 1);

INSERT INTO loc_matches (uid, loc_id, status)
VALUES ('4', '2', 1);
INSERT INTO loc_matches (uid, loc_id, status)
VALUES ('4', '3', 1);

INSERT INTO loc_matches (uid, loc_id, status)
VALUES ('5', '1', 1);
INSERT INTO loc_matches (uid, loc_id, status)
VALUES ('5', '2', 1);
INSERT INTO loc_matches (uid, loc_id, status)
VALUES ('5', '3', 1);

--closest locations
WITH cte AS (
SELECT 
loc_id,
city, 
name,
(
   6371 *
   acos(cos(radians(42)) * 
   cos(radians(lat)) * 
   cos(radians(long) - 
   radians(50)) + 
   sin(radians(42)) * 
   sin(radians(lat )))
) AS distance 
FROM locations )
SELECT *
FROM cte
WHERE cte.distance < 50
ORDER BY cte.distance LIMIT 20;

--return the user which you share the most interests with, are not matched with and that has said yes to that location
with cte as (
    SELECT uid, count(1) as common
      FROM (
          SELECT uid, unnest(profiles.interests) as intr
            FROM profiles WHERE profiles.uid = any(SELECT uid from loc_matches where loc_id = '1' and uid <> '1' and status = '1')
            and profiles.uid <> ALL(SELECT uid1 from user_matches where uid2 = '1')
            and profiles.uid <> ALL(SELECT uid2 from user_matches where uid1 = '1')
            )
x WHERE intr = any(
    SELECT unnest(interests) FROM profiles WHERE uid = '1') AND uid <> '1'
GROUP BY uid) 
SELECT cte.uid FROM cte INNER JOIN profiles p on cte.uid = p.uid ORDER BY common DESC LIMIT 1;

INSERT INTO user_matches (uid1, uid2, status)
VALUES ('1', '5', 1);

SELECT * FROM user_matches WHERE uid1 = '1' or uid2 = '1';