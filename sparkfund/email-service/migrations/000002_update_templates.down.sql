-- Drop new columns
ALTER TABLE templates DROP COLUMN body;
ALTER TABLE templates DROP COLUMN variables;
ALTER TABLE templates DROP COLUMN description;

-- Add back original column
ALTER TABLE templates ADD COLUMN content TEXT NOT NULL; 