export type DragAndDropProps = {
  checked: boolean;
  className?: string;
  filter?: File["type"][];
  multiple?: boolean;
  disabled?: boolean;
  processing?: boolean;
  i18n?: I18n;
  onChange?: (files: File[]) => void;
  onError?: (error: Error) => void;
};

export type I18n = {
  drop: string;
  processing: string;
  select: {
    label?: string;
    title?: string;
  };
  description?: string;
};

export class InvalidFileTypeError extends Error {
  filename: File["name"];

  type: File["type"];

  constructor(filename: File["name"], type: File["type"]) {
    super(`Invalid file type: ${filename} (${type})`);
    this.filename = filename;
    this.type = type;
  }
}
