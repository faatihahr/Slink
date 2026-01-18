-- Temporarily disable RLS to test registration
ALTER TABLE users DISABLE ROW LEVEL SECURITY;

-- Test if this fixes the issue
-- If registration works, we know it's purely a policy issue

-- Check RLS status
SELECT 
    schemaname,
    tablename,
    rowsecurity 
FROM pg_tables 
WHERE tablename = 'users';
