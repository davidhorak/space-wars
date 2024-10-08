import classNames from "classnames";

import style from "./dragAndDrop.module.css";
import { DragAndDropProps, InvalidFileTypeError } from "./types";
import { noop, partition } from "lodash/fp";
import { PropsWithChildren, useMemo, useRef, useState } from "react";
import { preventStop } from "../../utils/events";
import { Spinner } from "../spinner";
import { Button } from "../button";
import { DefaultI18n } from "./i18n";

import IconFileDrop from "../../icons/file-drop.svg?react";
import IconFile from "../../icons/file.svg?react";
import IconFolder  from "../../icons/folder.svg?react";

const DragAndDrop = ({
  className,
  i18n = DefaultI18n,
  filter,
  children,
  multiple = false,
  processing = false,
  disabled = false,
  onChange = noop,
  onError = noop,
}: PropsWithChildren<DragAndDropProps>) => {
  const [hasHighlight, setHasHighlight] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const icon = useMemo(
    () =>
      !disabled && hasHighlight ? (
        <IconFileDrop width={24} height={24} />
      ) : (
        <IconFile width={24} height={24} />
      ),
    [disabled, hasHighlight]
  );

  const onOpen = () => {
    if (!disabled) {
      inputRef.current?.click();
    }
  };

  const onFilesChanged = (fileList: FileList | null) => {
    if (disabled || processing || !fileList) {
      return;
    }

    const [validFiles, invalidFiles] = partition(
      (file: File) => !filter || !filter.length || filter.includes(file.type)
    )(Array.from(fileList));

    if (invalidFiles.length) {
      invalidFiles.forEach((file: File) =>
        onError(new InvalidFileTypeError(file.name, file.type))
      );
      return;
    }

    if (validFiles.length) {
      onChange(validFiles);
    }
  };

  const onDrop = (event: React.DragEvent) => {
    setHasHighlight(false);

    if (disabled) return;

    onFilesChanged(event.dataTransfer?.files ?? null);
  };

  return (
    <div
      className={classNames(
        "drag-and-drop",
        style["drag-and-drop"],
        { [style["drag-and-drop--highlighted"]]: hasHighlight },
        { [style["drag-and-drop--disabled"]]: disabled },
        className
      )}
    >
      {processing ? (
        <div
          className={classNames(
            style["drag-and-drop__processing"],
            "w-100 d-block rounded-md p-24 text-center"
          )}
        >
          <Spinner />
          <p className="text-uppercase pt-6 b6">{i18n.processing}</p>
        </div>
      ) : (
        <button
          className={classNames(
            style["drag-and-drop__upload"],
            "drag-and-drop__upload w-100 d-block rounded-md p-24 position-relative"
          )}
          type="button"
          onDragEnter={preventStop(() => setHasHighlight(true))}
          onDragOver={preventStop(() => setHasHighlight(true))}
          onDragLeave={preventStop(() => setHasHighlight(false))}
          onDrop={preventStop(onDrop)}
          onClick={preventStop(onOpen)}
        >
          <div
            className={classNames(
              style["drag-and-drop__preview"],
              "drag-and-drop__preview"
            )}
          >
            {children}
          </div>
          <div
            className={classNames(
              style["drag-and-drop__legend"],
              {
                [style["drag-and-drop__legend--overlay"]]: !!children,
                "rounded-md": !!children,
              },
              "drag-and-drop__legend d-flex align-items-center justify-content-center"
            )}
          >
            <div>
              {icon}
              <p className="text-uppercase pt-6 b6">{i18n.drop}</p>
            </div>
          </div>
        </button>
      )}
      <Button
        className={classNames(
          style["drag-and-drop__select"],
          "d-block mt-24 w-100 rounded-md"
        )}
        title={i18n.select.title}
        disabled={disabled || processing}
        onClick={onOpen}
      >
        <div className="d-flex justify-content-center">
          <IconFolder
            width={24}
            height={24}
            className="mr-6 align-self-center"
          />
          <span>{i18n.select.label}</span>
        </div>
      </Button>
      {!!i18n.description && (
        <div
          className={classNames(
            style["drag-and-drop__desc"],
            "drag-and-drop__desc d-flex mt-12"
          )}
        >
          <div
            className={classNames(
              style["drag-and-drop__desc__icon"],
              "rounded-circle p-4 mr-6"
            )}
          >
            <IconFile />
          </div>
          <p className="self-align-center">{i18n.description}</p>
        </div>
      )}
      <input
        className="d-none"
        ref={inputRef}
        type="file"
        multiple={multiple}
        accept={filter?.join(", ")}
        onChange={(event) => onFilesChanged(event.target.files)}
      />
    </div>
  );
};

export default DragAndDrop;
