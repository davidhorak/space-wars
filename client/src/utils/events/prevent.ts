import { noop } from 'lodash';

export const prevent =
  <T extends React.UIEvent | React.MouseEvent>(handler: (event: T) => void = noop) =>
  (event: T) => {
    event.preventDefault();
    handler(event);
  };
