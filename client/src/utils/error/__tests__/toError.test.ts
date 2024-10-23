import { toError } from '../toError';

describe('utils / error / toError', () => {
  it('returns an Error', () => {
    const error = new Error('test');
    expect(toError(error)).toBe(error);
  });

  it('returns an Error with the message', () => {
    const error = new Error('test');
    expect(toError(error)).toBe(error);
  });
});
