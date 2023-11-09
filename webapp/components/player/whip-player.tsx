import { useEffect, useRef, useState } from 'react'
import { useAtom } from 'jotai'
import {
  localStreamAtom,
  currentDeviceAudioAtom,
  currentDeviceVideoAtom,
} from '../../store/atom'
import Player from './player'
import WHIPClient from '../../lib/whip'

export default function WhipPlayer(props: { streamId: string }) {
  const refEnabled = useRef(false)
  const refPC = useRef<RTCPeerConnection | null>(null)
  const [localStream] = useAtom(localStreamAtom)
  const [connectionState, setConnectionState] = useState("unknown")

  const [currentDeviceAudio] = useAtom(currentDeviceAudioAtom)
  const [currentDeviceVideo] = useAtom(currentDeviceVideoAtom)

  const newPeerConnection = () => {
    const stream = localStream.stream
    if (stream) {
      const pc = new RTCPeerConnection()
      pc.onconnectionstatechange = () => setConnectionState(pc.connectionState)

      // NOTE: array audio index is: 0
      if (!stream.getAudioTracks().length) {
        pc.addTransceiver('audio', { 'direction': 'sendonly' })
      } else {
        stream.getAudioTracks().map(track => pc.addTrack(track))
      }

      // NOTE: array video index is: 1
      if (!stream.getVideoTracks().length) {
        pc.addTransceiver('video', { 'direction': 'sendonly' })
      } else {
        stream.getVideoTracks().map(track => pc.addTrack(track))
      }

      //pc.addTransceiver(stream.getVideoTracks()[0], {
      //  direction: 'sendonly',
      //  //sendEncodings: [
      //  //  { rid: 'a', scaleResolutionDownBy: 2.0 },
      //  //  { rid: 'b', scaleResolutionDownBy: 1.0, },
      //  //  { rid: 'c' }
      //  //]
      //})

      refPC.current = pc
    }
  }

  const start = async (resource: string) => {
    const stream = localStream.stream
    if (stream) {
      if (refPC.current) {
        const whip = new WHIPClient();
        const url = location.origin + `/whip/${resource}`
        const token = "xxx"
        await whip.publish(refPC.current, url, token);
      }
    }
  }

  const restart = async (resource: string) => {
    if (refPC.current) {
      refPC.current.close()
    }
    newPeerConnection()
    start(resource)
  }

  const init = () => {
    if (localStream.stream) {
      if (!refEnabled.current) {
        refEnabled.current = true
        newPeerConnection()
        start(props.streamId)
      }
    }
  }
  useEffect(() => {
    init()
  }, [])

  useEffect(() => {
    const mediaStream = localStream.stream
    // If WebRTC is connected, switch track
    // NOTE: array audio index is: 0
    refPC.current?.getSenders().filter((_, i) => i === 0).map(sender => {
      if (mediaStream) {
        mediaStream.getAudioTracks().map(track => sender.replaceTrack(track))
      }
    })
    init()
  }, [currentDeviceAudio])

  useEffect(() => {
    const mediaStream = localStream.stream
    // If WebRTC is connected, switch track
    // NOTE: array video index is: 1
    refPC.current?.getSenders().filter((_, i) => i === 1).map(sender => {
      if (mediaStream) {
        mediaStream.getVideoTracks().map(track => sender.replaceTrack(track))
      }
    })
    init()
  }, [currentDeviceVideo])

  return (
    <div className='flex flex-col'>
      <Player user={localStream} muted={true} />
      <center className='text-white flex flex-row justify-around'>
        <p className='rounded-xl p-2 b-1 hover:border-orange-300'>{connectionState}</p>
        <button className='btn-primary' disabled={connectionState === 'connected'} onClick={() => restart(props.streamId)}>restart</button>
      </center>
    </div>
  )
}
