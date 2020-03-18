import React from "react";
import ListGroup from "react-bootstrap/ListGroup";
import { NextAction } from "./NextAction";

export interface Action {
  id: string;
  name: string;
}

type NextActionsListProps = {
  actions: Action[];
};

export const NextActionsList: React.FC<NextActionsListProps> = ({
  actions,
}) => (
  <ListGroup>
    {actions.map((action: Action) => (
      <NextAction key={action.id} action={action} />
    ))}
  </ListGroup>
);
