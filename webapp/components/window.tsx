import { useState, useEffect } from 'react'
import { getInvitation } from '../lib/api'
import { setStorageMeeting } from '../lib/storage'
import { useAtom } from 'jotai'
import { locationAtom, meetingIdAtom } from '../store/atom'

interface InviteWindowProps {
  inviteeId: string
}

export default function InviteWindow({ inviteeId }: InviteWindowProps) {
  const [invitation, setInvitation] = useState<string | null>(null)
  const [isOpen, setIsOpen] = useState(false)
  const [_, setLoc] = useAtom(locationAtom)
  const [__, setAtomMeetingId] = useAtom(meetingIdAtom)

  type invitation = { value: string }

  const checkInvitation = async () => {
    try {
      const value = await getInvitation(inviteeId)
      console.log('checkInvitation', value)
      if (value) {
        setInvitation(value)
        setIsOpen(true)
      }
    } catch { /* empty */ }
  }

  useEffect(() => {
    checkInvitation()
    const interval = setInterval(checkInvitation, 5000)
    return () => clearInterval(interval)
  }, [inviteeId])

  const handleAccept = () => {
    console.log('Accepted the invitation')
    const invitationValue = invitation?.value
    const roomId = invitationValue.split(' ')[0]
    setStorageMeeting(roomId)
    setAtomMeetingId(roomId)
    setLoc(prev => ({ ...prev, pathname: `/${roomId}` }))
    setIsOpen(false)
  }

  const handleReject = () => {
    console.log('Rejected the invitation')
    setIsOpen(false)
  }

  return (
    <>
      {isOpen && invitation && (
        <div className="fixed top-4 right-4 bg-white p-4 rounded-lg shadow-lg max-w-xs w-64">
          <h3 className="font-bold mb-2">You have an invitation!</h3>
          <p className="mb-4">
            {invitation.value}
          </p>
          <div className="flex justify-between space-x-4">
            <button
              onClick={handleAccept}
              className="bg-green-500 text-white p-2 rounded-md"
            >
              Accept
            </button>
            <button
              onClick={handleReject}
              className="bg-red-500 text-white p-2 rounded-md"
            >
              Reject
            </button>
          </div>
        </div>
      )}
    </>
  )
}