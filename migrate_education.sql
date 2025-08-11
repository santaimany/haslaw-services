-- Update education field from string to JSON array format

-- For simple education entries (single entry)
UPDATE members 
SET education = JSON_ARRAY(education) 
WHERE education IS NOT NULL 
  AND education != '' 
  AND education NOT LIKE '[%'
  AND education NOT LIKE '%,%';

-- For education entries with commas (multiple entries)
UPDATE members 
SET education = JSON_ARRAY(
    TRIM(SUBSTRING_INDEX(education, ',', 1)),
    TRIM(SUBSTRING_INDEX(SUBSTRING_INDEX(education, ',', 2), ',', -1))
)
WHERE education IS NOT NULL 
  AND education != '' 
  AND education NOT LIKE '[%'
  AND education LIKE '%,%'
  AND CHAR_LENGTH(education) - CHAR_LENGTH(REPLACE(education, ',', '')) = 1;

-- For education entries with multiple commas (3+ entries)
UPDATE members 
SET education = CONCAT(
    '["', 
    REPLACE(education, ', ', '","'),
    '"]'
)
WHERE education IS NOT NULL 
  AND education != '' 
  AND education NOT LIKE '[%'
  AND CHAR_LENGTH(education) - CHAR_LENGTH(REPLACE(education, ',', '')) > 1;
