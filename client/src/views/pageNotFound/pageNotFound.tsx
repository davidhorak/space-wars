import classNames from "classnames";
import { Trans, useTranslation } from "react-i18next";
import { Link } from "react-router-dom";

import { Path } from "../../router/path";

const PageNotFoundView = (): JSX.Element => {
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
      <h2 className={classNames("h2")}>{t("views.pageNotFound.title")}</h2>
      <h3 className={classNames("h3")}>{t("views.pageNotFound.subtitle")}</h3>
      <p>
        <Trans
          i18nKey="views.pageNotFound.body"
          components={{ 1: <Link to={Path.Battlefield} /> }}
        />
      </p>
    </div>
  );
};

export default PageNotFoundView;
