import { useEffect, useState } from 'react'
import { useAtom } from 'jotai'
import { getUserOnlineStatus, updateUserStatus } from '../../lib/api'
import { getStorage } from '../../lib/storage'
import { meetingIdAtom } from '../../store/atom'
import Invite from './Invite'

export default function UserList() {
  const [userStatus, setUserStatus] = useState<{ [userId: string]: string }>({})
  const [isOpen, setIsOpen] = useState<boolean>(false)

  const [meeting] = useAtom(meetingIdAtom)

  const inviterId = getStorage()?.userId
  const meetingId = getStorage()?.meeting

  const fetchUserStatus = async () => {
    try {
      const status = await getUserOnlineStatus()
      setUserStatus(status)
    } catch (error) {
      console.error('Failed to fetch user status:', error)
    }
  }

  useEffect(() => {
    fetchUserStatus()
    const interval = setInterval(fetchUserStatus, 5000)
    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    const cleanup = async () => {
      updateUserStatus(inviterId, '0')
    }

    window.addEventListener('beforeunload', cleanup)
    window.addEventListener('unload', cleanup)
    return () => {
      window.removeEventListener('beforeunload', cleanup)
      window.removeEventListener('unload', cleanup)
    }
  }, [])

  const sortedUserStatus = Object.keys(userStatus)
    .sort((_, b) => (userStatus[b] === '1' ? 1 : -1))
    .map((userId) => ({
      userId,
      status: userStatus[userId],
    }))

  return (
    <div className="relative">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="absolute bottom-4 right-4 bg-blue-500 text-white p-3 rounded-full shadow-lg"
      >
        {isOpen ? (
          <span className="text-xl">▲</span>
        ) : (
          <span className="text-xl">▼</span>
        )}
      </button>

      {isOpen && (
        <div className="absolute bottom-16 right-4 bg-white p-4 rounded-lg shadow-lg max-w-xs w-full max-h-[200px] overflow-y-auto">
          <h3 className="font-bold mb-2">User Online Status</h3>
          <ul className="space-y-2">
            {sortedUserStatus.map(({ userId, status }) => (
              <li key={userId} className="flex items-center justify-between space-x-1">
                <span className="font-bold text-lg text-blue-600">{userId}</span>

                <div className="flex items-center space-x-2">
                  {status === '1' ? (
                    <span className="text-green-500">✔️</span>
                  ) : (
                    <span className="text-red-500">❌</span>
                  )}
                  <span className={status === 'true' ? 'text-green-500' : 'text-red-500'}>
                    {status === '1' ? 'Online' : 'Offline'}
                  </span>
                </div>

                {status === '1' && meeting ? (
                  <Invite
                    meetingId={meetingId}
                    inviterId={inviterId}
                    inviteeId={userId}
                  />
                ) : (
                  <span className="text-gray-500">Disabled</span>
                )}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  )
}
