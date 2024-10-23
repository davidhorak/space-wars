export const preventStop =
  <T extends React.UIEvent>(handler: (event: T) => void) =>
  (event: T) => {
    event.preventDefault();
    event.stopPropagation();
    handler(event);
  };
