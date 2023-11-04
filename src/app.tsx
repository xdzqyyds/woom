import { useAtom } from 'jotai'
import Welcome from './pages/welcome'
import Meeting from './pages/meeting'
import { meetingIdAtom } from './store/atom'
import backgroundImage from './assets/background.jpg'

export default function WOOM() {
  const [meetingId] = useAtom(meetingIdAtom)

  return (
    <div style={{
      height: '100vh',
      backgroundImage: `url(${backgroundImage})`,
    }}>
      {
        !meetingId
          ? <Welcome />
          : <Meeting meetingId={meetingId} />
      }
    </div>
  )
}
