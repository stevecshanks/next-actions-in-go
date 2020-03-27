import React from "react";
import ListGroup from "react-bootstrap/ListGroup";
import { NextAction } from "./NextAction";
import { Action } from "../models/Action";

type NextActionsListProps = {
  actions: Action[];
};

const maxDate = new Date(8640000000000000);

const compareByDueDate = (a: Action, b: Action): number => {
  const aDueBy = a.dueBy || maxDate;
  const bDueBy = b.dueBy || maxDate;

  return aDueBy.getTime() - bDueBy.getTime();
};

const sortActions = (actions: Action[]): Action[] =>
  [...actions].sort(compareByDueDate);

export const NextActionsList: React.FC<NextActionsListProps> = ({
  actions,
}) => {
  const sortedActions = sortActions(actions);

  return (
    <ListGroup>
      {sortedActions.map((action: Action) => (
        <NextAction key={action.id} action={action} />
      ))}
    </ListGroup>
  );
};
