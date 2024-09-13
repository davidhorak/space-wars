import classNames from "classnames";
import { ButtonProps } from "./types";
import style from "./button.module.css";
import { PropsWithChildren } from "react";
import { noop } from "lodash/fp";

const Button = ({
  onClick = noop,
  children,
  className,
  disabled = false,
}: PropsWithChildren<ButtonProps>) => {
  return (
    <button
      className={classNames(
        style.button,
        className,
        "border-0",
        "border-radius-24",
        "cursor-pointer",
        "d-inline-block",
        "text-decoration-none",
        "text-center",
        "py-12",
        "px-24",
        {
          [style['button--disabled']]: disabled,
        }
      )}
      onClick={onClick}
      disabled={disabled}
    >
      {children}
    </button>
  );
};

export default Button;
