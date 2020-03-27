import React from "react";
import Badge, { BadgeProps } from "react-bootstrap/Badge";
import ListGroup from "react-bootstrap/ListGroup";
import { Action } from "../models/Action";

type NextActionProps = {
  action: Action;
};

const twentyFourHours = 24 * 60 * 60 * 1000;

const isPast = (dueBy: Date): boolean => dueBy.getTime() < Date.now();

const isSoon = (dueBy: Date): boolean =>
  dueBy.getTime() < Date.now() + twentyFourHours;

const dateBadgeVariant = (dueBy: Date): BadgeProps["variant"] => {
  if (isPast(dueBy)) {
    return "danger";
  }
  if (isSoon(dueBy)) {
    return "warning";
  }
  return "primary";
};

export const NextAction: React.FC<NextActionProps> = ({ action }) => (
  <ListGroup.Item>
    {action.name}
    {action.dueBy ? (
      <Badge pill variant={dateBadgeVariant(action.dueBy)}>
        {action.dueBy.toLocaleDateString()}
      </Badge>
    ) : null}
  </ListGroup.Item>
);
