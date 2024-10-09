import { noop } from "lodash/fp";
import { PropsWithChildren } from "react";
import { Button } from "../button";
import { ModalProps } from "./types";
import classNames from "classnames";

import style from "./modal.module.css";

import IconClose from "../../icons/close.svg?react";
import { useTranslation } from "react-i18next";

const Modal = ({
  className,
  title,
  children,
  onClose = noop,
}: PropsWithChildren<ModalProps>) => {
  const { t } = useTranslation();

  return (
    <div
      className={classNames(
        style.modal,
        className,
        "position-fixed",
        "d-flex",
        "justify-content-center",
        "align-items-center",
        "w-100",
        "h-100",
        "top-0",
        "left-0",
        "right-0",
        "bottom-0"
      )}
      onClick={onClose}
    >
      <div className={classNames(style["modal__content"], "p-24", "border-radius-12")}>
        <div className={classNames(style["modal__header"], "position-relative")}>
          <h2 className={classNames("px-24", "text-center", "h5")}>{title}</h2>
          <button className={classNames("button-icon", "position-absolute", "top-0", "right-0")} onClick={onClose}>
            <IconClose width={24} height={24} />
          </button>
        </div>
        <div className={classNames(style["modal__body"], "mt-24")}>{children}</div>
        <div className={classNames(style["modal__footer"], "mt-24", "text-center")}>
          <Button onClick={onClose}>{t("modal.ok")}</Button>
        </div>
      </div>
    </div>
  );
};

export default Modal;
