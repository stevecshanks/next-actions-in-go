import React, { useState, useEffect } from "react";
import Alert from "react-bootstrap/Alert";
import { NextActionsList } from "./NextActionsList";
import { Action } from "../models/Action";
import "bootstrap/dist/css/bootstrap.min.css";
import "./App.css";

const ONE_HOUR = 60 * 60 * 1000;

type JsonAction = {
  id: string;
  name: string;
  url: string;
  dueBy?: string;
};

const actionsFromJson = (json: JsonAction[]): Action[] =>
  json.map(
    action =>
      new Action({
        ...action,
        dueBy: action.dueBy ? new Date(action.dueBy) : undefined,
      }),
  );

const App: React.FC = () => {
  const [actions, setActions] = useState<Action[]>([]);
  const [errorMessage, setErrorMessage] = useState<String | null>(null);

  const fetchActions = () => {
    fetch("api/actions")
      .then(response => response.json())
      .then(json => setActions(actionsFromJson(json.data)))
      .catch(() => setErrorMessage("An error occurred"));
  };

  const updateNotificationCount = () => {
    const notificationCount = actions.filter(
      action => action.isOverdue() || action.isDueSoon(),
    ).length;
    const notificationText = notificationCount ? `(${notificationCount}) ` : "";
    document.title = `${notificationText}Next Actions`;
  };

  useEffect(() => {
    fetchActions();
    setInterval(fetchActions, ONE_HOUR);
  }, []);

  useEffect(updateNotificationCount, [actions]);

  return (
    <>
      {errorMessage && <Alert variant="danger">{errorMessage}</Alert>}
      <NextActionsList actions={actions} />
    </>
  );
};

export default App;