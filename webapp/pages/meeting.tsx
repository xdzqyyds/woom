import { useAtom } from 'jotai'
import { meetingJoinedAtom } from '../store/atom'
import Layout from '../components/layout'
import Prepare from '../components/prepare'
import InviteWindow from '../components/login/window'
import { getStorage } from '../lib/storage'

export default function Meeting(props: { meetingId: string }) {
  const [meetingJoined] = useAtom(meetingJoinedAtom)
  const inviteeId = getStorage()?.userId
  return (
    <div className="flex flex-col justify-around min-h-screen">
      {meetingJoined
        ? <Layout meetingId={props.meetingId} />
        : <Prepare meetingId={props.meetingId} />
      }
      <InviteWindow inviteeId={inviteeId} />
    </div>
  )
}
