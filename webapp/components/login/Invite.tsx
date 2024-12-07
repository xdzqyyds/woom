import { useState } from 'react'
import { sendInvite } from '../../lib/api'

interface InviteProps {
  meetingId: string;
  inviterId: string;
  inviteeId: string;
}

export default function Invite({ meetingId, inviterId, inviteeId }: InviteProps) {
  const [isInvited, setIsInvited] = useState<boolean>(false)

  const handleInvite = () => {
    sendInvite(meetingId, inviterId, inviteeId)
      .catch((error) => {
        console.error('Error sending invite:', error)
      })
    setIsInvited(true)
    console.log('Invite sent successfully')
  }

  return (
    <div>
      <button
        onClick={handleInvite}
        className={`bg-blue-500 text-white p-2.8 rounded-md ${isInvited ? 'bg-green-500' : ''}`}
      >
        {isInvited ? 'Invited' : 'Invite'}
      </button>
    </div>
  )
}
