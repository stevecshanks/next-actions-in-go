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
    {action.dueBy ? (
      <Badge variant={dateBadgeVariant(action)}>
        {action.dueBy.toLocaleDateString()}
      </Badge>
    ) : null}
    <ActionImage url="https://trello-backgrounds.s3.amazonaws.com/SharedBackground/75x100/82e01e3a2ef841bebd6cdb064d2e95d3/photo-1527667250583-fcdbb18a523d.jpg" />
    {action.name}
  </ListGroup.Item>
);

const ActionImage: React.FC<{ url: string }> = ({ url }) => (
  <div
    className="action-image"
    style={{
      backgroundImage: `url("${url}")`,
    }}
  />
);
