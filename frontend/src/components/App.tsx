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
  imageUrl: string;
  projectName: string;
  dueBy?: string;
};

type JsonError = {
  detail: string;
};

type JsonResponse = {
  data?: JsonAction[];
  errors?: JsonError[];
};

const actionsFromJson = (json: JsonAction[]): Action[] =>
  json.map(
    (action) =>
      new Action({
        ...action,
        dueBy: action.dueBy ? new Date(action.dueBy) : undefined,
      })
  );

const errorsFromJson = (json: JsonError[]): String[] =>
  json.map((error) => `An error occurred: ${error.detail}`);

const App: React.FC = () => {
  const [actions, setActions] = useState<Action[]>([]);
  const [errorMessages, setErrorMessages] = useState<String[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const updateNotificationCount = () => {
    const notificationCount = actions.filter(
      (action) => action.isOverdue() || action.isDueSoon()
    ).length;
    const notificationText = notificationCount ? `(${notificationCount}) ` : "";
    document.title = `${notificationText}Next Actions`;
  };

  useEffect(() => {
    const fetchActions = async () => {
      setIsLoading(true);

      try {
        const response = await fetch("api/actions");
        const json = (await response.json()) as JsonResponse;
        setActions(actionsFromJson(json.data || []));
        setErrorMessages(errorsFromJson(json.errors || []));
      } catch {
        setErrorMessages(["An error occurred"]);
      }

      setIsLoading(false);
    };

    fetchActions();
    setInterval(fetchActions, ONE_HOUR);
  }, []);

  useEffect(updateNotificationCount, [actions]);

  return (
    <>
      {errorMessages.map((message, i) => (
        <Alert key={`error-${i}`} variant="danger">
          {message}
        </Alert>
      ))}
      <NextActionsList actions={actions} isLoading={isLoading} />
    </>
  );
};

export default App;
