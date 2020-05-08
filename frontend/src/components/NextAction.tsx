import React from "react";
import Badge, { BadgeProps } from "react-bootstrap/Badge";
import ListGroup from "react-bootstrap/ListGroup";
import Skeleton from "react-loading-skeleton";
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
    <ActionImage url={action.imageUrl} />
    <h1>{action.projectName}</h1>
    {action.name}
  </ListGroup.Item>
);

export const NextActionSkeleton: React.FC = () => (
  <ListGroup.Item action>
    <ActionImageSkeleton />
    <h1>
      <Skeleton width={"30%"} />
    </h1>
    <Skeleton width={"70%"} />
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

const ActionImageSkeleton: React.FC = () => (
  <div className="action-image-skeleton" />
);
