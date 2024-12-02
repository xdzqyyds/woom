import Login from '../components/login'
import Join from '../components/join'
import { useAtom } from 'jotai'
import { isLoggedInAtom } from '../store/atom'

export default function Welcome() {
  const [isLoggedIn] = useAtom(isLoggedInAtom)

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
            : <Join />
        }
      </center>
    </div>
  )
}
