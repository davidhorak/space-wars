import classNames from "classnames";
import React, { useEffect, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { useNavigation } from "react-router-dom";
import { MainLayoutProps } from "./types";
import { Header } from "../../components/header";
import { Footer } from "../../components/footer";

const MainLayout = <T extends object>(
  Component: React.ComponentType<T>,
  { className, i18n }: MainLayoutProps = {}
): React.ComponentType<T> =>
  function Render(props: T) {
    const { t } = useTranslation();
    const navigation = useNavigation();
    const disabled = useMemo(
      () => navigation.state !== "idle",
      [navigation.state]
    );

    useEffect(() => {
      document.title = t(i18n?.pageTitle ?? "title");
    }, [t]);

    return (
      <div
        className={classNames(
          "layout",
          "layout--main",
          "d-flex",
          "flex-column",
          "h-100",
          className
        )}
      >
        <Header />
        <Component disabled={disabled} {...props} />
        <Footer />
      </div>
    );
  };

export default MainLayout;
