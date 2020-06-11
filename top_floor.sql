SELECT * FROM rentInfo WHERE floor = max_floor AND is_new = 1
AND title NOT LIKE "%六張犁%" 
AND title NOT LIKE "%北醫%"
AND title NOT LIKE "%麟光%"
AND option_type NOT LIKE "%分租%" 