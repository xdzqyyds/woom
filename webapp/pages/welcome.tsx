import Join from '../components/join'
import Login from '../components/login/login'
import { useAtom } from 'jotai'
import { isLoggedInAtom } from '../store/atom'
// import { getStorage } from '../lib/storage'

export default function Welcome() {
  const [isLoggedIn] = useAtom(isLoggedInAtom)

  // const inviteeId = getStorage()?.userId
  return (
    <div className="flex flex-col justify-around min-h-screen">
      <center>
        <div>
          <h1 className="flex justify-center text-white text-5xl font-bold">WOOM</h1>
          <p className="text-white">A new meeting</p>
        </div >
        {
          !isLoggedIn
            ? <Login />
            : (
              <>
                <Join />
              </>
            )
        }
      </center>
    </div>
  )
}
