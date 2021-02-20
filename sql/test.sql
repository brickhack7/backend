INSERT INTO profiles (uid, lat, long, interests) VALUES ('1', 42, 50, ARRAY['anime', 'hiking']);
INSERT INTO profiles (uid, lat, long, interests) VALUES ('2', 42.0000001, 50.0000001, ARRAY['biking', 'hiking']);
INSERT INTO profiles (uid, lat, long, interests) VALUES ('3', 41.9999999, 49.9999999, ARRAY['anime', 'biking']);


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