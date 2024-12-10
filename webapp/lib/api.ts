/* eslint-disable @typescript-eslint/no-explicit-any */

interface Room {
  roomId: string,
  locked: false,
  owner: string,
  presenter?: string,
  streams: any,
}

/**
 * @see https://developer.mozilla.org/en-US/docs/Web/API/RTCPeerConnection/connectionState#value
 */
enum StreamState {
  New = 'new',
  Signaled = 'signaled',
  Connecting = 'connecting',
  Connected = 'connected',
  Disconnected = 'disconnected',
  Failed = 'failed',
  Closed = 'closed',
}

interface Stream {
  name: string,
  state: StreamState
  audio: boolean,
  video: boolean,
  screen: boolean,
}

interface User {
  streamId: string,
  token: string,
}

let token = ''
let roomId = ''

function setApiToken(str: string) {
  token = str
}

function setRoomId(str: string) {
  roomId = str
}

async function newUser(): Promise<User> {
  return (await fetch('/user/', {
    method: 'POST',
  })).json()
}

async function newRoom(): Promise<Room> {
  return (await fetch('/room/', {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
    method: 'POST',
  })).json()
}

async function getRoom(roomId: string): Promise<Room> {
  return (await fetch(`/room/${roomId}`, {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  })).json()
}

async function setRoom(roomId: string, data: any): Promise<Room> {
  return (await fetch(`/room/${roomId}`, {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    method: 'PATCH',
    body: JSON.stringify(data),
  })).json()
}

async function delRoom(roomId: string): Promise<void> {
  return (await fetch(`/room/${roomId}`, {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
    method: 'DELETE',
  })).json()
}

async function newStream(roomId: string): Promise<Stream> {
  return (await fetch(`/room/${roomId}/stream`, {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
    method: 'POST',
  })).json()
}

async function setStream(streamId: string, data: any): Promise<Stream> {
  return (await fetch(`/room/${roomId}/stream/${streamId}`, {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    method: 'PATCH',
    body: JSON.stringify(data),
  })).json()
}

async function delStream(roomId: string, streamId: string): Promise<any> {
  return fetch(`/room/${roomId}/stream/${streamId}`, {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
    method: 'DELETE',
    keepalive: true,
  })
}

async function login(userId: string, password: string): Promise<{ success: boolean; message: string }> {
  return (await fetch('/login/', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ userId, password }),
  })).json()
}

async function getUserOnlineStatus(): Promise<{ [userId: string]: string }> {
  return (await fetch('/login/userlist', {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  })).json()
}

async function updateUserStatus(userId: string, status: string): Promise<void> {
  await fetch('/login/offline', {
    method: 'PATCH',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      userId,
      status,
    }),
  })
}

function sendInvite(meetingId: string, inviterId: string, inviteeId: string): Promise<void> {
  return fetch('/login/invite', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      meetingId,
      inviterId,
      inviteeId,
    }),
  })
    .then(() => {
    })
    .catch((error) => {
      console.error('Error sending invite:', error)
    })
}


interface InvitationResponse {
  value: string;
}

async function getInvitation(inviteeId: string): Promise<InvitationResponse | null> {
  try {
    const response = await fetch('/login/invitee', {
      method: 'PATCH',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ inviteeId }),
    })

    const result = await response.json()

    if (!result.value) {
      return null
    }

    return result
  } catch (error) {
    console.error('Failed to check invitation:', error)
    return null
  }
}

async function signup(userId: string, password: string): Promise<{ success: boolean; message: string }> {
  const response = await fetch('/login/signup', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ userId, password }),
  })
  return response.json()
}

export {
  setRoomId,
  setApiToken,
  newUser,

  newRoom,
  getRoom,
  setRoom,
  delRoom,

  newStream,
  setStream,
  delStream,

  login,
  getUserOnlineStatus,
  updateUserStatus,
  sendInvite,
  getInvitation,
  signup,

  StreamState,
}

export type {
  Stream,
}
