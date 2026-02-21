import { describe, it, expect } from 'vitest';
import { generateULID, isValidULID } from './id';

describe('generateULID', () => {
  it('should generate a valid ULID', () => {
    const id = generateULID();
    expect(id).toHaveLength(26);
    expect(typeof id).toBe('string');
  });

  it('should generate unique ULIDs', () => {
    const ids = new Set();
    for (let i = 0; i < 100; i++) {
      ids.add(generateULID());
    }
    expect(ids.size).toBe(100);
  });
});

describe('isValidULID', () => {
  it('should return true for valid ULID', () => {
    const id = generateULID();
    expect(isValidULID(id)).toBe(true);
  });

  it('should return false for invalid ULID', () => {
    expect(isValidULID('')).toBe(false);
    expect(isValidULID('invalid')).toBe(false);
    expect(isValidULID('123')).toBe(false);
  });
});
