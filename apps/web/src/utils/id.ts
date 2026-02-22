// ID utility functions
const ULID_CHARS = '0123456789ABCDEFGHJKMNPQRSTVWXYZ'

export function generateULID(): string {
  const now = Date.now()
  let result = ''

  // Encode 48-bit timestamp (10 chars in base32)
  let ts = now
  for (let i = 9; i >= 0; i--) {
    result = ULID_CHARS[ts % 32] + result
    ts = Math.floor(ts / 32)
  }

  // Encode 80-bit randomness (16 chars in base32)
  for (let i = 0; i < 16; i++) {
    result += ULID_CHARS[Math.floor(Math.random() * 32)]
  }

  return result
}

export function isValidULID(id: string): boolean {
  return typeof id === 'string' && id.length === 26 && /^[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}$/i.test(id)
}
