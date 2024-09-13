import classNames from "classnames";
import { useTranslation } from "react-i18next";

const Footer = (): JSX.Element => {
  const { t } = useTranslation();
  return (
    <footer className={classNames("mx-24", "py-24")}>
      <p className={classNames("text-center")}>
        {t("copyright", { year: new Date().getFullYear() })}
      </p>
    </footer>
  );
};

export default Footer;
