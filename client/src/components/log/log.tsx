import classNames from "classnames";
import styles from "./Log.module.css";
import { LogProps } from "./types";
import { Trans } from "react-i18next";
import { isCollisionLog, isDamageLog, isKillLog } from "../../client/utils";
import { spaceshipColorClassName } from "./spaceshipColorClassName";

function Log({ log }: LogProps) {
  const time = log.time.replace(/^\d+.+\s/, "");

  const renderMessage = (log: Log) => {
    if (isDamageLog(log)) {
      return (
        <Trans
          i18nKey="views.home.log.damage"
          components={{
            1: <span className={classNames(spaceshipColorClassName(log.meta.who, 10))} />,
            2: <span className={classNames(spaceshipColorClassName(log.meta.whom, 10))} />,
          }}
          values={{
            damage: log.meta.damage,
            damageType: log.meta.damageType,
            who: log.meta.who,
            whom: log.meta.whom,
          }}
        />
      );
    }
    if (isCollisionLog(log)) {
      return (
        <Trans
          i18nKey="views.home.log.collision"
          components={{
            1: <span className={classNames(spaceshipColorClassName(log.meta.who, 10))} />,
            2: <span className={classNames(spaceshipColorClassName(log.meta.with, 10))} />,
          }}
          values={{
            who: log.meta.who,
            with: log.meta.with,
          }}
        />
      );
    }
    if (isKillLog(log)) {
      return (
        <Trans
          i18nKey="views.home.log.killed"
          components={{
            1: <span className={classNames(spaceshipColorClassName(log.meta.who, 10))} />,
            2: <span className={classNames(spaceshipColorClassName(log.meta.whom, 10))} />,
          }}
          values={{
            who: log.meta.who,
            whom: log.meta.whom,
          }}
        />
      );
    }
    return log.message;
  };

  return (
    <div className={classNames("d-flex")}>
      <div className={styles.log__time}>{time}</div>
      <div>{renderMessage(log)}</div>
    </div>
  );
}

export default Log;
