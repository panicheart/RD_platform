// ID utility functions
export function generateULID(): string {
  // TODO: implement actual ULID generation
  return '01H' + Math.random().toString(36).substring(2, 15).padEnd(23, '0');
}

export function isValidULID(id: string): boolean {
  return typeof id === 'string' && id.length === 26 && /^[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}$/i.test(id);
}
