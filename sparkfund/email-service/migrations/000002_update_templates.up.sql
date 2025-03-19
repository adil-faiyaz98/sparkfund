-- Drop existing columns
ALTER TABLE templates DROP COLUMN content;

-- Add new columns
ALTER TABLE templates ADD COLUMN body TEXT NOT NULL;
ALTER TABLE templates ADD COLUMN variables TEXT[] NOT NULL;
ALTER TABLE templates ADD COLUMN description TEXT; 