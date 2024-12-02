import { useState } from 'react'
import { sendInvite } from '../lib/api'

interface InviteProps {
  meetingId: string;
  inviterId: string;
  inviteeId: string;
}

export default function Invite({ meetingId, inviterId, inviteeId }: InviteProps) {
  const [isInvited, setIsInvited] = useState<boolean>(false)

  const handleInvite = async () => {
    try {
      const response = await sendInvite(meetingId, inviterId, inviteeId)
      if (response.success) {
        setIsInvited(true)
        console.log('Invite sent successfully')
      } else {
        console.error('Failed to send invite')
      }
    } catch (error) {
      console.error('Error sending invite:', error)
    }
  }

  return (
    <div>
      <button
        onClick={handleInvite}
        className={`bg-blue-500 text-white p-2 rounded-md ${isInvited ? 'bg-green-500' : ''}`}
      >
        {isInvited ? 'Invited' : 'Invite'}
      </button>
    </div>
  )
}
