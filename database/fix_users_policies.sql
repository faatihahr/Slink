-- Fix missing INSERT policy for users table
-- This allows new users to register

-- Drop existing policies first
DROP POLICY IF EXISTS "Users can view own profile" ON users;
DROP POLICY IF EXISTS "Users can update own profile" ON users;

-- Create comprehensive policies for users table
CREATE POLICY "Users can view own profile" ON users
    FOR SELECT USING (auth.uid() = id);

CREATE POLICY "Users can update own profile" ON users
    FOR UPDATE USING (auth.uid() = id);

CREATE POLICY "Users can insert own profile" ON users
    FOR INSERT WITH CHECK (auth.uid() = id);

-- Allow public registration (anyone can insert, but auth.uid() will be null for new users)
-- We need to handle this in the application by setting auth.uid() after user creation
CREATE POLICY "Enable insert for authentication" ON users
    FOR INSERT WITH CHECK (true);

-- Alternative approach: Disable RLS for inserts temporarily
-- ALTER TABLE users DISABLE ROW LEVEL SECURITY;

-- Or create a bypass for registration
CREATE POLICY "Allow registration" ON users
    FOR INSERT WITH CHECK (true);
