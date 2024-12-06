import { useAtom } from 'jotai'
import Welcome from './pages/welcome'
import Meeting from './pages/meeting'
import { meetingIdAtom, isLoggedInAtom } from './store/atom'
import UserList from './components/login/userlist'

export default function WOOM() {
  const [meetingId] = useAtom(meetingIdAtom)
  const [isLoggedIn] = useAtom(isLoggedInAtom)

  return (
    <div
      className="min-h-screen">
      {
        !meetingId
          ? <Welcome />
          : <Meeting meetingId={meetingId} />
      }
      {isLoggedIn && <UserList />}
    </div>
  )
}
