import { preventStop } from '../preventStop';

describe('utils / events / preventStop', () => {
  const createMockEvent = () => ({
    preventDefault: jest.fn(),
    stopPropagation: jest.fn(),
  });

  it('returns a function', () => {
    const handler = jest.fn();
    const wrappedHandler = preventStop(handler);
    expect(typeof wrappedHandler).toBe('function');
  });

  it('calls preventDefault and stopPropagation on the event', () => {
    const handler = jest.fn();
    const wrappedHandler = preventStop(handler);
    const mockEvent = createMockEvent();

    wrappedHandler(mockEvent);
    expect(mockEvent.preventDefault).toHaveBeenCalledTimes(1);
    expect(mockEvent.stopPropagation).toHaveBeenCalledTimes(1);
  });

  it('calls the provided handler with the event', () => {
    const handler = jest.fn();
    const wrappedHandler = preventStop(handler);
    const mockEvent = createMockEvent();

    wrappedHandler(mockEvent);
    expect(handler).toHaveBeenCalledTimes(1);
    expect(handler).toHaveBeenCalledWith(mockEvent);
  });
});
