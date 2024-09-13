import classNames from "classnames";
import { Button } from "../button";
import { useTranslation } from "react-i18next";

const GlobalError = ({ error }: { error?: Error }): JSX.Element => {
  const { t } = useTranslation();

  return (
    <div
      className={classNames(
        "d-flex",
        "flex-column",
        "align-items-center",
        "justify-content-center",
        "h-100"
      )}
    >
      <h2 className={classNames("h2")}>{t("error.title")}</h2>
      <p className={classNames("mt-6")}>
        {error?.message ?? t("error.message")}
      </p>
      <Button
        className={classNames("mt-24")}
        onClick={() => window.location.reload()}
      >
        {t("error.refresh.button")}
      </Button>
    </div>
  );
};

export default GlobalError;
