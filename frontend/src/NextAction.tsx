import React from "react";
import ListGroup from "react-bootstrap/ListGroup";
import { Action } from "./NextActionsList";

type NextActionProps = {
  action: Action;
};

export const NextAction: React.FC<NextActionProps> = ({ action }) => (
  <ListGroup.Item>{action.name}</ListGroup.Item>
);
