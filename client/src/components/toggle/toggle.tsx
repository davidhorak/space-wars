import classNames from "classnames";
import { ToggleProps } from "./types";
import style from "./toggle.module.css";

const Toggle = ({ checked, className, disabled = false, onChange }: ToggleProps) => {
  return (
    <label
      className={classNames(
        style.toggle,
        className,
        "d-inline-block",
        "position-relative",
        { [style["toggle--disabled"]]: disabled }
      )}
    >
      <input
        type="checkbox"
        checked={checked}
        onChange={() => onChange(!checked)}
        disabled={disabled}
      />
      <span
        className={classNames(
          style.toggle__slider,
          "position-absolute",
          "top-0",
          "left-0",
          "right-0",
          "bottom-0"
        )}
      />
    </label>
  );
};

export default Toggle;
