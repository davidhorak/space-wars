import classNames from "classnames";
import { Size, type SpinnerProps } from "./types";
import style from "./spinner.module.css";

const Spinner = ({ size = Size.Medium }: SpinnerProps): JSX.Element => {
  return (
    <div className={classNames(style.spinner, style[`spinner-${size.toLocaleLowerCase()}`])} />
  );
};

export default Spinner;
