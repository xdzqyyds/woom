import { useState } from 'react'
import { useAtom } from 'jotai'
import {
  locationAtom,
  meetingIdAtom,
} from '../store/atom'
import { getStorage, setStorage, delStorage, setStorageStream, setStorageMeeting } from '../lib/storage'
import { newUser, setApiToken, setRoomId } from '../lib/api'

export default function Join() {
  const [, setLoc] = useAtom(locationAtom)
  const [__, setAtomMeetingId] = useAtom(meetingIdAtom)
  //const [tmpId, setTmpId] = useState<string>('')
  const [selectedRoom, setSelectedRoom] = useState<string>('')

  const roomOptions = [
    { name: 'room 1', id: '10001' },
    { name: 'room 2', id: '10002' },
    { name: 'room 3', id: '10003' },
  ]

  const getLoginStatus = async () => {
    const user = getStorage()
    if (!user.token || !user.stream) {
      const res = await newUser()
      user.token = res.token
      user.stream = res.streamId
      setStorage(user)
    }

    setApiToken(user.token)
    if (user.stream) setStorageStream(user.stream)
  }

  const joinMeeting = async () => {
    await getLoginStatus()
    const meetingId = selectedRoom
    //await fetch(`/room/${meetingId}`, {
    //  method: "PATCH"
    //})
    enterMeeting(meetingId)
    setRoomId(meetingId)
  }

  const enterMeeting = (meetingId: string) => {
    setStorageMeeting(meetingId)
    setAtomMeetingId(meetingId)
    setLoc(prev => ({ ...prev, pathname: `/${meetingId}` }))
  }

  //useEffect(() => {
  //  const id = loc.pathname?.replace('/', '')
  //  if (id) {
  //    setTmpId(id)
  //  }
  //}, [location])

  return (
    <div className="flex flex-col justify-around bg-gray-800/80 p-6 my-4">
      <center className="flex flex-row flex-wrap justify-center">
        {/*<button className="btn-primary my-2" disabled={!!tmpId} onClick={() => { newMeeting() }}>New Meeting</button>*/}
        {/*<div className="mx-2 my-2">
          <input
            className="text-center text-4xl"
            placeholder="Enter Meeting id"
            value={tmpId}
            onChange={e => setTmpId(e.target.value)}
            maxLength={11}
          />
        </div>*/}
        <div className="mx-2 my-2">
          <select
            className="text-center font-semibold text-lg py-2 px-4 border rounded-md bg-white shadow-md hover:border-blue-500 focus:border-blue-500 focus:ring focus:ring-blue-200 focus:ring-opacity-50 transition duration-150 ease-in-out"
            value={selectedRoom}
            onChange={e => setSelectedRoom(e.target.value)}
            style={{ color: selectedRoom ? '#2d3748' : '#718096' }}
          >
            <option value="" disabled hidden style={{ color: '#718096' }}>
              Select a room
            </option>
            {roomOptions.map(room => (
              <option key={room.id} value={room.id} className="text-gray-800">
                {room.name}
              </option>
            ))}
          </select>
        </div>

        <button className="btn-primary my-2" disabled={!selectedRoom} onClick={() => { joinMeeting() }}>Join</button>
      </center>
      <center className="flex flex-row flex-wrap justify-center text-white">
        <p>If have some problems, Please click this:</p>
        <a className="mx-2 text-red-300 underline" onClick={delStorage}>Reset</a>
      </center>
    </div>
  )
}
