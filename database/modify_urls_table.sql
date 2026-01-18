-- Add user_id to URLs table for user data separation
-- Run this after creating the users table

-- Add user_id column to urls table
ALTER TABLE urls 
ADD COLUMN user_id UUID REFERENCES users(id) ON DELETE CASCADE;

-- Create index for user_id
CREATE INDEX idx_urls_user_id ON urls(user_id);

-- Update RLS policies for URLs table
-- Drop existing policies first
DROP POLICY IF EXISTS "Anyone can view URLs" ON urls;
DROP POLICY IF EXISTS "Anyone can create URLs" ON urls;
DROP POLICY IF EXISTS "Anyone can update hit count" ON urls;

-- Create new user-specific policies
CREATE POLICY "Users can view own URLs" ON urls
    FOR SELECT USING (auth.uid() = user_id);

CREATE POLICY "Users can create own URLs" ON urls
    FOR INSERT WITH CHECK (auth.uid() = user_id);

CREATE POLICY "Users can update own URLs" ON urls
    FOR UPDATE USING (auth.uid() = user_id);

CREATE POLICY "Users can delete own URLs" ON urls
    FOR DELETE USING (auth.uid() = user_id);

-- Public policy for URL redirects (read-only)
CREATE POLICY "Anyone can view URLs for redirect" ON urls
    FOR SELECT USING (true);
