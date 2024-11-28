import { useState } from 'react'
import { useAtom } from 'jotai'
import { userIdAtom, userPasswordAtom, isLoggedInAtom } from '../store/atom'
import Join from '../components/join'
import { login } from '../lib/api'

export default function Login() {
  const [userId, setUserId] = useAtom(userIdAtom)
  const [password, setPassword] = useAtom(userPasswordAtom)
  const [isLoggedIn, setIsLoggedIn] = useAtom(isLoggedInAtom)
  const [error, setError] = useState<string | null>(null)

  const handleLogin = async () => {
    if (userId === '' || password === '') {
      setError('User ID and password cannot be empty')
      return
    }
    try {
      const response = await login(userId, password)
      if (response.success) {
        setError(null)
        setIsLoggedIn(true)
      } else {
        setError(response.message || 'Login failed')
      }
    } catch {
      setError('Network error or server unavailable')
    }
  }

  if (isLoggedIn) {
    return <Join />
  }

  return (
    <div className="flex flex-col justify-around bg-gray-800/80 p-6 my-4 rounded-lg">
      <center className="flex flex-col items-center space-y-4">
        <input
          className="text-center text-lg border border-gray-300 rounded p-2"
          placeholder="Enter User ID"
          value={userId}
          onChange={(e) => setUserId(e.target.value)}
        />

        <input
          type="password"
          className="text-center text-lg border border-gray-300 rounded p-2"
          placeholder="Enter Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />

        {error && (
          <p className="text-red-500 text-sm">
            {error}
          </p>
        )}

        <button
          className="btn-primary"
          onClick={handleLogin}
        >
          Sign In
        </button>
      </center>
    </div>
  )
}
