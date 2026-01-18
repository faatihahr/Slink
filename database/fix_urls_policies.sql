-- Fix URLs table RLS policies to allow authenticated users to create URLs
-- This addresses the 42501 RLS violation error

-- Drop existing URL policies that might be causing conflicts
DROP POLICY IF EXISTS "Users can view own URLs" ON urls;
DROP POLICY IF EXISTS "Users can create own URLs" ON urls;
DROP POLICY IF EXISTS "Users can update own URLs" ON urls;
DROP POLICY IF EXISTS "Users can delete own URLs" ON urls;
DROP POLICY IF EXISTS "Anyone can view URLs for redirect" ON urls;

-- Create new comprehensive policies for URLs table

-- Allow authenticated users to create URLs (user_id will be set automatically)
CREATE POLICY "Authenticated users can create URLs" ON urls
    FOR INSERT WITH CHECK (auth.uid() IS NOT NULL);

-- Allow users to view their own URLs
CREATE POLICY "Users can view own URLs" ON urls
    FOR SELECT USING (auth.uid() = user_id);

-- Allow users to update their own URLs (for hit count, etc.)
CREATE POLICY "Users can update own URLs" ON urls
    FOR UPDATE USING (auth.uid() = user_id);

-- Allow users to delete their own URLs
CREATE POLICY "Users can delete own URLs" ON urls
    FOR DELETE USING (auth.uid() = user_id);

-- Public policy for URL redirects (anyone can access the original_url for redirect)
CREATE POLICY "Public URL access for redirects" ON urls
    FOR SELECT USING (true);

-- Verify the policies
SELECT 
    schemaname,
    tablename,
    policyname,
    permissive,
    roles,
    cmd,
    qual,
    with_check
FROM pg_policies 
WHERE tablename = 'urls'
ORDER BY policyname;
