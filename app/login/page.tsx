'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/lib/auth-context';
import { isSupabaseConfigured } from '@/lib/supabase';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

export default function LoginPage() {
  const router = useRouter();
  const { signIn, signUp } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [mode, setMode] = useState<'login' | 'signup'>('login');

  if (!isSupabaseConfigured) {
    return (
      <div className="flex min-h-[60vh] items-center justify-center">
        <Card title="Dev Mode" shadow="raised" className="w-full max-w-sm text-center">
          <p className="text-body mb-4">Auth is not configured. Running in dev mode.</p>
          <Button variant="primary" onClick={() => router.push('/')}>
            Continue as Demo User
          </Button>
        </Card>
      </div>
    );
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    const result = mode === 'login' ? await signIn(email, password) : await signUp(email, password);

    setLoading(false);
    if (result.error) {
      setError(result.error);
    } else {
      router.push('/');
    }
  };

  return (
    <div className="flex min-h-[60vh] items-center justify-center">
      <Card title={mode === 'login' ? 'Sign In' : 'Create Account'} shadow="raised" className="w-full max-w-sm">
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="Email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="you@example.com"
            required
          />
          <Input
            label="Password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="••••••••"
            required
            minLength={6}
          />

          {error && <p className="text-sm text-red-600">{error}</p>}

          <Button variant="primary" type="submit" loading={loading} className="w-full">
            {mode === 'login' ? 'Sign In' : 'Sign Up'}
          </Button>
        </form>

        <p className="mt-4 text-center text-sm text-body">
          {mode === 'login' ? "Don't have an account? " : 'Already have an account? '}
          <button
            type="button"
            onClick={() => { setMode(mode === 'login' ? 'signup' : 'login'); setError(''); }}
            className="font-medium text-accent hover:underline"
          >
            {mode === 'login' ? 'Sign Up' : 'Sign In'}
          </button>
        </p>
      </Card>
    </div>
  );
}
