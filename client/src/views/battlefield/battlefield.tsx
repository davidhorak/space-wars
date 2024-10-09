import classNames from "classnames";
import { Trans, useTranslation } from "react-i18next";
import { CANVAS_ID, CANVAS_WIDTH, CANVAS_HEIGHT, FPS } from "../../client";
import { engine as createEngine } from "../../client";
import { useEffect, useState } from "react";
import useAppStore from "../../store/app";
import { Toggle } from "../../components/toggle";
import { Button } from "../../components/button";
import { Engine } from "../../client/engine";
import type { ScoreboardEntry } from "../../client/utils/scoreboard";
import { cloneDeep, get } from "lodash/fp";
import { spaceshipColorClassName } from "../../components/log/spaceshipColorClassName";
import { GameState } from "../../../../spaceships";
import { useSearchParams } from "react-router-dom";
import { Log } from "../../components/log";
import styles from "./battlefield.module.css";
import { DragAndDrop } from "../../components/dragAndDrop";
import { toError } from "../../utils/error";
import { Modal } from "../../components/modal";

let hasLoaded = false;
let engine: Engine;

const BattlefieldView = (): JSX.Element => {
  const { t } = useTranslation();
  const { setError, setStatus } = useAppStore();
  const [searchParams] = useSearchParams();

  const [gameState, setGameState] =
    useState<GameState["status"]>("initialized");
  const [logs, setLogs] = useState<GameState["logs"]>([]);
  const [scoreboard, setScoreboard] = useState<ScoreboardEntry[]>([]);
  const [winner, setWinner] = useState<[string, number]>();
  const [overallScoreboard, setOverallScoreboard] = useState<
    Map<
      string,
      Omit<ScoreboardEntry, "id" | "name" | "destroyed"> & { destroyed: number }
    >
  >(new Map());
  const [overallKills, setOverallKills] = useState<[string, number][]>([]);
  const [overallDestroyed, setOverallDestroyed] = useState<[string, number][]>(
    []
  );
  const [overallScore, setOverallScore] = useState<[string, number][]>([]);

  const [showColliders, setShowColliders] = useState(false);
  const [showEnergy, setShowEnergy] = useState(false);
  const [showHealth, setShowHealth] = useState(false);
  const [showNames, setShowNames] = useState(false);
  const [autoReset, setAutoReset] = useState(false);
  const [totalRounds, setTotalRounds] = useState(0);

  const [stateFilename, setStateFilename] = useState<string>();
  const [isProcessingStateFile, setIsProcessingStateFile] =
    useState<boolean>(false);

  const [showModal, setShowModal] = useState(false);
  const [modalTitle, setModalTitle] = useState("");
  const [modalBody, setModalBody] = useState<React.ReactNode>();

  useEffect(() => {
    if (hasLoaded) return;
    hasLoaded = true;

    const setup = async () => {
      const fps = searchParams.has("fps")
        ? parseInt(searchParams.get("fps")!)
        : FPS;
      const width = searchParams.has("width")
        ? parseInt(searchParams.get("width")!)
        : CANVAS_WIDTH;
      const height = searchParams.has("height")
        ? parseInt(searchParams.get("height")!)
        : CANVAS_HEIGHT;
      const seed = searchParams.has("seed")
        ? parseInt(searchParams.get("seed")!)
        : undefined;

      try {
        const go = new Go();
        const result = await WebAssembly.instantiateStreaming(
          fetch("space-wars.wasm", { cache: "no-cache" }),
          go.importObject
        );
        go.run(result.instance);

        if (seed) {
          spaceWars.init(width, height, seed);
        } else {
          spaceWars.init(width, height);
        }

        engine = await createEngine({
          canvasId: CANVAS_ID,
          width: width,
          height: height,
          fps: fps,
        });
        engine.onStateChanged(setGameState);
        engine.onLogsChanged(setLogs);
        engine.onScoreboardChanged(setScoreboard);

        setShowColliders(false);
        setShowNames(true);
        setShowHealth(true);
        setShowEnergy(true);
      } catch (error) {
        setError(error as Error);
        setStatus("error");
      }
    };

    setup();
  }, [setError, setStatus, searchParams]);

  useEffect(() => {
    if (!engine) return;
    engine.showCollider(showColliders);
    engine.showEnergy(showEnergy);
    engine.showHealth(showHealth);
    engine.showNames(showNames);
  }, [showColliders, showEnergy, showHealth, showNames]);

  useEffect(() => {
    if (gameState !== "ended") return;
    setTotalRounds((rounds) => rounds + 1);

    const first = get(0)(scoreboard);
    if (!first) {
      return;
    }

    if (!first.destroyed) {
      setWinner([first.name, first.score]);
      return;
    }

    const winners = scoreboard
      .filter((state) => state.score === scoreboard[0].score)
      .map((state) => state.name)
      .join(", ");
    setWinner([winners, first.score]);
  }, [scoreboard, gameState]);

  useEffect(() => {
    if (!winner) return;

    const overall = cloneDeep(overallScoreboard);
    scoreboard.forEach((entry) => {
      const overallEntry = overall.get(entry.name) ?? {
        score: 0,
        destroyed: 0,
        kills: 0,
      };

      overallEntry.score += entry.score;
      overallEntry.destroyed += entry.destroyed ? 1 : 0;
      overallEntry.kills += entry.kills;
      overall.set(entry.name, overallEntry);
    });

    setOverallScoreboard(overall);

    const entries = Array.from(overall.entries());
    entries.sort((a, b) => b[1].score - a[1].score);
    setOverallScore(entries.map(([name, entry]) => [name, entry.score]));
    entries.sort((a, b) => b[1].kills - a[1].kills);
    setOverallKills(entries.map(([name, entry]) => [name, entry.kills]));
    entries.sort((a, b) => a[1].destroyed - b[1].destroyed);
    setOverallDestroyed(
      entries.map(([name, entry]) => [name, entry.destroyed])
    );

    if (autoReset) {
      engine.reset();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [winner]);

  const saveState = () => {
    const state = engine.state();
    if (!state) return;

    const blob = new Blob([JSON.stringify(state, null, 2)], {
      type: "text/json",
    });

    const link = document.createElement("a");
    link.download = `space-wars-state-${Date.now()}.json`;
    link.href = window.URL.createObjectURL(blob);
    link.dataset.downloadurl = ["text/json", link.download, link.href].join(
      ":"
    );

    const evt = new MouseEvent("click", {
      view: window,
      bubbles: true,
      cancelable: true,
    });

    link.dispatchEvent(evt);
    link.remove();
  };

  const onLoadState = async (files: File[]) => {
    const file = files[0];
    if (!file) return;

    setIsProcessingStateFile(true);
    engine.pause();

    try {
      const text = await file.text();
      engine.setStartState(text);
      setStateFilename(file.name);
    } catch (error) {
      setModalTitle("Error");
      setModalBody(
        <p className={classNames("text-center")}>{toError(error).message}</p>
      );
      setShowModal(true);
    }
    setIsProcessingStateFile(false);
  };

  const removeState = () => {
    engine.removeStartState();
    setStateFilename(undefined);
  };

  return (
    <div
      className={classNames(
        "d-flex",
        "h-100",
        "pt-24",
        "flex-grow-1",
        "overflow-hidden"
      )}
    >
      <div className={classNames(styles.body__left, "pl-24 pr-12")}>
        {/* Options */}
        <h2 className={classNames("h5")}>
          {t("views.battlefield.options.title")}
        </h2>
        <div className={classNames("d-flex mt-12")}>
          <Toggle
            checked={showColliders}
            onChange={setShowColliders}
            disabled={isProcessingStateFile}
          />
          <h3 className={classNames("h6 align-self-center pl-12 pt-2")}>
            {t("views.battlefield.options.colliders")}
          </h3>
        </div>
        <div className={classNames("d-flex mt-12")}>
          <Toggle
            checked={showEnergy}
            onChange={setShowEnergy}
            disabled={isProcessingStateFile}
          />
          <h3 className={classNames("h6 align-self-center pl-12 pt-2")}>
            {t("views.battlefield.options.energy")}
          </h3>
        </div>
        <div className={classNames("d-flex mt-12")}>
          <Toggle
            checked={showHealth}
            onChange={setShowHealth}
            disabled={isProcessingStateFile}
          />
          <h3 className={classNames("h6 align-self-center pl-12 pt-2")}>
            {t("views.battlefield.options.health")}
          </h3>
        </div>
        <div className={classNames("d-flex mt-12")}>
          <Toggle
            checked={showNames}
            onChange={setShowNames}
            disabled={isProcessingStateFile}
          />
          <h3 className={classNames("h6 align-self-center pl-12 pt-2")}>
            {t("views.battlefield.options.names")}
          </h3>
        </div>
        <div className={classNames("d-flex mt-12")}>
          <Toggle
            checked={autoReset}
            onChange={setAutoReset}
            disabled={isProcessingStateFile}
          />
          <h3 className={classNames("h6 align-self-center pl-12 pt-2")}>
            {t("views.battlefield.options.autoReset")}
          </h3>
        </div>
        {/* Actions */}
        <div className={classNames("mt-12")}>
          <Button
            className="w-100"
            onClick={() => engine.start()}
            disabled={
              isProcessingStateFile ||
              gameState === "running" ||
              gameState === "ended"
            }
          >
            {t("views.battlefield.actions.start")}
          </Button>
        </div>
        <div className={classNames("mt-12")}>
          <Button
            className="w-100"
            onClick={() => engine.pause()}
            disabled={isProcessingStateFile || gameState !== "running"}
          >
            {t("views.battlefield.actions.pause")}
          </Button>
        </div>
        <div className={classNames("mt-12")}>
          <Button
            className="w-100"
            onClick={() => engine.step()}
            disabled={
              isProcessingStateFile ||
              gameState === "running" ||
              gameState === "ended"
            }
          >
            {t("views.battlefield.actions.step")}
          </Button>
        </div>
        <div className={classNames("mt-12")}>
          <Button
            className="w-100"
            onClick={() => engine.reset()}
            disabled={isProcessingStateFile}
          >
            {t("views.battlefield.actions.restart")}
          </Button>
        </div>
        <div className={classNames("mt-12")}>
          <Button
            className="w-100"
            onClick={saveState}
            disabled={isProcessingStateFile || gameState !== "paused"}
          >
            {t("views.battlefield.actions.saveState")}
          </Button>
        </div>
        {/* State file */}
        {stateFilename && (
          <>
            <h3 className={classNames("h6", "mt-24")}>
              {t("views.battlefield.load.loadedTitle")}
            </h3>
            <p className={classNames("mt-6")}>{stateFilename}</p>
            <div className={classNames("mt-12")}>
              <Button className="w-100" onClick={removeState}>
                {t("views.battlefield.load.remove")}
              </Button>
            </div>
          </>
        )}
        {!stateFilename && (
          <>
            <h3 className={classNames("h6", "mt-24")}>
              {t("views.battlefield.load.title")}
            </h3>
            <div className={classNames("mt-12")}>
              <DragAndDrop
                multiple={false}
                i18n={{
                  drop: t("views.battlefield.load.drop"),
                  processing: t("views.battlefield.load.processing"),
                  select: {
                    label: t("views.battlefield.load.select.label"),
                    title: t("views.battlefield.load.select.title"),
                  },
                }}
                filter={["application/json"]}
                processing={isProcessingStateFile}
                onChange={onLoadState}
              />
            </div>
          </>
        )}
      </div>
      <div
        className={classNames(
          styles.body__center,
          "flex-grow-1",
          "overflow-y-auto",
          "d-flex",
          "flex-column",
          "position-relative"
        )}
      >
        {/* Canvas */}
        <canvas
          className={classNames("align-self-center")}
          id={CANVAS_ID}
          width={CANVAS_WIDTH}
          height={CANVAS_HEIGHT}
        />
        {/* Logs */}
        <h2 className={classNames("h5", "mt-24")}>
          {t("views.battlefield.log.title")}
        </h2>
        <div className={classNames(styles.body__center__logs, "mt-2")}>
          {logs.map((log) => (
            <Log key={log.id} log={log} />
          ))}
        </div>
        {/* Game Over */}
        {gameState === "ended" && winner && (
          <div
            className={classNames(
              styles["body__center__game-over"],
              "position-absolute",
              "top-0",
              "left-0",
              "right-0",
              "text-center",
              "mt-96"
            )}
          >
            <h2 className={classNames("h2")}>
              {t("views.battlefield.gameOver.title")}
            </h2>
            <p className={classNames("h6")}>
              <Trans
                i18nKey="views.battlefield.gameOver.winner"
                components={{
                  1: (
                    <span
                      className={classNames(
                        spaceshipColorClassName(winner[0], 10)
                      )}
                    />
                  ),
                }}
                values={{
                  winner: winner[0],
                }}
              />
            </p>
            <p className={classNames("h6")}>
              {t("views.battlefield.gameOver.score", {
                score: winner[1],
              })}
            </p>
          </div>
        )}
      </div>
      <div className={classNames(styles.body__right, "pr-24", "pl-12")}>
        <h3 className={classNames("h5")}>
          {t("views.battlefield.scoreboard.title")}
        </h3>
        <div className={classNames("mt-12")}>
          {scoreboard.map((state) => (
            <div
              key={state.id}
              className={classNames(
                styles.scoreboard__item,
                {
                  [styles["scoreboard__item--destroyed"]]: state.destroyed,
                },
                "mt-2"
              )}
            >
              <span
                className={classNames(spaceshipColorClassName(state.name, 10))}
              >
                {state.name}
              </span>
              <span>{state.score}</span>
            </div>
          ))}
        </div>
        {totalRounds > 0 && (
          <>
            <h3 className={classNames("h5", "mt-24")}>
              {t("views.battlefield.scoreboard.overall")}
            </h3>
            <div className={classNames("mt-6")}>
              <span>
                {t("views.battlefield.scoreboard.totalRounds", { totalRounds })}
              </span>
            </div>
            <div className={classNames("mt-6")}>
              {overallScore.map(([name, score]) => (
                <div
                  key={name}
                  className={classNames(styles.scoreboard__item, "mt-2")}
                >
                  <span
                    className={classNames(spaceshipColorClassName(name, 10))}
                  >
                    {name}
                  </span>
                  <span>{score}</span>
                </div>
              ))}
            </div>
            <h3 className={classNames("h6", "mt-12")}>
              {t("views.battlefield.scoreboard.kills")}
            </h3>
            <div className={classNames("mt-6")}>
              {overallKills.map(([name, kills]) => (
                <div
                  key={name}
                  className={classNames(styles.scoreboard__item, "mt-2")}
                >
                  <span
                    className={classNames(spaceshipColorClassName(name, 10))}
                  >
                    {name}
                  </span>
                  <span>{kills}</span>
                </div>
              ))}
            </div>
            <h3 className={classNames("h6", "mt-12")}>
              {t("views.battlefield.scoreboard.destroyed")}
            </h3>
            <div className={classNames("mt-6")}>
              {overallDestroyed.map(([name, destroyed]) => (
                <div
                  key={name}
                  className={classNames(styles.scoreboard__item, "mt-2")}
                >
                  <span
                    className={classNames(spaceshipColorClassName(name, 10))}
                  >
                    {name}
                  </span>
                  <span>{destroyed}</span>
                </div>
              ))}
            </div>
          </>
        )}
      </div>
      {/* Modal */}
      {showModal && (
        <Modal title={modalTitle} onClose={() => setShowModal(false)}>
          {modalBody}
        </Modal>
      )}
    </div>
  );
};

export default BattlefieldView;
