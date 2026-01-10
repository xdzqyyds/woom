const MeetingKey = 'meeting'
const StreamKey = 'stream'
const TokenKey = 'token'
const NameKey = 'name'
const UserId = 'userId'

interface Storage {
  meeting?: string,
  stream?: string,
  token?: string,
  name?: string,
  userId?: string,
}

function setStorageMeeting(value: string) { localStorage.setItem(MeetingKey, value) }
function setStorageStream(value: string) { localStorage.setItem(StreamKey, value) }
function setStorageToken(value: string) { localStorage.setItem(TokenKey, value) }
function setStorageName(value: string) { localStorage.setItem(NameKey, value) }
function setStorageUserId(value: string) { localStorage.setItem(UserId, value) }

function getStorageMeeting(): string { return localStorage.getItem(MeetingKey) || '' }
function getStorageStream(): string { return localStorage.getItem(StreamKey) || '' }
function getStorageToken(): string { return localStorage.getItem(TokenKey) || '' }
function getStorageName(): string { return localStorage.getItem(NameKey) || '' }
function getStorageUserId(): string { return localStorage.getItem(UserId) || '' }

function delStorageMeeting() { localStorage.removeItem(MeetingKey) }
function delStorageStream() { localStorage.removeItem(StreamKey) }
function delStorageToken() { localStorage.removeItem(TokenKey) }
function delStorageName() { localStorage.removeItem(NameKey) }
function delStorageUserId() { localStorage.removeItem(UserId) }

function setStorage(opt: Storage) {
  if (opt.meeting) setStorageMeeting(opt.meeting)
  if (opt.stream) setStorageStream(opt.stream)
  if (opt.token) setStorageToken(opt.token)
  if (opt.name) setStorageName(opt.name)
  if (opt.userId) setStorageUserId(opt.userId)
}

function getStorage(): Storage {
  return {
    meeting: getStorageMeeting(),
    stream: getStorageStream(),
    token: getStorageToken(),
    name: getStorageName(),
    userId: getStorageUserId(),
  } as Storage
}

function delStorage() {
  delStorageMeeting()
  delStorageStream()
  delStorageToken()
  delStorageName()
  delStorageUserId()
}

export {
  setStorageMeeting,
  setStorageStream,
  setStorageToken,
  setStorageName,
  setStorageUserId,

  getStorageMeeting,
  getStorageStream,
  getStorageToken,
  getStorageName,
  getStorageUserId,

  setStorage,
  getStorage,
  delStorage,
}
