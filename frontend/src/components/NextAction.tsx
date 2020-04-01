import React from "react";
import Badge, { BadgeProps } from "react-bootstrap/Badge";
import ListGroup from "react-bootstrap/ListGroup";
import { Action } from "../models/Action";

type NextActionProps = {
  action: Action;
};

const dateBadgeVariant = (action: Action): BadgeProps["variant"] => {
  if (action.isOverdue()) {
    return "danger";
  }
  if (action.isDueSoon()) {
    return "warning";
  }
  return "primary";
};

export const NextAction: React.FC<NextActionProps> = ({ action }) => (
  <ListGroup.Item action href={action.url} target="_blank">
    {action.name}
    {action.dueBy ? (
      <Badge pill variant={dateBadgeVariant(action)}>
        {action.dueBy.toLocaleDateString()}
      </Badge>
    ) : null}
  </ListGroup.Item>
);
