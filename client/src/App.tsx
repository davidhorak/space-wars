import i18n from "i18next";
import { RouterProvider } from "react-router-dom";
import { useEffect, useState } from "react";
import useAppStore from "./store/app";
import { GlobalError } from "./components/error";
import { initReactI18next } from "react-i18next";
import i18nSpaceWarsEn from "./i18n/space-wars.en";
import classNames from "classnames";
import { createRouter, routes, toRouterRoutes } from "./router";
import "./App.css";
import { Spinner } from "./components/spinner";

const development: boolean =
  !process.env.NODE_ENV || process.env.NODE_ENV === "development";
let hasLoaded = false;

function App() {
  const { status, setStatus, setError, error } = useAppStore();
  const [router, setRouter] = useState<ReturnType<typeof createRouter>>();

  useEffect(() => {
    if (hasLoaded) return;
    hasLoaded = true;

    const setup = async () => {
      setRouter(createRouter("", toRouterRoutes(routes)));

      try {
        await i18n.use(initReactI18next).init({
          resources: {
            en: { translation: i18nSpaceWarsEn },
          },
          lng: "en",
          interpolation: {
            escapeValue: false,
          },
          saveMissing: true,
          missingKeyHandler: (_, __, key: string) => {
            const error = new Error(`missing localization key ${key}`);
            if (development) {
              throw error;
            }
          },
        });
        setStatus("ready");
      } catch (error) {
        setError(error instanceof Error ? error : new Error(`${error}`));
        setStatus("error");
      }
    };

    setup();
  }, [setStatus, setError]);

  return (
    <>
      {status === "ready" && router && <RouterProvider router={router} />}
      {status === "initialization" && (
        <div
          className={classNames(
            "h-100",
            "d-flex",
            "align-items-center",
            "justify-content-center"
          )}
        >
          <Spinner />
        </div>
      )}
      {status === "error" && <GlobalError error={error} />}
    </>
  );
}

export default App;
