-- Apply this SQL in Supabase SQL Editor to fix user registration
-- This adds the missing INSERT policy for the users table

-- Add INSERT policy to allow user registration
CREATE POLICY "Allow user registration" ON users
    FOR INSERT WITH CHECK (true);

-- Verify all policies are in place
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
WHERE tablename = 'users';
