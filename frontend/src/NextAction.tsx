import React from "react";
import Badge from "react-bootstrap/Badge";
import ListGroup from "react-bootstrap/ListGroup";
import { Action } from "./NextActionsList";

type NextActionProps = {
  action: Action;
};

export const NextAction: React.FC<NextActionProps> = ({ action }) => (
  <ListGroup.Item>
    {action.name}
    {action.dueBy ? (
      <Badge pill variant="primary">
        {action.dueBy.toLocaleDateString()}
      </Badge>
    ) : null}
  </ListGroup.Item>
);
