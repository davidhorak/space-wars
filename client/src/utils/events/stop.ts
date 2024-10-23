export const stop =
  <T extends React.UIEvent & React.MouseEvent>(handler: (event: T) => void) =>
  (event: T) => {
    event.stopPropagation();
    handler(event);
  };
