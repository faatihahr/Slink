-- Remove conflicting INSERT policies and keep only one clear policy
DROP POLICY IF EXISTS "Users can insert own profile" ON users;
DROP POLICY IF EXISTS "Enable insert for authentication" ON users;
DROP POLICY IF EXISTS "Allow registration" ON users;
DROP POLICY IF EXISTS "Allow user registration" ON users;

-- Keep only the necessary policies
CREATE POLICY "Users can view own profile" ON users
    FOR SELECT USING (auth.uid() = id);

CREATE POLICY "Users can update own profile" ON users
    FOR UPDATE USING (auth.uid() = id);

-- Add single clear INSERT policy for registration
CREATE POLICY "Allow user registration" ON users
    FOR INSERT WITH CHECK (true);

-- Verify final policies
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
WHERE tablename = 'users'
ORDER BY policyname;
