import classNames from "classnames";
import { useTranslation } from "react-i18next";
import TactableLogo from "../../assets/tactable.svg";

const Header = (): JSX.Element => {
  const { t } = useTranslation();

  return (
    <header className={classNames("text-center", "mx-24")}>
      <img src={TactableLogo} alt="Tactable Logo" height={40} className={classNames("mt-24")} />
      <h1 className={classNames("h1")}>{t("title")}</h1>
    </header>
  );
};

export default Header;
